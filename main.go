package main

import (
    "fmt"
    "os"
    "github.com/davidw1457/gator/internal/config"
)

type state struct {
    cfg *config.Config
}

func main() {
    cfg, err := config.Read()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    s := state{cfg: &cfg}
    cmds := commands{
        registeredCommands: make(map[string]func(*state, command) error),
    }

    cmds.register("login", handlerLogin)

    if len(os.Args) < 2 {
        fmt.Println("Usage: gator <command> [args...]")
        os.Exit(1)
    }

    err = cmds.run(&s, command{name: os.Args[1], args: os.Args[2:]})

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
