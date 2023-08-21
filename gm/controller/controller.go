package controller

import (
	"net/http"
	"strconv"
	"sync"

	"hk4e/common/config"
	"hk4e/common/rpc"
	"hk4e/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	gmClientMap     map[uint32]*rpc.GMClient
	gmClientMapLock sync.RWMutex
}

func NewController() (r *Controller) {
	r = new(Controller)
	r.gmClientMap = make(map[uint32]*rpc.GMClient)
	go r.registerRouter()
	return r
}

func (c *Controller) authorize() gin.HandlerFunc {
	return func(context *gin.Context) {
		if true {
			// 验证通过
			context.Next()
			return
		}
		// 验证不通过
		context.Abort()
		context.JSON(http.StatusOK, gin.H{
			"code": "10001",
			"msg":  "没有访问权限",
		})
	}
}

func (c *Controller) registerRouter() {
	if config.GetConfig().Logger.Level == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.Use(c.authorize())
	engine.POST("/gm/cmd", c.gmCmd)
	port := config.GetConfig().HttpPort
	addr := ":" + strconv.Itoa(int(port))
	err := engine.Run(addr)
	if err != nil {
		logger.Error("gin run error: %v", err)
	}
}
