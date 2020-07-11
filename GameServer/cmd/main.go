package main

import (
	"Server/api"
	"Server/model"
	ws "Server/websocket"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)



func main(){
	go ws.WebsocketManager.Start()

	model.MysqlInit()
	//model.RedisInit()

	r := gin.Default()
	api.SetRouter(r)

	r.Run(":8080")


}