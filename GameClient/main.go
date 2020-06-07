package main

import (
	"client/Process"
	"fmt"
	"time"
)

func main() {
	Process.Sigin()
	Process.Play()
	fmt.Println("十秒后自动关闭客户端")
	time.Sleep(time.Second*10)
}


