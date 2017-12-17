package toclient

import (
	"github.com/genofire/yaja/server/extension"
	"github.com/genofire/yaja/server/state"
	"github.com/genofire/yaja/server/utils"
)

// SendingClient state
type SendingClient struct {
	Next state.State
}

// Process messages
func (state *SendingClient) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "normal")
	client.Log.Debug("sending")
	// sending
	go func() {
		select {
		case msg := <-client.Messages:
			err := client.Out.Encode(msg)
			if err != nil {
				client.Log.Warn(err)
			}
		case <-client.OnClose():
			return
		}
	}()
	client.Log.Debug("receiving")
	return state.Next, client
}

// ReceivingClient state
type ReceivingClient struct {
	Extensions extension.Extensions
}

// Process messages
func (state *ReceivingClient) Process(client *utils.Client) (state.State, *utils.Client) {
	element, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	state.Extensions.Process(element, client)
	return state, client
}
