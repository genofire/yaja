package tester

import (
	"fmt"
	"strings"

	"github.com/FreifunkBremen/yanic/lib/jsontime"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

type Output struct {
	Timestamp jsontime.Time `json:"timestamp"`
	Status    []*Status     `json:"nodes"`
	Links     []*Link       `json:"links"`
}

type Link struct {
	Source     string `json:"source"`
	SourceJID  string `json:"source_jid"`
	Target     string `json:"target"`
	TargetJID  string `json:"target_jid"`
	FromSource bool   `json:"from_source"`
	FromTarget bool   `json:"from_target"`
}

func (t *Tester) Output() *Output {
	output := &Output{
		Timestamp: jsontime.Now(),
		Status:    make([]*Status, 0),
		Links:     make([]*Link, 0),
	}
	links := make(map[string]*Link)

	t.mux.Lock()
	defer t.mux.Unlock()

	for from, status := range t.Status {
		output.Status = append(output.Status, status)
		if !status.Login {
			continue
		}
		for to, linkOK := range status.Connections {
			var key string
			// keep source and target in the same order
			switchSourceTarget := strings.Compare(from, to) > 0
			if switchSourceTarget {
				key = fmt.Sprintf("%s-%s", from, to)
			} else {
				key = fmt.Sprintf("%s-%s", to, from)
			}
			if link := links[key]; link != nil {
				if switchSourceTarget {
					link.FromTarget = linkOK
				} else {
					link.FromSource = linkOK
				}
				continue
			}
			toJID := xmppbase.NewJID(to)
			link := &Link{
				Source:     status.JID.Domain,
				SourceJID:  status.JID.Bare().String(),
				Target:     toJID.Domain,
				TargetJID:  toJID.Bare().String(),
				FromSource: linkOK,
				FromTarget: false,
			}
			if switchSourceTarget {
				link.Source = toJID.Domain
				link.SourceJID = toJID.Bare().String()
				link.Target = status.JID.Domain
				link.TargetJID = status.JID.Bare().String()
				link.FromSource = false
				link.FromTarget = linkOK
			}
			links[key] = link
			output.Links = append(output.Links, link)
		}
	}
	return output
}
