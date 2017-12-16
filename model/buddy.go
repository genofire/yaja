package model

const (
	SubscriptionNone = iota
	SubscriptionTo
	SubscriptionFrom
	SubscriptionBoth
	AskNone = iota
	AskSubscribe
)

type Buddy struct {
	Name         string   `json:"name"`
	Groups       []string `json:"groups"`
	Subscription int      `json:"subscription"`
	Ask          int      `json:"ask"`
}
