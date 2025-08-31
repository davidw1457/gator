package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/davidw1457/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(
		context.Background(),
		"https://www.wagslane.dev/index.xml",
	)
	if err != nil {
		return fmt.Errorf("handlerAgg: %w", err)
	}

	feed.printFeed()

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("proper usage: gator addfeed <name> <url>")
	}

	feed, err := s.qry.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("handlerAddFeed: %w", err)
	}

	_, err = s.qry.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("handlerAddFeed: %w", err)
	}

	fmt.Println("Feed added")
	fmt.Println("******************************")
	printFeed(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.qry.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("handlerFeeds: %w", err)
	}

	printFeeds(feeds)

	return nil
}

func printFeeds(feeds []database.GetFeedsRow) {
	if len(feeds) == 0 {
		fmt.Println("no feeds added")
		return
	} else {
		fmt.Printf("%d feeds found:\n", len(feeds))
	}

	for _, f := range feeds {
		fmt.Println("******************************")
		printFeed(database.Feed{Name: f.Name, Url: f.Url})
		fmt.Printf("Feed user: %s\n", f.UserName)
	}
}

func printFeed(f database.Feed) {
	fmt.Printf("Feed name: %s\n", f.Name)
	fmt.Printf("Feed url: %s\n", f.Url)
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("proper usage: gator follow <url>")
	}

	feed, err := s.qry.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("handlerFollow: %w", err)
	}

	feedFollow, err := s.qry.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("handlerFollow: %w", err)
	}

	fmt.Printf("Feed name: %s\n", feedFollow.FeedName)
	fmt.Printf("Followed by: %s\n", feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.qry.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("handlerFollowing: %w", err)
	}

	fmt.Printf("%s is following:\n", s.cfg.CurrentUserName)
	for _, feed := range feedFollows {
		fmt.Println(feed.FeedName)
	}

	return nil
}
