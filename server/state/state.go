package state

import "github.com/genofire/yaja/server/utils"

// State processes the stream and moves to the next state
type State interface {
	Process(client *utils.Client) (State, *utils.Client)
}
