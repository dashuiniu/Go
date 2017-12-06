// Msg project main.go
package main

import (
	"Msg/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) //	生产模式
	route := gin.Default()

	route.GET("/", func(c *gin.Context) {
		c.String(200, "home")
	})

	//增加消息
	go route.POST("/addMsg", controllers.AddMsg)
	//消息列表
	go route.GET("/msgList/:uid/:mtype/:type/:page", controllers.MsgList)
	//消息
	go route.GET("/msgInfo/:uid/:mid", controllers.MsgInfo)

	route.Run(":8099")
}
