package main

import (
    "context"
    "fmt"
)

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) < 1 {
        return fmt.Errorf("No username provided")
    }
    
    _, err := s.qry.GetUser(context.Background(), cmd.args[0])
    if err != nil {
        return fmt.Errorf("handlerLogin(%v, %v): %w", s, cmd, err)
    }

    err = s.cfg.SetUser(cmd.args[0])
    if err != nil {
        return fmt.Errorf(
            "handlerLogin(%v, %v): %w",
            s,
            cmd,
            err,
        )
    }

    fmt.Printf("User set: %v\n", s.cfg)

    return nil
}

