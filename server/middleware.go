package server

import (
	"log"
	"srv/errtrace"

	"github.com/gin-gonic/gin"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				err := errtrace.WrapRecoveredError(err)
				log.Print("recover: ", err.String())
			}
		}()
		c.Next()
	}
}
