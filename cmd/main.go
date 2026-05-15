package main

import (
    "log"

    "github.com/imirazimi/graph/config"
    apphttp "github.com/imirazimi/graph/internal/task/handler"
    "github.com/imirazimi/graph/internal/infra/postgres"
    "fmt"
)

func main() {
    cfg := config.LoadConfig()

    db := postgres.NewPostgresConnection(cfg.DatabaseURL())
    defer db.Close()

    router := apphttp.SetupRouter()

    log.Printf("server started on port %s", cfg.AppPort)

    if err := router.Run(":" + cfg.AppPort); err != nil {
        log.Fatal(err)
    }
}