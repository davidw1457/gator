package main

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"

    "github.com/davidw1457/gator/internal/database"
)

func handlerUsers(s *state, cmd command) error {
    users, err := s.qry.GetUsers(context.Background())
    if err != nil {
        return fmt.Errorf("handleUsers(%v, %v): %w", s, cmd, err)
    }

    for _, u := range users {
        if u.Name == s.cfg.CurrentUserName {
            fmt.Printf("* %s (current)\n", u.Name)
        } else {
            fmt.Printf("* %s\n", u.Name)
        }
    }

    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) < 1 {
        return fmt.Errorf("")
    }

    _, err := s.qry.GetUser(context.Background(), cmd.args[0])
    if err == nil {
        return fmt.Errorf("user %s already exists", cmd.args[0])
    }

    user, err := s.qry.CreateUser(
        context.Background(),
        database.CreateUserParams{
            ID: uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Name: cmd.args[0],
        },
    )
    if err != nil {
        return fmt.Errorf("")
    }

    err = s.cfg.SetUser(cmd.args[0])
    if err != nil {
        return fmt.Errorf("handerRegister(%v, %v): %w", s, cmd, err)
    }

    fmt.Printf("user %s created\n", cmd.args[0])
    fmt.Printf("%v\n", user)

    return nil
}

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
        return fmt.Errorf("handlerLogin(%v, %v): %w", s, cmd, err)
    }

    fmt.Printf("User set: %v\n", s.cfg)

    return nil
}
