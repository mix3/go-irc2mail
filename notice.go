package irc2mail

import (
	"fmt"
	"log"
	"strings"

	"github.com/thoj/go-ircevent"
)

type Notice struct {
	GmailService *GmailService
	Key          string
	Channel      string
	Config       ConfigNotice
	IsAway       bool
	DebugLog     bool
}

func (n *Notice) Debugf(format string, args ...interface{}) {
	if n.DebugLog {
		log.Printf(fmt.Sprintf("[%s][%s] ", n.Key, n.Channel)+format, args...)
	}
}

func (n *Notice) Init(e *irc.Event) {
	for _, v := range strings.Split(e.Arguments[3], " ") {
		if v == n.Config.AwayName {
			n.Debugf("away")
			n.IsAway = true
		}
	}
}

func (n *Notice) Part(e *irc.Event) {
	if n.Channel == e.Arguments[0] {
		n.Debugf("join %s", n.Channel)
		e.Connection.Join(n.Channel)
	}
}

func (n *Notice) Nick(e *irc.Event) {
	if e.Nick == n.Config.AwayName {
		n.Debugf("back")
		n.IsAway = false
	}
	for _, v := range e.Arguments {
		if v == n.Config.AwayName {
			n.Debugf("away")
			n.IsAway = true
		}
	}
}

func (n *Notice) Match(e *irc.Event) bool {
	if len(e.Arguments) <= 0 {
		n.Debugf("e.Arguments empty")
		return false
	}

	if e.Arguments[0] != n.Channel {
		n.Debugf("not match channel: %s", e.Arguments[0])
		return false
	}

	str := strings.ToUpper(e.Message())
	for _, v := range n.Config.Keywords {
		if 0 <= strings.Index(str, strings.ToUpper(v)) {
			return true
		}
	}

	n.Debugf("not match keyword")
	return false
}

func (n *Notice) Msg(e *irc.Event) {
	if n.Match(e) && (n.IsAway || !n.Config.CheckAway) {
		n.Debugf("send message")
		_, err := n.GmailService.Send(
			n.Config.To,
			fmt.Sprintf("%s %s %s", n.Config.SubjectPrefix, n.Key, n.Channel),
			e.Message(),
		)
		if err != nil {
			n.Debugf("error: %v", err)
		}
	} else {
		n.Debugf("skip message")
	}
}
