package main

import (
    "github.com/mykykh/concerts-api/internal/api"
)

func main() {
    api := api.Init()

    api.Run()
}
