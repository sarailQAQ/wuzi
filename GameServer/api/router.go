package api

import (
	"Server/Middleware"
	ws "Server/Websocket"
	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine){
	r.POST("/login",Login)
	r.POST("/register",Register)
	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/:channel",Middleware.LoginStatus, ws.WebsocketManager.WsClient)//
	}

}
