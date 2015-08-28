package irc2mail

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/thoj/go-ircevent"
)

type Bot struct {
	Config  Config
	Notices []*Notice
}

func NewBot(c Config) (*Bot, error) {
	notices := []*Notice{}
	for key, v := range c.Notice {
		for channel, noticeConfig := range v {
			gs, err := NewGmailService(
				c.OAuth2.ClientID,
				c.OAuth2.ClientSecret,
				c.OAuth2.RefreshToken,
			)
			if err != nil {
				return nil, err
			}

			notices = append(notices, &Notice{
				GmailService: gs,
				Key:          key,
				Channel:      channel,
				Config:       noticeConfig,
				DebugLog:     c.DebugLog,
				IsAway:       false,
			})
		}
	}
	return &Bot{
		Config:  c,
		Notices: notices,
	}, nil
}

func (b *Bot) newIRCConn() (*irc.Connection, error) {
	irccon := irc.IRC(
		b.Config.IRCEvent.NickName,
		b.Config.IRCEvent.UserName,
	)
	irccon.Debug = b.Config.IRCEvent.Debug
	irccon.UseTLS = b.Config.IRCEvent.UseSSL
	if irccon.UseTLS {
		irccon.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	irccon.Password = b.Config.IRCEvent.Password

	err := irccon.Connect(fmt.Sprintf(
		"%s:%d",
		b.Config.IRCEvent.Host,
		b.Config.IRCEvent.Port,
	))
	if err != nil {
		return nil, err
	}

	irccon.AddCallback("001", func(e *irc.Event) {
		for _, channel := range b.Channels() {
			b.Debugf("join: %s", channel)
			irccon.Join(channel)
		}
	})
	irccon.AddCallback("PART", func(e *irc.Event) {
		for _, notice := range b.Notices {
			notice.Part(e)
		}
	})
	irccon.AddCallback("353", func(e *irc.Event) {
		for _, notice := range b.Notices {
			notice.Init(e)
		}
	})
	irccon.AddCallback("NICK", func(e *irc.Event) {
		for _, notice := range b.Notices {
			notice.Nick(e)
		}
	})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		for _, notice := range b.Notices {
			notice.Msg(e)
		}
	})

	return irccon, nil
}

func (b *Bot) Debugf(format string, args ...interface{}) {
	if b.Config.DebugLog {
		log.Printf(format, args...)
	}
}

func (b *Bot) Run() error {
	irccon, err := b.newIRCConn()
	if err != nil {
		return err
	}
	for {
		irccon.Loop()
	}
	return nil
}

func (b *Bot) Channels() []string {
	channelMap := make(map[string]struct{})
	channels := []string{}
	for _, v := range b.Config.Notice {
		for channel, _ := range v {
			channelMap[channel] = struct{}{}
		}
	}
	for channel, _ := range channelMap {
		channels = append(channels, channel)
	}
	return channels
}
