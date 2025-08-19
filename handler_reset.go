package main

import (
    "context"
    "fmt"
)

func handlerReset(s *state, cmd command) error {
    err := s.qry.ResetUsers(context.Background())
    if err != nil {
        fmt.Println("user table NOT reset")
        return fmt.Errorf("handleReset(%v, %v): %w", s, cmd, err)
    }

    fmt.Println("user table reset")

    return nil
}
