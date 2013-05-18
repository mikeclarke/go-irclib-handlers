package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"github.com/mikeclarke/go-irclib"
)

var message_RE = regexp.MustCompile(`(?:youtube|yt)(?: me)? (.*)`)

func YouTube(event *irc.Event) {
	matches := message_RE.FindStringSubmatch(event.Raw)

	if len(matches) < 2 {
		return
	}

	client := event.Client
	buffer := event.Arguments[0]

	query := matches[1]
	baseUrl := "http://gdata.youtube.com/feeds/api/videos"

	params := url.Values{}
	params.Add("orderBy", "relevance")
	params.Add("max-results", "15")
	params.Add("alt", "json")
	params.Add("q", query)

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseUrl, params.Encode()))

	if err != nil {
		client.Privmsg(buffer, fmt.Sprintf("error hitting youtube: %s", err))
		return
	}

	// Decode response
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	// Pick a random video
	feed := data["feed"].(map[string]interface{})
	videos := feed["entry"].([]interface{})
	video := videos[rand.Intn(len(videos))].(map[string]interface{})
	links := video["link"].([]interface{})

	for _, item := range links {
		link := item.(map[string]interface{})
		if link["rel"] == "alternate" && link["type"] == "text/html" {
			client.Privmsg(buffer, link["href"].(string))
		}
	}
}
