package main

import (
    "context"
    "fmt"
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
