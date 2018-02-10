package client

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"dev.sum7.eu/genofire/yaja/messages"
)

func (client *Client) auth(password string) error {
	f, err := client.startStream()
	if err != nil {
		return err
	}
	//auth:
	mechanism := ""
	for _, m := range f.Mechanisms.Mechanism {
		if m == "PLAIN" {
			mechanism = m
			// Plain authentication: send base64-encoded \x00 user \x00 password.
			raw := "\x00" + client.JID.Local + "\x00" + password
			enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
			base64.StdEncoding.Encode(enc, []byte(raw))
			fmt.Fprintf(client.conn, "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n", messages.NSSASL, enc)
			break
		}
		if m == "DIGEST-MD5" {
			mechanism = m
			// Digest-MD5 authentication
			fmt.Fprintf(client.conn, "<auth xmlns='%s' mechanism='DIGEST-MD5'/>\n", messages.NSSASL)
			var ch string
			if err := client.ReadElement(&ch); err != nil {
				return err
			}
			b, err := base64.StdEncoding.DecodeString(string(ch))
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
			digest := saslDigestResponse(client.JID.Local, realm, password, nonce, cnonceStr, "AUTHENTICATE", digestURI, nonceCount)
			message := "username=\"" + client.JID.Local + "\", realm=\"" + realm + "\", nonce=\"" + nonce + "\", cnonce=\"" + cnonceStr +
				"\", nc=" + nonceCount + ", qop=" + qop + ", digest-uri=\"" + digestURI + "\", response=" + digest + ", charset=" + charset

			fmt.Fprintf(client.conn, "<response xmlns='%s'>%s</response>\n", messages.NSSASL, base64.StdEncoding.EncodeToString([]byte(message)))

			err = client.ReadElement(&ch)
			if err != nil {
				return err
			}
			_, err = base64.StdEncoding.DecodeString(ch)
			if err != nil {
				return err
			}
			fmt.Fprintf(client.conn, "<response xmlns='%s'/>\n", messages.NSSASL)
			break
		}
	}
	if mechanism == "" {
		return fmt.Errorf("PLAIN authentication is not an option: %v", f.Mechanisms.Mechanism)
	}
	element, err := client.Read()
	if err != nil {
		return err
	}
	if element.Name.Local != "success" {
		return errors.New("auth failed: " + element.Name.Local)
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
