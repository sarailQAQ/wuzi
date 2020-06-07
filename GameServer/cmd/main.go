package main

import (
	ws "Server/Websocket"
	"Server/api"
	"Server/model"
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main(){
	go ws.WebsocketManager.Start()

	//go ws.TestSendGroup()
	//go ws.TestSendAll()

	model.MysqlInit()
	//model.RedisInit()

	r := gin.Default()
	api.SetRouter(r)

	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Start Error: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}
	log.Println("Server Shutdown")

}



func NewData(){
	dd:=model.DB.CreateTable(model.User{})
	fmt.Println(dd.Error)
	dd = model.DB.CreateTable(model.Chat{})
	fmt.Println(dd.Error)
}