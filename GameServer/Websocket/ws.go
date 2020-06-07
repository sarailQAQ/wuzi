package ws

import (
	"Server/Struct"
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
)

// Client 单个 websocket 信息
type Client struct {
	Id,Group,Username string
	Uid,Ready uint
	Socket *websocket.Conn
	Message chan []byte
	Data chan *interface{}
}

// messageData 单个发送数据信息
type MessageData struct {
	Id,Group,Sender string
	Message []byte
}

// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group,Sender string
	Message []byte
}

type BroadCastMessageData struct {
	Sender string
	Message []byte
}

// Manager 所有 websocket 信息
type Manager struct {
	Group map[string]map[string]*Client
	groupCount, clientCount uint
	Lock sync.Mutex
	Register, UnRegister chan *Client
	Message chan *MessageData
	GroupMessage chan *GroupMessageData
	BroadCastMessage chan *BroadCastMessageData
	GameMsg map[string]chan *GameMessageData
}

// 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group: make(map[string]map[string]*Client),
	Register:    make(chan *Client, 128),
	UnRegister:  make(chan *Client, 128),
	GroupMessage:   make(chan *GroupMessageData, 128),
	Message:   make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	GameMsg: make(map[string]chan *GameMessageData),
	groupCount: 0,
	clientCount: 0,
}

type GameMessageData struct {
	User string
	group string
	Type string
	Px,Py int
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close();err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Username, err)
		}
	}()

	for {
		var data interface{}
		err := c.Socket.ReadJSON(&data)
		if err != nil 	{ fmt.Println(err);break }
		if dataMap,ok := data.(map[string]interface{}); ok {
			typ := Struct.TypeString(dataMap["type"])
			switch typ {
			case "ready":WebsocketManager.ReadyClient(c)
			case "over":
				px,_ := strconv.Atoi(Struct.TypeString(dataMap["px"]))
				py,_ := strconv.Atoi(Struct.TypeString(dataMap["py"]))
				data := GameMessageData{
					User:  Struct.TypeString(dataMap["user"]),
					group: c.Group,
					Type:  typ,
					Px:    px,
					Py:    py,
				}
				WebsocketManager.GameOver(c.Group,data)
				return

			case "play":

				WebsocketManager.SendPlay(c.Group,Struct.OptData{
					Type: typ,
					User: Struct.TypeString(dataMap["user"]),
					Px:   Struct.TypeString(dataMap["px"]),
					Py:   Struct.TypeString(dataMap["py"]),
				})
			}

		}
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Printf("client [%s] disconnect", c.Username)
		if err := c.Socket.Close();err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Username, err)
		}
	}()

	for {
		select {
		case data := <- c.Data :

			log.Println("client write message:", c.Username, data)
			err := c.Socket.WriteJSON(data)
			if err != nil {
				log.Printf("client [%s] writemessage err: %s", c.Username, err)
			}
		}
	}
}

func (c *Client) MsgUnion(msg []byte) []byte {
	return bytes.Join([][]byte{[]byte(c.Username+"  "),[]byte(time.Now().Format("2006-01-02 15:04:05")),[]byte("\n"),msg},[]byte(""))
}

func (c *Client) WriteJSON(data interface{}) error {
	err := c.Socket.WriteJSON(data)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

