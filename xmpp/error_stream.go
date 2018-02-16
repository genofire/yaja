package xmpp

import "encoding/xml"

// StreamErrorGroup implements RFC 6120  A.2  Stream Error Namespace
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
	SeeOtherHost           string    `xml:"urn:ietf:params:xml:ns:xmpp-streams see-other-host,omitempty"`
	SystemShutdown         *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams system-shutdown"`
	UndefinedCondition     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams undefined-condition"`
	UnsupportedEncoding    *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-encoding"`
	UnsupportedStanzaType  *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-stanza-type"`
	UnsupportedVersion     *xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams unsupported-version"`
}

// StreamError implements RFC 6120  A.2  Stream Error Namespace
type StreamError struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Text    *Text
	StreamErrorGroup

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
