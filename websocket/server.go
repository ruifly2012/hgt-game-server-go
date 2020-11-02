package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"server/app"
	"server/dto"
	"server/protobuf"
	"server/util"
	"time"
)

type ClientManager struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	UserDTO    dto.UserDTO
	socket     *websocket.Conn
	Send       chan map[string]interface{}
	InsertTime string
	Protocol   int64
	RoomId     string
}

var Manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[string]*Client),
}

// 服务
func Server(c *gin.Context) {
	//解析一个连接
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	token := c.Query("Authorization")
	if token == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Illegal connection"))
		conn.Close()
		return
	}
	go Manager.start()
	// 解析token
	userDto := util.CheckToken(token)

	//初始化一个客户端对象
	client := &Client{
		UserDTO:    userDto,
		socket:     conn,
		Send:       make(chan map[string]interface{}),
		InsertTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	//把这个对象发送给 管道
	Manager.register <- client
	// 协程接收输出信息
	go client.write()
	go client.read()
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: //新客户端加入
			// 判断用户是否还在
			if client, ok := manager.clients[conn.UserDTO.UserId]; ok {
				client.socket = conn.socket
			} else {
				manager.clients[conn.UserDTO.UserId] = conn
			}
			fmt.Println("新用户加入:"+conn.UserDTO.Username, "加入时间："+conn.InsertTime)
			fmt.Println("当前总用户数量register：", len(manager.clients))
		case conn := <-manager.unregister: // 离开
			if _, ok := manager.clients[conn.UserDTO.UserId]; ok {
				close(conn.Send)
				delete(manager.clients, conn.UserDTO.UserId)
				fmt.Println("用户离开：" + conn.UserDTO.Username)
				fmt.Println("当前总用户数量unregister：", len(manager.clients))
			}
			//case message := <-manager.broadcast: //读到广播管道数据后的处理
			//	fmt.Println("当前总用户数量broadcast：", len(manager.clients))
			//	for _, conn := range manager.clients {
			//		select {
			//		case conn.Send <- message: //调用发送给全体客户端
			//		default:
			//			// 重新上来之后挤掉了 @todo
			//			// 关闭连接
			//			close(conn.Send)
			//			delete(manager.clients, conn.UserDTO.UserId)
			//		}
			//	}
		}
	}
}

// 广播数据 除了自己
//func (manager *ClientManager) send(message []byte, ignore *Client) {
//	for _, conn := range manager.clients {
//		if conn != ignore {
//			conn.Send <- message //发送的数据写入所有的 websocket 连接 管道
//		}
//	}
//}

// 写入管道后激活这个进程
func (c *Client) write() {
	defer func() {
		if err := recover(); err != nil {
			// 错误记录
			app.GameServerRecover(err)
			// 恢复
			go c.write()
			// 给用户推送500错误
			c.Send <- map[string]interface{}{
				"protocol": -c.Protocol,
				"code":     500,
			}
		} else {
			// 用户正常退出
			Manager.unregister <- c
			c.socket.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.Send: //这个管道有了数据 写这个消息出去
			if !ok {
				// 发送关闭提示
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			protocol, _ := message["protocol"].(int64)
			code, _ := message["code"].(int)
			if _, ok := ProtocolHandler[protocol]; ok {
				var res []byte
				if _, ok := message["data"]; ok {
					childMessage := message["data"].(proto.Message)
					childBytes, _ := proto.Marshal(childMessage)

					baseMessage := &protobuf.Message{
						Protocol: protocol,
						Code:     int64(code),
						Data:     childBytes,
					}
					res, _ = proto.Marshal(baseMessage)
				} else {
					baseMessage := &protobuf.Message{
						Protocol: protocol,
						Code:     int64(code),
					}
					res, _ = proto.Marshal(baseMessage)
				}
				err := c.socket.WriteMessage(websocket.BinaryMessage, res)
				if err != nil {
					// 程序退出 关闭链接
					return
				}
			} else {
				fmt.Println("返回没有找到对应的 struct")
			}
		}
	}
}

// 客户端发送数据处理逻辑
func (c *Client) read() {
	defer func() {
		if err := recover(); err != nil {
			// 错误记录
			app.GameServerRecover(err)
			// 给用户推送500错误
			c.Send <- map[string]interface{}{
				"protocol": -c.Protocol,
				"code":     500,
			}
			// 恢复
			go c.read()
		} else {
			// 用户正常退出
			Manager.unregister <- c
			c.socket.Close()
		}
	}()

	for {
		// 监听从 socket 获取数据
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			// 数据获取错误 退出登录
			fmt.Println("read 关闭")
			return
		}
		// 基础结构体
		baseMessage := &protobuf.Message{}
		// proto解码
		proto.Unmarshal(message, baseMessage)
		fmt.Println(baseMessage.Protocol)
		// 找到对应的协议操作
		if info, ok := ProtocolHandler[baseMessage.Protocol]; ok {
			c.Protocol = baseMessage.Protocol
			infoMessage := reflect.New(info.messageType.Elem()).Interface()
			proto.Unmarshal(baseMessage.Data, infoMessage.(proto.Message))
			info.messageHandler(c, infoMessage)
		} else {
			panic("找不到协议对应的结构体")
		}
		//激活start 程序 入广播管道
		//websocketManager.broadcast <- message
	}
}
