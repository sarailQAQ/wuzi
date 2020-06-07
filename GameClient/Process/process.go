package Process

import (
	"bytes"
	"client/Struct"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const N = 20
var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var token,user string
var a [N+1][N+1]int
var lock sync.Mutex

func Sigin() {
	var u Struct.LoginForm
	var address string
	fmt.Println("请输入服务器地址")
	fmt.Scanln(&address)
	addr = flag.String("addr",address,"http service address")
	for  {
		fmt.Println("请输入用户名密码")
		fmt.Scanln(&u.Username,&u.Password)
		data,err := json.Marshal(u)
		if err != nil {
			fmt.Println("登陆时发生意外错误 ",err)
			continue
		}
		reader := bytes.NewReader(data)
		url := address+"/login"
		req,err := http.NewRequest("POST",url,reader)
		if err != nil {
			fmt.Println("登陆时发生意外错误1",err)
			continue
		}
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("登陆时发生意外错误2",err)
			continue
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("登陆时发生意外错误3",err)
			continue
		}

		var respJson gin.H

		json.Unmarshal(respBytes,&respJson)
		respData := respJson["data"]
		if r,ok := respData.(map[string]interface{});ok {
			temp := r["token"]
			if str,ok := temp.(string);ok {token = str}
			temp = r["username"]
			if str,ok := temp.(string);ok {user = str}
		}
		break
	}
	fmt.Println("Hello,",user)
}

func Play() {
	var dialer *websocket.Dialer
	var group string
	fmt.Println("请输入房间号：")
	fmt.Scanln(&group)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/"+group}
	var header http.Header = make(map[string][]string)
	//cheader.Set("token",token)
	header.Add("token",token)
	conn,_,err := dialer.Dial(u.String(),header)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//go timeWriter(conn)


	//游戏准备阶段
	status := 0
	for {
		if status == 0 {
			fmt.Println("请输入1准备游戏")
			temp := 0
			for temp != 1 {
				fmt.Scanln(&temp)
			}
			status = (status + 1) % 2

		}

		data := Struct.TransData{
			Type:"ready",
		}
		err := SendWS(conn,data)
		if err != nil {
			fmt.Println(err)
			return
		}
		break
	}
	Read(conn)
	fmt.Println("ready")

	//游戏开始
	var data interface{}
	err = conn.ReadJSON(&data)
	if err != nil {
	 	fmt.Println("err1:",err)
	 	return
	}
	if dataMap,ok := data.(map[string]interface{});ok {
		status,_ = strconv.Atoi(Struct.TypeString(dataMap["type"]))
		status -= 1000
		fmt.Println(Struct.TypeString(dataMap["message"]))
	}
	var opt Struct.OptData
	fmt.Println("start",status)
	BoardInit()

	for {
		PrintBoard()
		if status == 1 {
			for {
				fmt.Println("请输入你要落子的坐标x,y")
				var px, py int
				fmt.Scanln(&px, &py)
				if a[px][py] != 0 || px <= 0 || py<=0 || px>N || py>N{
					continue
				}
				opt = Struct.OptData{
					Type: "play",
					User: user,
					Px:   strconv.Itoa(px),
					Py:   strconv.Itoa(py),
				}
				ok := Judge(px,py)
				if ok { opt.Type = "over" }
				err := SendWS(conn, opt)
				if err != nil {
					fmt.Println("err",err)
					return
				}
				a[px][py] = 2
				break
			}
		}

		err := conn.ReadJSON(&data)
		if err != nil {

			break
		}
		if dataMap,ok := data.(map[string]interface{});ok {
			opt.Type = Struct.TypeString(dataMap["type"])
			opt.User = Struct.TypeString(dataMap["user"])
			opt.Px = Struct.TypeString(dataMap["px"])
			opt.Py = Struct.TypeString(dataMap["py"])
		}else {break}

		px,_ := strconv.Atoi(opt.Px); py,_ := strconv.Atoi(opt.Py)
		if a[px][py] != 2 {a[px][py] = 1}
		if opt.Type == "over" {
			if opt.User == user {
				fmt.Println("You win! 五子，行")
			} else {fmt.Println("You lose! 没有五子，不行" )}
			break
		}
		status = (status + 1) % 2
	}
}

func Read(conn *websocket.Conn) {
	for {
		var data interface{}
		err := conn.ReadJSON(&data)
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		if dataMap,ok := data.(map[string]interface{}); ok {
			fmt.Println(dataMap["message"])
			if dataMap["message"] == "Game Start!" {break}
		}

	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}

func PrintBoard() {
	//cmd := exec.Command("cmd", "/c", "cls")
	//cmd.Stdout = os.Stdout
	//cmd.Run()
	var ch string
	for i := 0; i <= N; i++ {
		fmt.Printf("%2d",i)
	}
	fmt.Printf("\n")
	for i := 1;i <= N; i++ {
		for j := 0;j <= N; j++ {
			if j == 0 {
				fmt.Printf("%2d",i)
				continue
			}
			if a[i][j] == 0 {
				ch = "+"
			} else if a[i][j] == 1 {
				ch = "X"
			}else { ch = "O"}
			fmt.Printf(" %s",ch)
		}
		fmt.Printf("\n")
	}
}

func SendWS(conn *websocket.Conn,data interface{}) error{
	err := conn.WriteJSON(data)
	return err
}

func Judge(x int,y int) bool {
	var l[2] int
	dx := []int{1,-1,0,0,1,-1,1,-1}
	dy := []int{0,0,1,-1,1,-1,-1,1}
	for k := 0;k < 8; k++ {
		var xx,yy,i int
		xx = int(x)
		yy = int(y)
		i = 0
		for {
			i++
			xx += dx[k]
			yy += dy[k]
			if xx<=0 || xx>N || yy<=0 || yy>N {break}
			if a[xx][yy] == 2 {
				l[k%2] = i
			}else {break}
		}
		if k%2==1 && l[1]+l[0]+1>=5 {
			return true
		}
	}
	return false
}

func BoardInit() {
	for i:=  0; i<= N; i++ {
		for j := 0; j <= N; j++ {
			a[i][j] = 0
		}
	}
}