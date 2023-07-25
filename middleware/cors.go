package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowHeaders = []string{"*"}
	corsCfg.AllowCredentials = true
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowWildcard = true
	return cors.New(corsCfg)
}
