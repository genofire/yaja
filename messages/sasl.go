package messages

import (
	"encoding/xml"
)

// RFC 6120 - A.4 SASL Namespace
type SASLMechanisms struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism []string `xml:"mechanism"`
}

// SASLAuth element
type SASLAuth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Body      string   `xml:",chardata"`
}

// SASLChallenge element
type SASLChallenge struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Body    string   `xml:",chardata"`
}

// SASLResponse element
type SASLResponse struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Body    string   `xml:",chardata"`
}

// SASLSuccess element
type SASLSuccess struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
	Body    string   `xml:",chardata"`
}

// SASLAbout element
type SASLAbout struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

// RFC 6120 - A.4 SASL Namespace
type SASLFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`

	Aborted              *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl aborted"`
	AccountDisabled      *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl account-disabled"`
	CredentialsExpired   *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl credentials-expired"`
	EncryptionRequired   *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl encryption-required"`
	IncorrectEncoding    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl incorrect-encoding"`
	InvalidAuthzid       *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl invalid-authzid"`
	InvalidMechanism     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl invalid-mechanism"`
	MalformedRequest     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl malformed-request"`
	MechanismTooWeak     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanism-too-weak"`
	NotAuthorized        *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl not-authorized"`
	TemporaryAuthFailure *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl temporary-auth-failure"`

	Text *Text
}
