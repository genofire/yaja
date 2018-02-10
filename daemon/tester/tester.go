package tester

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/client"
	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type Tester struct {
	mainClient *client.Client
	Accounts   map[string]string  `json:"accounts"`
	Status     map[string]*Status `json:"-"`
}

func NewTester() *Tester {
	return &Tester{
		Accounts: make(map[string]string),
		Status:   make(map[string]*Status),
	}
}

func (t *Tester) Start(mainClient *client.Client, password string) {

	t.mainClient = mainClient

	status := NewStatus(mainClient.JID, password)
	status.client = mainClient
	status.Login = true
	status.Update()

	t.Status[mainClient.JID.Bare()] = status
	go t.StartBot(status)

	for jidString, passwd := range t.Accounts {
		t.Connect(jidString, passwd)
	}
}
func (t *Tester) Close() {
	for _, s := range t.Status {
		s.Login = false
		s.client.Close()
	}
}

func (t *Tester) Connect(jidString, password string) {
	logCTX := log.WithField("jid", jidString)
	jid := model.NewJID(jidString)
	status, ok := t.Status[jidString]
	if !ok {
		status = NewStatus(jid, password)
		t.Status[jidString] = status
	} else if status.JID == nil {
		status.JID = jid
	}
	if status.Login {
		logCTX.Warn("is already loggedin")
		return
	}
	c, err := client.NewClient(jid, password)
	if err != nil {
		logCTX.Warnf("could not connect client: %s", err)
	} else {
		logCTX.Info("client connected")
		status.Login = true
		status.client = c
		status.Update()
		go t.StartBot(status)
	}
}

func (t *Tester) UpdateConnectionStatus(from, to *model.JID, recvmsg string) {
	logCTX := log.WithFields(log.Fields{
		"jid":     to.Full(),
		"from":    from.Full(),
		"recvmsg": recvmsg,
	})

	status, ok := t.Status[from.Bare()]
	if !ok {
		logCTX.Warn("recv wrong msg")
		return
	}
	msg, ok := status.MessageForConnection[to.Bare()]
	logCTX = logCTX.WithField("msg", msg)
	if !ok || msg != recvmsg {
		logCTX.Warn("recv wrong msg")
		return
	}
	delete(status.MessageForConnection, to.Bare())
	status.Connections[to.Bare()] = true
	logCTX.Info("recv msg")

}

func (t *Tester) CheckStatus() {
	send := 0
	online := 0
	connection := 0
	for ownJID, own := range t.Status {
		logCTX := log.WithField("jid", ownJID)
		if !own.Login {
			pass, ok := t.Accounts[ownJID]
			if ok {
				t.Connect(ownJID, pass)
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
				logCTXTo = logCTXTo.WithField("old-msg", msg)
				own.Connections[jid] = false
				if ok, exists := own.Connections[jid]; !exists || ok {
					logCTXTo.Warn("could not recv msg")
				} else {
					logCTXTo.Debug("could not recv msg")
				}
			}
			msg = utils.CreateCookieString()
			logCTXTo = logCTXTo.WithField("msg", msg)

			own.client.Send(&messages.MessageClient{
				Body: "checkmsg " + msg,
				Type: messages.ChatTypeChat,
				To:   s.JID,
			})
			own.MessageForConnection[jid] = msg
			logCTXTo.Info("test send")
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
