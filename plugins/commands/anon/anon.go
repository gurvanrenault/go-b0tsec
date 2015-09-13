package anon

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"
)

// Plugin is the anon plugin. Exposed as anon.Plugin.
type Plugin struct{}

// Help provides some help on the usage of the plugin.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    Allows to send anonymous messages on the channel where the bot is connected.")
	ib.Privmsgf(from, "    Example : /msg %s !command Hello everyone.", configuration.Config.BotName)
}

// Get actually sends the data to the channel.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		ib.Privmsg(configuration.Config.Channel, strings.Join(args, " "))
	}
}
