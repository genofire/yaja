package xmpp

const (
	// NSStream implements RFC 6120 - A.1 Stream Namespace
	NSStream = "http://etherx.jabber.org/streams"

	// NSStreamError implements RFC 6120 - A.2 Stream Error Namespace
	NSStreamError = "urn:ietf:params:xml:ns:xmpp-streams"

	// NSStartTLS implements RFC 6120 - A.3 StartTLS Namespace
	NSStartTLS = "urn:ietf:params:xml:ns:xmpp-tls"

	// NSSASL implements RFC 6120 - A.4 SASL Namespace
	NSSASL = "urn:ietf:params:xml:ns:xmpp-sasl"

	// NSClient implements RFC 6120 - A.5 Client Namespace
	NSClient = "jabber:client"

	// NSServer implements RFC 6120 - A.6 Server Namespace
	NSServer = "jabber:server"

	// NSBind implements RFC 6120 - A.7 Resource Binding Namespace
	NSBind = "urn:ietf:params:xml:ns:xmpp-bind"

	// NSStanzaError implements RFC 6120 - A.8 Stanza Error Binding Namespace
	NSStanzaError = "urn:ietf:params:xml:ns:xmpp-stanzas"
)
