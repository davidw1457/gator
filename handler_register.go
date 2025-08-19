package main

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"

    "github.com/davidw1457/gator/internal/database"
)

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

    s.cfg.SetUser(cmd.args[0])
    fmt.Printf("user %s created\n", cmd.args[0])
    fmt.Printf("%v\n", user)

    return nil
}
