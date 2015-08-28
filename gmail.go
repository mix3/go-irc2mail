package irc2mail

import (
	"encoding/base64"

	"google.golang.org/api/gmail/v1"
	"gopkg.in/jpoehls/gophermail.v0"

	"golang.org/x/oauth2"
)

type GmailService struct {
	*gmail.Service
}

func NewGmailService(client_id, client_secret, refresh_token string) (*GmailService, error) {
	conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
	tok := &oauth2.Token{
		RefreshToken: refresh_token,
	}

	svc, err := gmail.New(conf.Client(oauth2.NoContext, tok))

	if err != nil {
		return nil, err
	}

	return &GmailService{svc}, nil
}

func (gs *GmailService) Send(to, subject, body string) (*gmail.Message, error) {
	m := &gophermail.Message{}
	m.SetFrom("from@dummy.com")
	m.AddTo(to)
	m.Subject = subject
	m.Body = body

	b, err := m.Bytes()
	if err != nil {
		return nil, err
	}

	return gs.Users.Messages.Send("me", &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(b),
	}).Do()
}
