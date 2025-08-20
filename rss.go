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
        Title string `xml:"title"`
        Link string `xml:"link"`
        Description string `xml:"description"`
        Item []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type RSSItem struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    client := http.Client{Timeout: 10 * time.Second}

    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        return nil, fmt.Errorf("fetchFeed(%v, %v): %w", ctx, feedURL, err)
    }

    req.Header.Set("User-Agent", "gator")
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetchFeed(%v, %v): %w", ctx, feedURL, err)
    }
    defer resp.Body.Close()

    dat, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("fetchFeed(%v, %v): %w", ctx, feedURL, err)
    }

    rssFeed := RSSFeed{}
    err = xml.Unmarshal(dat, &rssFeed)
    if err != nil {
        return nil, fmt.Errorf("fetchFeed(%v, %v): %w", ctx, feedURL, err)
    }

    unescapeFeed(&rssFeed)

    return &rssFeed, nil
}

func unescapeFeed(f *RSSFeed) {
    f.Channel.Title = html.UnescapeString(f.Channel.Title)
    f.Channel.Description = html.UnescapeString(f.Channel.Description)
    for i, _ := range f.Channel.Item {
        f.Channel.Item[i].Title = html.UnescapeString(
            f.Channel.Item[i].Title,
        )
        f.Channel.Item[i].Description = html.UnescapeString(
            f.Channel.Item[i].Description,
        )
    }
}
