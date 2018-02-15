package tester

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/client"
	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

type Tester struct {
	mainClient     *client.Client
	Timeout        time.Duration       `json:"-"`
	Accounts       map[string]*Account `json:"accounts"`
	Status         map[string]*Status  `json:"-"`
	mux            sync.Mutex
	LoggingClients *log.Logger     `json:"-"`
	LoggingBots    log.Level       `json:"-"`
	Admins         []*xmppbase.JID `json:"-"`
}

func NewTester() *Tester {
	return &Tester{
		Accounts: make(map[string]*Account),
		Status:   make(map[string]*Status),
	}
}

func (t *Tester) Start(mainClient *client.Client, password string) {

	t.mainClient = mainClient

	status := NewStatus(mainClient, &Account{
		JID:      mainClient.JID,
		Password: password,
	})
	status.client = mainClient
	status.Login = true
	status.Update(t.Timeout)

	t.mux.Lock()
	defer t.mux.Unlock()

	t.Status[mainClient.JID.Bare().String()] = status
	go t.StartBot(status)

	for _, acc := range t.Accounts {
		t.Connect(acc)
	}
}
func (t *Tester) Close() {
	for _, s := range t.Status {
		s.Disconnect("yaja tester stopped")
	}
}

func (t *Tester) Connect(acc *Account) {
	logCTX := log.WithField("jid", acc.JID.Full().String())
	bare := acc.JID.Bare().String()
	status, ok := t.Status[bare]
	if !ok {
		status = NewStatus(t.mainClient, acc)
		t.Status[bare] = status
	} else if status.JID == nil {
		status.JID = acc.JID
	}
	if status.Login {
		logCTX.Warn("is already loggedin")
		return
	}
	c := &client.Client{
		Timeout: t.Timeout,
		JID:     acc.JID,
		Logging: t.LoggingClients,
	}
	err := c.Connect(acc.Password)
	if err != nil {
		logCTX.Warnf("could not connect client: %s", err)
	} else {
		logCTX.Info("client connected")
		status.Login = true
		status.client = c
		status.account.JID = c.JID
		status.JID = c.JID
		status.Update(t.Timeout)
		go t.StartBot(status)
	}
}

func (t *Tester) UpdateConnectionStatus(from, to *xmppbase.JID, recvmsg string) {
	logCTX := log.WithFields(log.Fields{
		"jid":      to.Full(),
		"from":     from.Full(),
		"msg-recv": recvmsg,
	})

	t.mux.Lock()
	defer t.mux.Unlock()

	status, ok := t.Status[from.Bare().String()]
	if !ok {
		logCTX.Warn("recv msg without receiver")
		return
	}
	toBare := to.Bare().String()
	msg, ok := status.MessageForConnection[toBare]
	logCTX = logCTX.WithField("msg-send", msg)
	if !ok || msg != recvmsg || msg == "" || recvmsg == "" {
		logCTX.Warn("recv wrong msg")
		return
	}
	delete(status.MessageForConnection, toBare)
	status.Connections[toBare] = true
	logCTX.Info("recv msg")

}

func (t *Tester) CheckStatus() {
	send := 0
	online := 0
	connection := 0

	t.mux.Lock()
	defer t.mux.Unlock()

	for ownJID, own := range t.Status {
		logCTX := log.WithField("jid", ownJID)
		if !own.Login {
			acc, ok := t.Accounts[ownJID]
			if ok {
				t.Connect(acc)
			}
			if !own.Login {
				continue
			}
		}
		online++
		for jid, s := range t.Status {
			logCTXTo := logCTX.WithField("to", jid)
			if jid == ownJID {
				continue
			}
			connection++
			if own.MessageForConnection == nil {
				own.MessageForConnection = make(map[string]string)
			}
			msg, ok := own.MessageForConnection[jid]
			if ok {
				logCTXTo = logCTXTo.WithField("msg-old", msg)
				own.Connections[jid] = false
				if ok, exists := own.Connections[jid]; !exists || ok {
					logCTXTo.Info("could not recv msg")
				} else {
					logCTXTo.Debug("could not recv msg")
				}
			}
			msg = xmpp.CreateCookieString()
			logCTXTo = logCTXTo.WithField("msg-send", msg)

			own.client.Send(&xmpp.MessageClient{
				Body: "checkmsg " + msg,
				Type: xmpp.MessageTypeChat,
				To:   s.JID,
			})
			own.MessageForConnection[s.JID.Bare().String()] = msg
			logCTXTo.Debug("test send")
			send++
		}
	}

	log.WithFields(log.Fields{
		"online":     online,
		"connection": connection,
		"send":       send,
		"all":        len(t.Status),
	}).Info("checked complete")
}
