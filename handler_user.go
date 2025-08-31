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
		return fmt.Errorf("handleUsers: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("proper usage: gator register <name>")
	}

	_, err := s.qry.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return fmt.Errorf("user %s already exists", cmd.args[0])
	}

	user, err := s.qry.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("handlerRegister: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("handerRegister: %w", err)
	}

	fmt.Printf("user %s created\n", cmd.args[0])
	printUser(user)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username provided")
	}

	_, err := s.qry.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("handlerLogin: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("handlerLogin: %w", err)
	}

	fmt.Printf("User set: %s\n", s.cfg.CurrentUserName)

	return nil
}

func printUser(user database.User) {
	fmt.Println("******************************")
	fmt.Printf("User id: %v\n", user.ID)
	fmt.Printf("User name: %s\n", user.Name)
}
