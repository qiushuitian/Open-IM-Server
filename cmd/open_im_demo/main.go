package main

import (
	"Open_IM/internal/demo/bottle"
	"Open_IM/internal/demo/register"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(utils.CorsHandler())

	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/code", register.SendVerificationCode)
		authRouterGroup.POST("/verify", register.Verify)
		authRouterGroup.POST("/password", register.SetPassword)
		authRouterGroup.POST("/login", register.Login)
		authRouterGroup.POST("/register", register.Register)
	}

	bottleGroup := r.Group("/bottle")
	{
		bottleGroup.POST("/throw", bottle.ThrowBottle)
		bottleGroup.POST("/pick", bottle.PickBottle)
		bottleGroup.POST("/check_upgrade", bottle.CheckUpgrade)
		bottleGroup.POST("/terms", bottle.Terms)
		bottleGroup.POST("/report", bottle.Report)
		bottleGroup.POST("/feedback", bottle.Feedback)

	}

	log.NewPrivateLog("demo")
	ginPort := flag.Int("port", 42233, "get ginServerPort from cmd,default 42233 as port")
	flag.Parse()
	r.Run(utils.ServerIP + ":" + strconv.Itoa(*ginPort))
}
