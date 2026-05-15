package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/imirazimi/graph/internal/infra/gin"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
	"fmt"
	"os"
	"time"
	"os/signal"
	"context"
	"syscall"

)

type Server struct {
	router   ginrouter.Router
	handler      Handler
}

func NewServer(router ginrouter.Router, handler Handler) Server {
	server := Server{
		router:		router,
		handler:	handler,
	}
	server.RegisterRoutes()
	return server
}

func (s *Server) RegisterRoutes() {
	
	s.router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })
	
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


    tasks := s.router.Group("/tasks")

    tasks.POST("", s.handler.Create)
    tasks.GET("", s.handler.List)
    tasks.GET("/:id", s.handler.GetByID)
    tasks.PUT("/:id", s.handler.Update)
    tasks.DELETE("/:id", s.handler.Delete)
}


func (s *Server) Serve() {
	httpserver := &http.Server{
		Addr: ":8080",
		Handler: s.router,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		fmt.Println("Starting server on port 8080...")
		if err := httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	fmt.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpserver.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped gracefully.")
}

