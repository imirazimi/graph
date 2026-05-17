package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/imirazimi/graph/internal/infra/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http/pprof"
    ginSwagger "github.com/swaggo/gin-swagger"

)

type Server struct {
	Router   ginrouter.Router
	Handler      Handler
}

func NewServer(router ginrouter.Router, handler Handler) Server {
	server := Server{
		Router:		router,
		Handler:	handler,
	}
	server.RegisterRoutes()
	return server
}

func (s *Server) RegisterRoutes() {
	
	
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
	pprofGroup := s.Router.Group("/debug/pprof")

    pprofGroup.GET("/", gin.WrapF(pprof.Index))
    pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
    pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
    pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
	
	s.Router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })
	
	s.Router.Use(otelgin.Middleware("task-manager"))
	
	s.Router.Use(MetricMiddleware())

	s.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))

    tasks := s.Router.Group("/tasks")

    tasks.POST("", s.Handler.Create)
    tasks.GET("", s.Handler.List)
    tasks.GET("/:id", s.Handler.GetByID)
    tasks.PUT("/:id", s.Handler.Update)
    tasks.DELETE("/:id", s.Handler.Delete)
}