package handlers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"github.com/mikeclarke/go-irclib"
	"github.com/garyburd/redigo/redis"
)

var lookupRE = regexp.MustCompile(`^\?\?(\S+)`)
var commandRE = regexp.MustCompile(`^\?(learn|forget)\s.+`)

var commandMap = map[string] string {
	"?learn": "SADD",
	"?forget": "SREM",
}

type DocBot struct {
	address  string
	password string
	conn     redis.Conn
}

func (bot *DocBot) connect() error {
	c, err := redis.Dial("tcp", bot.address)
	if err != nil {
		log.Printf("error connecting: %v", err)
	}

	_, err = c.Do("AUTH", bot.password);

	if err != nil {
		log.Printf("Unable to AUTH: %v", err)
		c.Close()
		return err
	}

	bot.conn = c
	return nil
}

func (bot *DocBot) HandleEvent(event *irc.Event) {
	// Parse the message
	if event.Command != "PRIVMSG" || len(event.Arguments) < 2 {
		return
	}

	client := event.Client
	buffer := event.Arguments[0]

	// Create redis connection
	bot.connect()
	defer bot.conn.Close()

	switch {
	case lookupRE.MatchString(event.Arguments[1]):

		// Grab lookup value
		key := lookupRE.FindString(event.Arguments[1])
		key = key[2:]

		// Check redis set
		results, _ := redis.Values(bot.conn.Do("SMEMBERS", key))

		// Check if no results
		if len(results) > 0 {
			for _, url := range results {
				client.Privmsg(buffer, fmt.Sprintf(":: %s", url))
			}
		} else {
			client.Privmsg(buffer, "no results found!")
		}

	case commandRE.MatchString(event.Arguments[1]):
		args := strings.Split(event.Arguments[1], " ")
		cmd, url, keys := args[0], args[1], args[2:]

		for _, key := range keys {
			bot.conn.Do(commandMap[cmd], key, url)
		}
		client.Privmsg(buffer, "command was successful.")

	default:
		return
	}
}

func NewDocBot(addr string, pass string) *DocBot {
	bot := &DocBot{
		address: addr,
		password: pass,
	}

	return bot
}
