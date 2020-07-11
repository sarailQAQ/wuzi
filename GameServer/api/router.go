package api

import (
	"Server/middle_ware"
	ws "Server/websocket"
	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine){
	r.POST("/login",Login)
	r.POST("/register",Register)
	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/:channel", middle_ware.LoginStatus, ws.WebsocketManager.WsClient) //
	}

}
