package state

import "dev.sum7.eu/genofire/yaja/server/utils"

// State processes the stream and moves to the next state
type State interface {
	Process() State
}

// Start state
type Debug struct {
	Next   State
	Client *utils.Client
}

// Process message
func (state *Debug) Process() State {
	state.Client.Log = state.Client.Log.WithField("state", "debug")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	state.Client.Log.Info(element)

	return state.Next
}
