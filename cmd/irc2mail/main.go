package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/mix3/go-irc2mail"
)

func main() {
	var opts struct {
		ConfigPath string `short:"c" long:"config" default:"config.yaml"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
	}
	c, err := irc2mail.LoadConfig(opts.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := irc2mail.NewBot(c)
	if err != nil {
		log.Fatal(err)
	}
	bot.Run()
}
