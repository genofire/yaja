package messages

import "encoding/xml"

type StreamErrorGroup struct {
	BadFormat              *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams bad-format"`
	BadNamespacePrefix     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams bad-namespace-prefix"`
	Conflict               *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams conflict"`
	ConnectionTimeout      *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams connection-timeout"`
	HostGone               *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams host-gone"`
	HostUnknown            *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams host-unknown"`
	ImproperAddressing     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams improper-addressing"`
	InternalServerError    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams internal-server-error"`
	InvalidFrom            *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-from"`
	InvalidID              *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-id"`
	InvalidNamespace       *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-namespace"`
	InvalidXML             *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams invalid-xml"`
	NotAuthorized          *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams not-authorized"`
	NotWellFormed          *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams not-well-formed"`
	PolicyViolation        *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams policy-violation"`
	RemoteConnectionFailed *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams remote-connection-failed"`
	Reset                  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams reset"`
	ResourceConstraint     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams resource-constraint"`
	RestrictedXML          *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams restricted-xml"`
	SeeOtherHost           string    `xml:"urn:ietf:params:xml:ns:xmpp-streams see-other-host"`
	SystemShutdown         *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams system-shutdown"`
	UndefinedCondition     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams undefined-condition"`
	UnsupportedEncoding    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-encoding"`
	UnsupportedStanzaType  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-stanza-type"`
	UnsupportedVersion     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-version"`
}

type StreamError struct {
	XMLName xml.Name     `xml:"http://etherx.jabber.org/streams error"`
	Text    string       `xml:"urn:ietf:params:xml:ns:xmpp-streams text"`
	Other   []XMLElement `xml:",any"`
	StreamErrorGroup
}

type ErrorClientType string

const (
	ErrorClientTypeAuth     ErrorClientType = "auth"
	ErrorClientTypeCancel   ErrorClientType = "cancel"
	ErrorClientTypeContinue ErrorClientType = "continue"
	ErrorClientTypeModify   ErrorClientType = "motify"
	ErrorClientTypeWait     ErrorClientType = "wait"
)

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
}

// ErrorClient element
type ErrorClient struct {
	XMLName xml.Name        `xml:"jabber:client error"`
	Code    string          `xml:"code,attr"`
	Type    ErrorClientType `xml:"type,attr"`
	Text    string          `xml:"text"`
	Other   []XMLElement    `xml:",any"`
	StanzaErrorGroup
}
