package ioc

import (
	"ShareSphere/V0/internal/web"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitWebServer(mids []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mids...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddleWares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 不加这个会导致前端无法获取到返回的header
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 允许带cookie
		AllowCredentials: true,
		// 允许的域名
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 本地开发环境
				return true
			}
			return strings.Contains(origin, "lnu/suu.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
