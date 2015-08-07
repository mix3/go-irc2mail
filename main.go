package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/sirsean/go-mailgun/mailgun"
	"github.com/thoj/go-ircevent"
)

var (
	mg      *mailgun.Client
	irccon  *irc.Connection
	matcher = []string{}
)

func init() {
	// matcher
	for _, v := range strings.Split(os.Getenv("IRC_MATCHER"), " ") {
		matcher = append(matcher, strings.ToUpper(v))
	}

	// Mailgun
	apiKey := os.Getenv("MAILGUN_APIKEY")
	domain := os.Getenv("MAILGUN_DOMAIN")
	mg = mailgun.NewClient(apiKey, domain)

	// mail
	mailFromname := os.Getenv("MAIL_FROMNAME")
	mailFromaddr := os.Getenv("MAIL_FROMADDR")
	mailToaddr := os.Getenv("MAIL_TOADDR")
	mailSubject := os.Getenv("MAIL_SUBJECT")

	// ircevent
	channels := getChannels(os.Getenv("IRC_CHANNELS"))
	ircNickname := os.Getenv("IRC_NICKNAME")
	ircUsername := os.Getenv("IRC_USERNAME")
	ircPassword := os.Getenv("IRC_PASSWORD")
	ircAddr := os.Getenv("IRC_ADDR")
	irccon = irc.IRC(ircNickname, ircUsername)
	irccon.UseTLS = true
	irccon.Password = ircPassword
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.Connect(ircAddr)
	irccon.AddCallback("001", func(e *irc.Event) {
		for _, v := range channels {
			irccon.Join(v)
		}
	})
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		if match(event.Message()) {
			sendMail(mailFromname, mailFromaddr, mailToaddr, mailSubject, event.Message())
		}
	})
	go irccon.Loop()

	// 寝ないように定期的に自分を殴る
	go ping()
}

func main() {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, os.Getenv("IRC_MATCHER"))
	})
	router.GET("/ping", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "pong")
	})
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	n.Run(fmt.Sprintf("%s:%s", host, port))
}

func getChannels(str string) []string {
	r := []string{}
	for _, v := range strings.Split(str, ",") {
		r = append(r, fmt.Sprintf("#%s", strings.Trim(v, "#")))
	}
	return r
}

func sendMail(name, from, to, subject, body string) {
	mg.Send(mailgun.Message{
		FromName:    name,
		FromAddress: from,
		ToAddress:   to,
		Subject:     subject,
		Body:        body,
	})
}

func ping() {
	t := time.NewTicker(time.Minute * 50)
	for {
		select {
		case <-t.C:
			http.Get(os.Getenv("PING_URL"))
		}
	}
}

func match(str string) bool {
	str = strings.ToUpper(str)
	for _, v := range matcher {
		if 0 <= strings.Index(str, v) {
			return true
		}
	}
	return false
}
