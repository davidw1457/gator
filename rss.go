package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: %w", err)
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: %w", err)
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: %w", err)
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: %w", err)
	}

	unescapeFeed(&rssFeed)

	return &rssFeed, nil
}

func unescapeFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, _ := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(
			feed.Channel.Item[i].Title,
		)
		feed.Channel.Item[i].Description = html.UnescapeString(
			feed.Channel.Item[i].Description,
		)
	}
}

func (feed RSSFeed) printFeed() {
	fmt.Printf("Feed title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed description: %s\n", feed.Channel.Description)
	fmt.Printf("Feed link: %s\n", feed.Channel.Link)
	for _, item := range feed.Channel.Item {
		item.printItem()
	}
}

func (item RSSItem) printItem() {
	fmt.Println("******************************")
	fmt.Printf("Item title: %s\n", item.Title)
	fmt.Printf("Item date: %s\n", item.PubDate)
	fmt.Printf("Item description: %s\n", item.Description)
	fmt.Printf("Item link: %s\n", item.Link)
}
