package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Atom    string   `xml:"atom,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func fetchFeed(ctx context.Context, feedURL string) (*Rss, error) {
	var RssResponse *Rss
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Printf("ERROR MAKING GET REQUEST:%v\n", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("ERROR GETTING HTTP RESPONSE:%v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR COULD NOT READ HTTP RESPONSE:%v\n", err)
	}
	err = xml.Unmarshal(body, &RssResponse)
	if err != nil {
		fmt.Printf("ERROR UNMARSHALLING XML:%v\n", err)
		return nil, err
	}
	RssResponse.Channel.Title = html.UnescapeString(RssResponse.Channel.Title)
	RssResponse.Channel.Description = html.UnescapeString(RssResponse.Channel.Description)
	for _, i := range RssResponse.Channel.Item {
		i.Title = html.UnescapeString(i.Title)
		i.Description = html.UnescapeString(i.Description)
	}
	return RssResponse, nil
}
