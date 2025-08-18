package main

import (
    "fmt"
    "os"
    "github.com/davidw1457/gator/internal/config"
)

func main() {
    cfg, err := config.Read()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    err = cfg.SetUser("david")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    cfg, err = config.Read()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Printf("%v\n", cfg)
}
