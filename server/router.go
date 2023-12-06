package server

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	e := gin.Default()

	e.Use(recovery())

	e.NoRoute(fileServe())

	_ = e.Group("/", fallback(false))

	return e
}
