package main

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	feed := RSSFeed{}

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return feed, err
	}
	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return feed, err
	}

	return feed, nil
}
