package main

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"

    "github.com/davidw1457/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
    f, err := fetchFeed(
        context.Background(), 
        "https://www.wagslane.dev/index.xml",
    )
    if err != nil {
        return fmt.Errorf("handlerAgg(%v, %v): %w", s, cmd, err)
    }

    fmt.Println(f)

    return nil
}

func handlerAddFeed(s *state, cmd command) error {
    if len(cmd.args) < 2 {
        return fmt.Errorf("proper usage: gator addfeed <name> <url>")
    }

    u, err := s.qry.GetUser(context.Background(), s.cfg.CurrentUserName)
    if err != nil {
        return fmt.Errorf("handlerAddFeed(%v, %v): %w", s, cmd, err)
    }

    f, err := s.qry.CreateFeed(
        context.Background(),
        database.CreateFeedParams{
            ID: uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Name: cmd.args[0],
            Url: cmd.args[1],
            UserID: u.ID,
        },
    )
    if err != nil {
        return fmt.Errorf("handlerAddFeed(%v, %v): %w", s, cmd, err)
    }

    fmt.Println(f)

    return nil
}
