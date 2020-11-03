package websocket

import (
	"github.com/golang/protobuf/proto"
	"reflect"
	"server/protobuf"
)

var ProtocolHandler = make(map[int64]MessageInfo)

type MessageHandler func(userInfo UserInfo, c *Client, msg interface{})

// 查询大厅房间数据
const ProtocolRoomHallReq int64 = 2001
const ProtocolRoomHallRes int64 = -2001
// 创建房间
const ProtocolCreateRoomReq int64 = 2002
const ProtocolCreateRoomRes int64 = -2002
// 加入房间
const ProtocolJoinRoomReq int64 = 2003
const ProtocolJoinRoomRes int64 = -2003
// 离开房间
const ProtocolLeaveRoomReq int64 = 2004
const ProtocolLeaveRoomRes int64 = -2004
// 游戏准备/取消
const ProtocolPrepareReq int64 = 2005
const ProtocolPrepareRes int64 = -2005

// 聊天或提问
const ProtocolChatReq int64 = 2008
const ProtocolChatRes int64 = -2008
// MC回答
const ProtocolAnswerReq int64 = 2009
const ProtocolAnswerRes int64 = -2009
// 游戏结束
const ProtocolEndReq int64 = 2010
const ProtocolEndRes int64 = -2010
// 选择问题
const ProtocolSelectQuestionReq int64 = 2011
const ProtocolSelectQuestionRes int64 = -2011
// 读取数据
const ProtocolLoadReq int64 = 2012
const ProtocolLoadRes int64 = -2012

// 房间推送消息
const ProtocolRoomPush int64 = 2901

// 消息信息
type MessageInfo struct {
	// 结构体类型
	messageType    reflect.Type
	// 对应函数
	messageHandler MessageHandler
}

// 加载协议
func LoadProtocol() {
	// 查询大厅数据
	RegisterMessage(ProtocolRoomHallReq, &protobuf.RoomHallReq{}, RoomHall)
	RegisterMessage(ProtocolRoomHallRes, &protobuf.RoomHallRes{}, nil)
	// 创建房间
	RegisterMessage(ProtocolCreateRoomReq, &protobuf.CreateRoomReq{}, CreateRoom)
	RegisterMessage(ProtocolCreateRoomRes, &protobuf.CreateRoomRes{}, nil)
	// 加入房间
	RegisterMessage(ProtocolJoinRoomReq, &protobuf.JoinRoomReq{}, JoinRoom)
	RegisterMessage(ProtocolJoinRoomRes, &protobuf.JoinRoomRes{}, nil)
	// 离开房间
	RegisterMessage(ProtocolLeaveRoomReq, &protobuf.LeaveRoomReq{}, LeaveRoom)
	RegisterMessage(ProtocolLeaveRoomRes, &protobuf.LeaveRoomRes{}, nil)
	// 游戏准备/取消
	RegisterMessage(ProtocolPrepareReq, &protobuf.PrepareReq{}, Prepare)
	RegisterMessage(ProtocolPrepareRes, &protobuf.PrepareRes{}, nil)

	// 聊天或提问
	RegisterMessage(ProtocolChatReq, &protobuf.ChatReq{}, Chat)
	RegisterMessage(ProtocolChatRes, &protobuf.ChatRes{}, nil)
	// MC回答
	RegisterMessage(ProtocolAnswerReq, &protobuf.AnswerReq{}, Answer)
	RegisterMessage(ProtocolAnswerRes, &protobuf.AnswerRes{}, nil)
	// 游戏结束
	RegisterMessage(ProtocolEndReq, &protobuf.EndReq{}, End)
	RegisterMessage(ProtocolEndRes, &protobuf.EndRes{}, nil)
	// 选择问题
	RegisterMessage(ProtocolSelectQuestionReq, &protobuf.SelectQuestionReq{}, SelectQuestion)
	RegisterMessage(ProtocolSelectQuestionRes, &protobuf.SelectQuestionRes{}, nil)

	// 读取数据
	RegisterMessage(ProtocolLoadReq, &protobuf.LoadReq{}, Load)
	RegisterMessage(ProtocolLoadRes, &protobuf.LoadRes{}, nil)

	// 房间推送消息
	RegisterMessage(ProtocolRoomPush, &protobuf.RoomPush{}, nil)

	// 测试
	RegisterMessage(3000, &protobuf.CreateRoomReq{}, Test)
}

// 注册协议对应方法
func RegisterMessage(protocol int64, msg interface{}, handler MessageHandler) {
	var info MessageInfo
	info.messageType = reflect.TypeOf(msg.(proto.Message))
	info.messageHandler = handler
	ProtocolHandler[protocol] = info
}