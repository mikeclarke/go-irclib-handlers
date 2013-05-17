package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"github.com/mikeclarke/go-irclib"
)

var message_RE = regexp.MustCompile(`(?:youtube|yt)(?: me)? (.*)`)

func Youtube(event *irc.Event) {
	matches := message_RE.FindStringSubmatch(event.Raw)

	if len(matches) > 1 {
		client := event.Client
		query := matches[1]
		baseUrl := "http://gdata.youtube.com/feeds/api/videos"

		params := url.Values{}
		params.Add("orderBy", "relevance")
		params.Add("max-results", "15")
		params.Add("alt", "json")
		params.Add("q", query)

		resp, err := http.Get(fmt.Sprintf("%s?%s", baseUrl, params.Encode()))

		if err != nil {
			client.SendRawf("error hitting youtube: %s", err)
			return
		}

		// Decode response
		decoder := json.NewDecoder(resp.Body)
		log.Print(decoder)
	}
}
