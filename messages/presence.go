package messages

import "encoding/xml"

type PresenceType string

const (
	PresenceTypeUnavailable  PresenceType = "unavailable"
	PresenceTypeSubscribe    PresenceType = "subscribe"
	PresenceTypeSubscribed   PresenceType = "subscribed"
	PresenceTypeUnsubscribe  PresenceType = "unsubscribe"
	PresenceTypeUnsubscribed PresenceType = "unsubscribed"
	PresenceTypeProbe        PresenceType = "probe"
	PresenceTypeError        PresenceType = "error"
)

// PresenceClient element
type PresenceClient struct {
	XMLName xml.Name `xml:"jabber:client presence"`
	From    string   `xml:"from,attr,omitempty"`
	ID      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`

	Show     string `xml:"show,omitempty"`   // away, chat, dnd, xa
	Status   string `xml:"status,omitempty"` // sb []clientText
	Priority string `xml:"priority,omitempty"`
	// Caps     *ClientCaps  `xml:"c"`
	Error *ErrorClient `xml:"error"`
	// Delay    Delay        `xml:"delay"`
}

// MessageClient element
type MessageClient struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr,omitempty"`
	ID      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Subject string   `xml:"subject"`
	Body    string   `xml:"body"`
	Thread  string   `xml:"thread"`
}
