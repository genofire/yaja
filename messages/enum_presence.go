package messages

type PresenceType string

// RFC 6120 part of A.5 Client Namespace and A.6 Server Namespace
const (
	PresenceTypeError        PresenceType = "error"
	PresenceTypeProbe        PresenceType = "probe"
	PresenceTypeSubscribe    PresenceType = "subscribe"
	PresenceTypeSubscribed   PresenceType = "subscribed"
	PresenceTypeUnavailable  PresenceType = "unavailable"
	PresenceTypeUnsubscribe  PresenceType = "unsubscribe"
	PresenceTypeUnsubscribed PresenceType = "unsubscribed"
)

type PresenceShow string

// RFC 6120 part of A.5 Client Namespace and A.6 Server Namespace
const (
	PresenceShowAway PresenceShow = "away"
	PresenceShowChat PresenceShow = "chat"
	PresenceShowDND  PresenceShow = "dnd"
	PresenceShowXA   PresenceShow = "xa"
)
