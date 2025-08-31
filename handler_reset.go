package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.qry.ResetUsers(context.Background())
	if err != nil {
		fmt.Println("database NOT reset")
		return fmt.Errorf("handleReset: %w", err)
	}

	fmt.Println("database reset")

	return nil
}
