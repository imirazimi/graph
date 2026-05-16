package task

import (
	"github.com/imirazimi/graph/internal/infra/gin"
	"github.com/imirazimi/graph/internal/infra/postgres"
	"github.com/imirazimi/graph/internal/task/service"
	"github.com/imirazimi/graph/internal/task/repository"
	"github.com/imirazimi/graph/internal/task/handler/http"
	"context"
	"github.com/imirazimi/graph/internal/infra/tracing"
	"github.com/imirazimi/graph/config"
)

type App struct {
	server  handler.Server
}


func NewApp(router ginrouter.Router,postgres postgres.Connection,cfg config.Config) App {
	return App {
		handler.NewServer(
			router,
			handler.NewHandler(
				service.NewService(
					repository.NewRepository(
						postgres,
					),
				),
			),
		),
	}
}

func (a App) Serve () {
	shutdown := tracing.InitTracer()
    defer shutdown(context.Background())
	
	a.server.Serve()
}
