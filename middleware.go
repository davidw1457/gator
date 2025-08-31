package main

import (
	"context"
	"fmt"

	"github.com/davidw1457/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.qry.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("middlewareLoggedIn: %w", err)
		}

		err = handler(s, cmd, user)
		if err != nil {
			err = fmt.Errorf("middlewareLoggedIn: %w", err)
		}

		return err
	}
}
