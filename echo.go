package handlers

import (
	"github.com/mikeclarke/go-irclib"
	"log"
)

func Echo(event *irc.Event) {
	log.Printf("<-- %v, %v, %v", event.Prefix, event.Command, event.Arguments)
}
