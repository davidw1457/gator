package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/davidw1457/gator/internal/database"
	"github.com/google/uuid"
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

func (feed RSSFeed) print() {
	fmt.Printf("Feed title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed description: %s\n", feed.Channel.Description)
	fmt.Printf("Feed link: %s\n", feed.Channel.Link)
	for _, item := range feed.Channel.Item {
		item.print()
	}
}

func (item RSSItem) print() {
	fmt.Println("******************************")
	fmt.Printf("Item title: %s\n", item.Title)
	// fmt.Printf("Item date: %s\n", item.PubDate)
	// fmt.Printf("Item description: %s\n", item.Description)
	// fmt.Printf("Item link: %s\n", item.Link)
}

func scrapeFeeds(ctx context.Context, s *state) error {
	feed, err := s.qry.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("scrapeFeeds: %w", err)
	}

	err = s.qry.MarkFeedFetched(
		ctx,
		database.MarkFeedFetchedParams{
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: time.Now(),
			ID:        feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("scrapeFeeds: %w", err)
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("scrapeFeeds: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		var publishedDate time.Time
		publishedDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("Unsupported date format: %s\n", item.PubDate)
			continue
		}
		_, err = s.qry.CreatePost(
			ctx,
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				PublishedAt: publishedDate,
				FeedID:      feed.ID,
			},
		)
		if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
			continue
		} else if err != nil {
			fmt.Printf("Error inserting: %v\n", err)
		}
	}

	return nil
}
