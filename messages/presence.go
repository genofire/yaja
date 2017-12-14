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

// Presence element
type Presence struct {
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
	Error *Error `xml:"error"`
	// Delay    Delay        `xml:"delay"`
}
