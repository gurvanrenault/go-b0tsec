package main

import (
	"crypto/tls"
	"log"
	"strings"

	flag "github.com/ogier/pflag"
	"github.com/thoj/go-ircevent"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/plugins/cleverbot"
	_ "github.com/depado/go-b0tsec/pluginsinit"
)

func main() {
	var err error

	cnf := configuration.Config

	// Argument parsing
	configuration.ConfPath = flag.StringP("config", "c", "conf.yml", "Local path to configuration file.")
	flag.Parse()

	// Load the configuration of the bot
	configuration.Load()

	if err = cleverbot.Clever.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Storage initialization and defering Close
	if err := database.BotStorage.Open(); err != nil {
		log.Fatalf("something went wrong with the databse : %v", err)
	}
	defer database.BotStorage.Close()

	// Plugins initialization and defering Stop
	plugins.Start()
	defer plugins.Stop()

	// Bot initialization
	ib := irc.IRC(cnf.BotName, cnf.BotName)
	if cnf.TLS {
		ib.UseTLS = true
		if cnf.InsecureTLS {
			ib.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
	if err = ib.Connect(cnf.Server); err != nil {
		log.Fatal(err)
	}

	// Callback on 'Connected' event
	ib.AddCallback("001", func(e *irc.Event) {
		ib.Join(cnf.Channel)
	})

	// Callback on 'Invite' event
	ib.AddCallback("INVITE", func(e *irc.Event) {
		ib.Join(cnf.Channel)
	})

	// Callback on 'Message' event
	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		from := e.Nick
		to := e.Arguments[0]
		m := e.Message()

		for _, c := range plugins.Middlewares {
			c.Get(ib, from, to, m)
		}

		if strings.HasPrefix(m, cnf.CommandCharacter) {
			if len(m) > 1 {
				splitted := strings.Fields(m[1:])
				command := splitted[0]
				args := splitted[1:]
				if p, ok := plugins.Commands[command]; ok {
					p.Get(ib, from, to, args)
				}
			}
		}
	})
	ib.Loop()
}
