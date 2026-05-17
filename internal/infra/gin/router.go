package ginrouter

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(port string) Router {
	return Router{
		gin.Default(),
	}
}
