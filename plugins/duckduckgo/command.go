package duckduckgo

import (
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	apiURL  = "https://api.duckduckgo.com/?q=%s&format=json"
	command = "ddg"
)

type message struct {
	Definition       string
	DefinitionSource string
	Heading          string
	AbstractText     string
	Abstract         string
	AbstractSource   string
	Image            string
	Type             string
	AnswerType       string
	Redirect         string
	DefinitionURL    string
	Answer           string
	AbstractURL      string
	Results          []relatedTopic
	RelatedTopics    []relatedTopic
}

type relatedTopic struct {
	Result string
	Icon   struct {
		URL    string
		Height interface{}
		Width  interface{}
	}
	FirstURL string
	Text     string
}

// Command is the duckduckgo plugin.
type Command struct {
	Started bool
}

func init() {
	plugins.Commands[command] = new(Command)
}

// Help provides some help on the plugin
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Search directly on DuckDuckGo.")
	ib.Privmsg(from, "Example : !ddg Who is James Cameron ?")
}

// Get actually sends the data to the channel
func (c *Command) Get(ircbot *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		res, err := c.fetch(strings.Join(args, " "))
		if err != nil || res == "" {
			if err != nil {
				log.Println(err)
			}
			return
		}
		ircbot.Privmsg(to, res)
	}
}

// Start starts the plugin and returns any occurred error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(command, configuration.Config.Commands) {
		c.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occurred error, nil otherwise
func (c *Command) Stop() error {
	c.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (c *Command) IsStarted() bool {
	return c.Started
}

func (c *Command) fetch(query string) (string, error) {
	var t message
	url := utils.EncodeURL(apiURL, query)
	if err := utils.FetchURL(url, &t); err != nil {
		return "", err
	}
	return t.Abstract, nil
}
