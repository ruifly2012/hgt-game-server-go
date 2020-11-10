package websocket

import (
	"context"
	"server/model/mongo"
	"server/protobuf"
)

// 游戏

// 读取数据
func Load(user UserInfo, c *Client, msg interface{}) {
	if user.RoomId != "" {
		// 房间是否存在
		room, exist := RoomManage.GetRoomInfo(user.RoomId)
		if exist {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolLoadRes,
				"code":     CodeSuccess,
				"data": &protobuf.LoadRes{
					Reconnect: true,
					RoomId:    room.RoomId,
					Password:  room.Password,
				},
			}
			return
		}
	}
	c.Send <- map[string]interface{}{
		"protocol": ProtocolLoadRes,
		"code":     CodeSuccess,
		"data":     &protobuf.LoadRes{},
	}
}

// 心跳机制
func HeartBeat(user UserInfo, c *Client, msg interface{}) {
	c.Send <- map[string]interface{}{
		"protocol": ProtocolHeartBeatRes,
		"code":     CodeSuccess,
		"data":     &protobuf.HeartBeatRes{},
	}
}

// 保存游戏对局数据
func (round *RoundInfo) SaveRoundData() {
	// 记录对局数据
	roundRecord := make(map[string]interface{})
	roundRecord["RoomId"] = round.RoomId
	roundRecord["McUserId"] = round.McUserId
	roundRecord["QuestionId"] = round.QuestionId
	memberList := make(map[string]interface{})
	for userId, member := range round.GetRoundMemberInfoMap() {
		memberList[userId] = member
	}
	roundRecord["Member"] = memberList
	roundId, _ := mongo.RoundRecord().InsertOne(context.TODO(), roundRecord)
	// 记录聊天数据
	chatList := make(map[string]interface{})
	for messageId, message := range round.GetRoundChatMap() {
		chatList[messageId] = message
	}
	roundChat := map[string]interface{}{
		"roundId":  roundId, // 记录对局的mongoId
		"roomId":   round.RoomId,
		"chatList": chatList,
	}
	_, _ = mongo.RoundChat().InsertOne(context.TODO(), roundChat)
	// @todo 记录笔记数据

	// 删除对局数据
	RoundManage.Remove(round.RoomId)
}

// 记录日志数据
func RecordLog() {

}
