package main

import (
    "github.com/imirazimi/graph/config"
    "github.com/imirazimi/graph/internal/infra/postgres"
    "github.com/imirazimi/graph/internal/infra/gin"
    "github.com/imirazimi/graph/internal/task"
    _ "github.com/imirazimi/graph/docs"
)

// @title Task Manager API
// @version 1.0
// @description Interview Task Manager Service
// @host localhost:8080
// @BasePath /

func main() {
    cfg := config.LoadConfig()
    task.NewApp(
        ginrouter.NewRouter(cfg.AppPort),
        postgres.NewConnection(cfg.DatabaseURL()),
        cfg,
    ).Serve()
}