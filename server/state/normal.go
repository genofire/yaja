package state

import (
	"dev.sum7.eu/genofire/yaja/server/extension"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

// SendingClient state
type SendingClient struct {
	Next   State
	Client *utils.Client
}

// Process messages
func (state *SendingClient) Process() State {
	state.Client.Log = state.Client.Log.WithField("state", "normal")
	state.Client.Log.Debug("sending")
	// sending
	go func() {
		select {
		case msg := <-state.Client.Messages:
			err := state.Client.Out.Encode(msg)
			if err != nil {
				state.Client.Log.Warn(err)
			}
		case <-state.Client.OnClose():
			return
		}
	}()
	state.Client.Log.Debug("receiving")
	return state.Next
}

// ReceivingClient state
type ReceivingClient struct {
	Extensions extension.Extensions
	Client     *utils.Client
}

// Process messages
func (state *ReceivingClient) Process() State {
	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	state.Extensions.Process(element, state.Client)
	return state
}
