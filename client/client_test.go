package client

/*
func TestClient(t *testing.T) {
	assert := assert.New(t)

	jid := xmppbase.NewJID("test@example.net")

	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	client := &Client{
		JID:     jid,
		Timeout: time.Millisecond * 500,
		Logging: logger.WithField("jid", jid.String()),
	}
	// close nil connected
	assert.NoError(client.Close())

	err := client.Connect("password")
	assert.Error(err)
	assert.Contains(err.Error(), "timeout")

	jid.Domain = "chat.sum7.eu"

	// invalid auth
	client, err = NewClient(jid, "password")
	assert.NotNil(client)
	assert.Error(err)
	assert.Contains(err.Error(), "not-authorized : ")
	// already closed
	assert.Error(client.Close())

	client.Logging = logger.WithField("jid", jid.String())

	err = client.Connect("FqzMp6bevlHlt8d")
	assert.NoError(err)

}
*/
