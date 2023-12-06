package server

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(recovery())
	e.NoRoute(fileServe())
	_ = e.Group("/", fallback(false))
	return e
}
