package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	DB *gorm.DB
	Conn redis.Conn
)

func MysqlInit() {
	sql,err :=gorm.Open("mysql","root:weweixiao228@tcp(127.0.0.1:3306)/ongorm?charset=utf8&parseTime=true")
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	DB=sql
	if !DB.HasTable(User{}) {
		dd:=DB.CreateTable(User{})
		if dd.Error != nil { log.Println(dd.Error) }
	}
}

func RedisInit() {
	c,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	Conn = c

}


