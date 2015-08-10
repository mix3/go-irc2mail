package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/sirsean/go-mailgun/mailgun"
	"github.com/thoj/go-ircevent"
)

var (
	// mailgun
	MAILGUN_APIKEY = GetEnvString("MAILGUN_APIKEY")
	MAILGUN_DOMAIN = GetEnvString("MAILGUN_DOMAIN")

	// irc
	IRC_CHANNELS = func() []string {
		r := []string{}
		for _, v := range strings.Split(GetEnvString("IRC_CHANNELS"), " ") {
			r = append(r, fmt.Sprintf("#%s", strings.Trim(v, "#")))
		}
		return r
	}()
	IRC_NICKNAME   = GetEnvString("IRC_NICKNAME")
	IRC_USERNAME   = GetEnvString("IRC_USERNAME")
	IRC_PASSWORD   = GetEnvString("IRC_PASSWORD")
	IRC_ADDR       = GetEnvString("IRC_ADDR")
	IRC_USE_SSL    = GetEnvBool("IRC_USE_SSL")
	IRC_DEBUG      = GetEnvBool("IRC_DEBUG")
	IRC_MATCH_LIST = func() []string {
		result := []string{}
		for _, v := range strings.Split(GetEnvString("IRC_MATCH_LIST"), " ") {
			result = append(result, strings.ToUpper(v))
		}
		return result
	}()

	// mail
	MAIL_FROM_NAME = GetEnvString("MAIL_FROM_NAME")
	MAIL_FROM_ADDR = GetEnvString("MAIL_FROM_ADDR")
	MAIL_TO_ADDR   = GetEnvString("MAIL_TO_ADDR")
	MAIL_SUBJECT   = GetEnvString("MAIL_SUBJECT")

	// server
	HOST        = os.Getenv("HOST")
	PORT        = os.Getenv("PORT")
	PING_MINUTE = GetEnvInt("PING_MINUTE")

	TARGET_NAME = GetEnvString("TARGET_NAME")
	AWAY_NAME   = GetEnvString("AWAY_NAME")
	APP_URL     = GetEnvString("APP_URL")
	DEBUG       = GetEnvBool("DEBUG")
	PING_PATH   = "/ping"
)

func debugf(f string, v ...interface{}) {
	if DEBUG {
		log.Printf(f, v...)
	}
}

func matcher(str string) bool {
	str = strings.ToUpper(str)
	for _, v := range IRC_MATCH_LIST {
		if 0 <= strings.Index(str, v) {
			return true
		}
	}
	return false
}

func GetEnvBool(key string) bool {
	v, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		log.Fatalf("env %s parse error: %v", key, err)
	}
	return v
}

func GetEnvString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env %s", key)
	}
	return v
}

func GetEnvInt(key string) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Fatalf("env %s parse error: %v", key, err)
	}
	return v
}

func NewIRCConn() *irc.Connection {
	// mailgun
	mg := mailgun.NewClient(MAILGUN_APIKEY, MAILGUN_DOMAIN)

	// ircevent
	irccon := irc.IRC(IRC_NICKNAME, IRC_USERNAME)
	irccon.Debug = IRC_DEBUG
	irccon.UseTLS = IRC_USE_SSL
	if irccon.UseTLS {
		irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	irccon.Password = IRC_PASSWORD
	irccon.Connect(IRC_ADDR)

	isAway := false

	irccon.AddCallback("001", func(e *irc.Event) {
		for _, v := range IRC_CHANNELS {
			irccon.Join(v)
		}
	})
	irccon.AddCallback("353", func(event *irc.Event) {
		for _, v := range strings.Split(event.Arguments[3], " ") {
			if v == AWAY_NAME {
				debugf("%s is away", TARGET_NAME)
				isAway = true
			}
		}
	})
	irccon.AddCallback("NICK", func(event *irc.Event) {
		if event.Nick == AWAY_NAME {
			debugf("%s is back", TARGET_NAME)
			isAway = false
		}
		for _, v := range event.Arguments {
			if v == AWAY_NAME {
				debugf("%s is away", TARGET_NAME)
				isAway = true
			}
		}
	})
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		if matcher(event.Message()) && isAway {
			debugf("send message")
			mg.Send(mailgun.Message{
				FromName:    MAIL_FROM_NAME,
				FromAddress: MAIL_FROM_ADDR,
				ToAddress:   MAIL_TO_ADDR,
				Subject:     MAIL_SUBJECT,
				Body:        event.Message(),
			})
		} else {
			debugf("skip send message")
		}
	})
	return irccon
}

func ping(url string) {
	t := time.NewTicker(time.Duration(PING_MINUTE) * time.Minute)
	for {
		select {
		case <-t.C:
			http.Get(url)
		}
	}
}

func main() {
	irccon := NewIRCConn()
	go func() {
		for {
			irccon.Loop()
		}
	}()
	go ping(APP_URL + PING_PATH)

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "%v", IRC_MATCH_LIST)
	})
	router.GET(PING_PATH, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "pong")
	})
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)
	n.Run(fmt.Sprintf("%s:%s", HOST, PORT))
}
