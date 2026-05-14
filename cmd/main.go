package main

import (
    "log"

    "github.com/imirazimi/graph/config"
    apphttp "github.com/imirazimi/graph/internal/task/handler"
)

func main() {
    cfg := config.LoadConfig()
    router := apphttp.SetupRouter()

    log.Printf("server started on port %s", cfg.AppPort)

    if err := router.Run(":" + cfg.AppPort); err != nil {
        log.Fatal(err)
    }
}