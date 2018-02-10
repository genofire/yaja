package client

import "crypto/tls"

func (client *Client) TLSConnectionState() *tls.ConnectionState {
	if tlsconn, ok := client.conn.(*tls.Conn); ok {
		state := tlsconn.ConnectionState()
		return &state
	}
	return nil
}
