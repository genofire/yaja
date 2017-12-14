package server

// State processes the stream and moves to the next state
type State interface {
	Process(client *Client) (State, *Client)
}
