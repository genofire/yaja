package xmppiq

const (
	// NSDiscoInfo implements XEP 0030: Service Discovery - 11.1 disco#info
	NSDiscoInfo = "http://jabber.org/protocol/disco#info"

	// NSDiscoItems implements XEP 0030: Service Discovery - 11.2 disco#items
	NSDiscoItems = "http://jabber.org/protocol/disco#items"

	// NSPrivateXMLStorage implements XEP-0049: Private XML Storage - 7
	NSPrivateXMLStorage = "jabber:iq:private"

	// NSVCard implements XEP-0054: vcard-temp - 14 (WIP)
	NSVCard = "vcard-temp"

	// NSVersion implements XEP-0092: Software Version - 4
	NSVersion = "jabber:iq:version"

	// NSPing implements XEP-0199: XMPP Ping - 10
	NSPing = "urn:xmpp:ping"

	// NSRegister implements XEP-0077: In-Band Registration - 14.1 jabber:iq:register
	NSRegister = "jabber:iq:register"

	// NSFeatureRegister implements XEP-0077: In-Band Registration - 14.2 Stream Feature
	NSFeatureRegister = "http://jabber.org/features/iq-register"
)
