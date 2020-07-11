package ws

import (
	"Server/data_form"
	"Server/middle_ware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (manager *Manager) Start() {
	for {
		select {

		case client := <-manager.Register :
			manager.Lock.Lock()

			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount=1
			}

			if len(manager.Group[client.Group]) >= 2{
				err := client.Socket.WriteMessage(websocket.BinaryMessage,[]byte("房间已满"))
				if err != nil {	log.Println(err)}
			}else {
				manager.Group[client.Group][client.Username] = client
				manager.clientCount += 1
				manager.SendGroup(client.Group,[]byte(client.Username+"已进入房间"))
			}

			manager.Lock.Unlock()

		case client := <-manager.UnRegister:
			manager.Lock.Lock()
			if _,ok := manager.Group[client.Group];ok {
				close(client.Message)
				delete(manager.Group[client.Group],client.Username)
				manager.clientCount -= 1

			}
			manager.Lock.Unlock()

		//发送广播
		case data := <-manager.BroadCastMessage :
			var temp interface{}
			temp = data
			for _,v := range manager.Group {
				for _,conn := range v {
					conn.Data <- &temp
				}
			}

		//发送群聊消息
		case data := <-manager.GroupMessage:
			var temp interface{}
			temp = map[string]interface{}{"message":string(data.Message)}
			if groupMap,ok := manager.Group[data.Group]; ok {
				for _,conn :=range groupMap {
					conn.Data <- &temp
				}
			}

		//发送个人消息
		case data := <-manager.Message:
			var temp interface{}
			temp = data
			if groupMap, ok := manager.Group[data.Group]; ok {
				if conn, ok := groupMap[data.Id]; ok {
					conn.Data <- &temp
				}
			}
		}
	}
}

func (manager *Manager) WsClient(ctx *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin:func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn,err := upGrader.Upgrade(ctx.Writer, ctx.Request,nil)
	if err != nil {
		log.Printf("websocket connect error: %s", ctx.Param("channel"))
		return
	}

	var user middle_ware.UserClaim
	t,_ := ctx.Get("user")
	if  uc,ok := t.(middle_ware.UserClaim); !ok{
		//resps.Error(ctx,1001,errors.New("Not login yet"))
		user = middle_ware.UserClaim{
			Id:       0,
			Username: "customer",
		}
	}else {
		user = uc
	}

	client := &Client{
		Id:      uuid.NewV4().String(),
		Username: user.Username,
		Uid:     user.Id ,
		Group:   ctx.Param("channel"),
		Socket:  conn,
		Message: make(chan []byte,1024),
		Data:	 make(chan *interface{},1024),
		Ready: 0,
	}

	manager.RegisterClient(client)
	go client.Read()
	go client.Write()
	time.Sleep(time.Second*15)
}

// 向指定的 client 发送数据
func (manger *Manager) Send(id string,group string,msg []byte) {
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: msg,
	}
	manger.Message <- data
}

// 向指定的 Group 广播
func (manger *Manager) SendGroup(group string,msg []byte) {
	data := &GroupMessageData{
		Group: group,
		Message: msg,
	}
	manger.GroupMessage <- data
}

func (manger *Manager) SendGroupData(group string,data GroupMessageData) {
	manger.GroupMessage <- &data
}

// 广播
func (manager *Manager) SendAll(msg []byte,sender string) {
	data := & BroadCastMessageData{
		Sender:sender,
		Message:  msg,
	}
	manager.BroadCastMessage <- data
}

func  (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

func (manager *Manager) ReadyClient(client *Client) {
	client.Ready = (client.Ready + 1) % 2
	msg := client.Username
	if client.Ready == 1 {
		msg += "已经准备！"
	} else {
		msg += "取消了准备！"
	}
	manager.SendGroup(client.Group,[]byte(msg))

	//游戏是否可以开始
	if len(manager.Group[client.Group]) == 2 {
		for _,v := range manager.Group[client.Group] {
			if v.Ready == 0 {return}
		}
	}else {return}
	go manager.GameStart(client.Group)
}

func (manager *Manager) GameOver(group string,data GameMessageData) {
	manager.GameMsg[group] <- &data
}

func (manager *Manager) SendPlay(group string,opt data_form.OptData) {
	px,_ := strconv.Atoi(opt.Px)
	py,_ := strconv.Atoi(opt.Py)
	data := &GameMessageData{
		User:  opt.User,
		group: group,
		Type:  opt.Type,
		Px:   px,
		Py:   py,
	}
	manager.GameMsg[group] <- data
}

func (manager *Manager) GameStart(group string) {
	manager.SendGroup(group,[]byte("Game Start!"))
	time.Sleep(time.Millisecond*100)
	var users[2] string
	var clients[2] *Client
	p := 0
	for k,v := range manager.Group[group] {
		if p >= 2 {break}
		users[p] = k
		clients[p] = v
		p++
	}
	fmt.Println(users)
	fmt.Println(clients[0].Username)
	manager.GameMsg[group] = make(chan *GameMessageData,1024)
	defer close(manager.GameMsg[group])

	//决定先手，并将消息发送给客户端
	rand.Seed(time.Now().UnixNano())
	p = rand.Intn(2)
	p = 0
	var data1,data2 interface{}
	data1 = data_form.OptData{
		Type:     strconv.Itoa(1000+1),
		Message:  users[p]+" first hand.",
	}
	clients[p].Data <- &data1
	data2 = data_form.OptData{
		Type:     strconv.Itoa(1000),
		Message:  users[p^1]+" second hand.",
	}
	clients[p^1].Data <- &data2

	for {
		select {
		case data := <- manager.GameMsg[group]:
			if data.Type == "over"{
				optData := &data_form.OptData{
					Type: data.Type,
					User: data.User,
					Px:   strconv.Itoa(data.Px),
					Py:   strconv.Itoa(data.Py),
				}
				data1 = optData
				manager.SendGroupJSON(group,&data1)
				time.Sleep(time.Millisecond*100)
				break
			}else if data.User != users[p] {
				// 不是当前要下棋的选手
				continue
			} else {
				p = (p + 1) % 2
				opt := data_form.OptData{
					Type: "play",
					User: data.User,
					Px:   strconv.Itoa(data.Px),
					Py:   strconv.Itoa(data.Py),
				}
				data1 = opt
				manager.SendGroupJSON(group,&data1)
			}
		}
	}
}

func (manager *Manager) SendGroupJSON(group string,data *interface{}) {
	for _,v := range manager.Group[group] {
		v.Data <- data
	}
}

// 处理单个 client 发送数据
func (manager *Manager)	SendService() {

	for {
		select {
		case data := <-manager.Message :
			for _,v := range manager.Group {
				for _,conn := range v {
					conn.Message <- data.Message
				}
			}
		}
	}

}

// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}
