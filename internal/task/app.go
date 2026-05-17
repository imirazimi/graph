package task

import (
	"github.com/imirazimi/graph/internal/infra/gin"
	"github.com/imirazimi/graph/internal/infra/redis"
	"github.com/imirazimi/graph/internal/infra/postgres"
	"github.com/imirazimi/graph/internal/task/service"
	"github.com/imirazimi/graph/internal/task/repository"
	"github.com/imirazimi/graph/internal/task/handler/http"
	"context"
	"time"
	"fmt"
	"os/signal"
	"net/http"
	"os"
	"syscall"

	"github.com/imirazimi/graph/internal/infra/tracing"
	"github.com/imirazimi/graph/config"
)

type App struct {
	Server  handler.Server
}

func ServeApp(
	router ginrouter.Router,
	postgres postgres.Connection,
	redis redis.RedisClient,
	cfg config.Config,
) App {

	app := App{
		Server: handler.NewServer(
			router,
			handler.NewHandler(
				service.NewService(
					repository.NewCacheRepository(
						redis,
						repository.NewRepository(postgres),
					),
				),
			),
		),
	}

	// ---- Tracing ----
	tracerShutdown := tracing.InitTracer()

	// ---- HTTP Server ----
	httpServer := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      app.Server.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ---- Run Server ----
	go func() {
		fmt.Printf("🚀 Server started on port %s\n", cfg.AppPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ server error: %v\n", err)
		}
	}()

	// ---- OS Signal ----
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Shutting down server gracefully...")

	// ---- Graceful Shutdown Context ----
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ---- HTTP Shutdown ----
	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("❌ HTTP shutdown failed: %v\n", err)
	}

	// ---- Close external resources ----
	postgres.Close()

	if err := redis.Close(); err != nil {
		fmt.Printf("❌ Redis close error: %v\n", err)
	}

	// ---- Tracer Shutdown ----
	if err := tracerShutdown(ctx); err != nil {
		fmt.Printf("❌ Tracer shutdown error: %v\n", err)
	}

	fmt.Println("✅ Server stopped gracefully")

	return app
}
