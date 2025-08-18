package main

import "fmt"

type command struct {
    name string
    args []string
}

type commands struct {
    registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
    f, ok := c.registeredCommands[cmd.name]
    if !ok {
        return fmt.Errorf("%s is not a valid command", cmd.name)
    }

    err := f(s, cmd)
    if err != nil {
        return fmt.Errorf("commands.run(%v, %v): %w", s, cmd, err)
    }

    return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.registeredCommands[name] = f
}
