package xmpp

import "encoding/xml"

// StanzaErrorGroup implements RFC 6120  A.8  Resource binding namespace
type StanzaErrorGroup struct {
	BadRequest            *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas bad-request"`
	Conflict              *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas conflict"`
	FeatureNotImplemented *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas feature-not-implemented"`
	Forbidden             *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas forbidden"`
	Gone                  string    `xml:"urn:ietf:params:xml:ns:xmpp-stanzas gone"`
	InternalServerError   *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas internal-server-error"`
	ItemNotFound          *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas item-not-found"`
	JIDMalformed          *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas jid-malformed"`
	NotAcceptable         *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-acceptable"`
	NotAllowed            *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-allowed"`
	NotAuthorized         *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-authorized"`
	PolicyViolation       *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas policy-violation"`
	RecipientUnavailable  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas recipient-unavailable"`
	Redirect              string    `xml:"urn:ietf:params:xml:ns:xmpp-stanzas redirect"`
	RegistrationRequired  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas registration-required"`
	RemoteServerNotFound  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas remote-server-not-found"`
	RemoteServerTimeout   *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas remote-server-timeout"`
	ResourceConstraint    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas resource-constraint"`
	ServiceUnavailable    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas service-unavailable"`
	SubscriptionRequired  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas subscription-required"`
	UndefinedCondition    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas undefined-condition"`
	UnexpectedRequest     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas unexpected-request"`

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
