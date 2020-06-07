package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
	Conn redis.Conn
)

func MysqlInit(){
	sql,err :=gorm.Open("mysql","root:@tcp(127.0.0.1:3306)/ongorm?charset=utf8&parseTime=true")
	if err!=nil {
		fmt.Println(err.Error())
	}
	DB=sql
}

func RedisInit(){
	c,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	Conn = c

}


