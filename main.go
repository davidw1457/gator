package main

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/lib/pq"
    "github.com/davidw1457/gator/internal/config"
    "github.com/davidw1457/gator/internal/database"
)

type state struct {
    cfg *config.Config
    qry *database.Queries
}

func main() {
    cfg, err := config.Read()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    db, err := sql.Open("postgres", cfg.DbUrl)
    if err != nil  {
        fmt.Println(err)
        os.Exit(1)
    }

    s := state{cfg: &cfg, qry: database.New(db)}
    cmds := commands{
        registeredCommands: make(map[string]func(*state, command) error),
    }

    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)
    cmds.register("reset", handlerReset)
    cmds.register("users", handlerUsers)

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
