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

type YouTubeLink struct {
	Rel      string
	Href     string
	MimeType string `json:"type"`
}

type YouTubeFeed struct {
	Feed struct {
		Entry []struct {
			Link []YouTubeLink
		}
	}
}

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
	var data YouTubeFeed
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	// Pick a random video
	videos := data.Feed.Entry
	video := videos[rand.Intn(len(videos))]
	links := video.Link

	for _, link := range links {
		if link.Rel == "alternate" && link.MimeType == "text/html" {
			client.Privmsg(buffer, link.Href)
		}
	}
}
