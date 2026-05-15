package main

import (
    "github.com/imirazimi/graph/config"
    "github.com/imirazimi/graph/internal/infra/postgres"
    "github.com/imirazimi/graph/internal/infra/gin"
    "github.com/imirazimi/graph/internal/task"
)

func main() {
    cfg := config.LoadConfig()
    task.NewApp(
        ginrouter.NewRouter(cfg.AppPort),
        postgres.NewConnection(cfg.DatabaseURL()),
        cfg,
    ).Serve()
}