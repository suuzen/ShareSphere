//go:build wireinject

package main

import (
	"ShareSphere/V0/internal/repository"
	"ShareSphere/V0/internal/repository/cache"
	"ShareSphere/V0/internal/repository/dao"
	"ShareSphere/V0/internal/service"
	"ShareSphere/V0/internal/service/sms"
	"ShareSphere/V0/internal/web"
	"ShareSphere/V0/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		dao.NewGORMUserDao,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		sms.NewSmsService,
		web.NewUserHandler,
		ioc.InitMiddleWares,
		ioc.InitWebServer,
	)
	return &gin.Engine{}

}
