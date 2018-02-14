package client

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"dev.sum7.eu/genofire/yaja/xmpp"
)

func (client *Client) auth(password string) error {
	f, err := client.startStream()
	if err != nil {
		return err
	}
	//auth:
	mechanism := ""
	challenge := &xmpp.SASLChallenge{}
	response := &xmpp.SASLResponse{}
	for _, m := range f.Mechanisms.Mechanism {
		client.Logging.Debugf("try auth with '%s'", m)
		if m == "SCRAM-SHA-1" {
			/*
				mechanism = m
				TODO
				break
			*/
		}

		if m == "DIGEST-MD5" {
			mechanism = m
			// Digest-MD5 authentication
			client.Send(&xmpp.SASLAuth{
				Mechanism: m,
			})
			if err := client.ReadDecode(challenge); err != nil {
				return err
			}
			b, err := base64.StdEncoding.DecodeString(challenge.Body)
			if err != nil {
				return err
			}
			tokens := map[string]string{}
			for _, token := range strings.Split(string(b), ",") {
				kv := strings.SplitN(strings.TrimSpace(token), "=", 2)
				if len(kv) == 2 {
					if kv[1][0] == '"' && kv[1][len(kv[1])-1] == '"' {
						kv[1] = kv[1][1 : len(kv[1])-1]
					}
					tokens[kv[0]] = kv[1]
				}
			}
			realm, _ := tokens["realm"]
			nonce, _ := tokens["nonce"]
			qop, _ := tokens["qop"]
			charset, _ := tokens["charset"]
			cnonceStr := cnonce()
			digestURI := "xmpp/" + client.JID.Domain
			nonceCount := fmt.Sprintf("%08x", 1)
			digest := saslDigestResponse(client.JID.Node, realm, password, nonce, cnonceStr, "AUTHENTICATE", digestURI, nonceCount)
			message := "username=\"" + client.JID.Node + "\", realm=\"" + realm + "\", nonce=\"" + nonce + "\", cnonce=\"" + cnonceStr +
				"\", nc=" + nonceCount + ", qop=" + qop + ", digest-uri=\"" + digestURI + "\", response=" + digest + ", charset=" + charset

			response.Body = base64.StdEncoding.EncodeToString([]byte(message))
			client.Send(response)
			break
		}
		if m == "PLAIN" {
			mechanism = m
			// Plain authentication: send base64-encoded \x00 user \x00 password.
			raw := "\x00" + client.JID.Node + "\x00" + password
			enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
			base64.StdEncoding.Encode(enc, []byte(raw))
			client.Send(&xmpp.SASLAuth{
				Mechanism: "PLAIN",
				Body:      string(enc),
			})

			break
		}
	}
	if mechanism == "" {
		return fmt.Errorf("PLAIN authentication is not an option: %s", f.Mechanisms.Mechanism)
	}
	client.Logging.Infof("used auth with '%s'", mechanism)

	element, err := client.Read()
	if err != nil {
		return err
	}
	fail := xmpp.SASLFailure{}
	if err := client.Decode(&fail, element); err == nil {
		if txt := fail.Text; txt != nil {
			return errors.New(xmpp.XMLChildrenString(fail) + " : " + txt.Body)
		}
		return errors.New(xmpp.XMLChildrenString(fail))
	}
	if err := client.Decode(&xmpp.SASLSuccess{}, element); err != nil {
		return errors.New("auth failed - with unexpected answer")
	}
	return nil
}

func saslDigestResponse(username, realm, passwd, nonce, cnonceStr, authenticate, digestURI, nonceCountStr string) string {
	h := func(text string) []byte {
		h := md5.New()
		h.Write([]byte(text))
		return h.Sum(nil)
	}
	hex := func(bytes []byte) string {
		return fmt.Sprintf("%x", bytes)
	}
	kd := func(secret, data string) []byte {
		return h(secret + ":" + data)
	}

	a1 := string(h(username+":"+realm+":"+passwd)) + ":" + nonce + ":" + cnonceStr
	a2 := authenticate + ":" + digestURI
	response := hex(kd(hex(h(a1)), nonce+":"+nonceCountStr+":"+cnonceStr+":auth:"+hex(h(a2))))
	return response
}

func cnonce() string {
	randSize := big.NewInt(0)
	randSize.Lsh(big.NewInt(1), 64)
	cn, err := rand.Int(rand.Reader, randSize)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%016x", cn)
}
