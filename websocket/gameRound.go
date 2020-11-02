package websocket

import (
	"context"
	"fmt"
	"github.com/orcaman/concurrent-map"
	"server/app"
	"server/model/mongo"
	"server/protobuf"
)

// 聊天回答 0:未回答 1:不相关 2:是 3:否 4:半对
const (
	AnswerStatusUnanswered = 0
	AnswerStatusNotRelate  = 1
	AnswerStatusYes        = 2
	AnswerStatusNot        = 3
	AnswerStatusHalf       = 4
)

var AnswerStatusManage = map[uint32]bool{
	AnswerStatusUnanswered: true,
	AnswerStatusNotRelate:  true,
	AnswerStatusYes:        true,
	AnswerStatusNot:        true,
	AnswerStatusHalf:       true,
}

type RoundInfo struct {
	// 房间id
	RoomId string
	// MC用户id
	McUserId string
	// 问题id
	QuestionId string
	// 聊天列表 map[messageId]protobuf.ChatMessageRes
	ChatList cmap.ConcurrentMap
	// 成员列表
	Member cmap.ConcurrentMap
}

type RoundManageStruct struct {
	cmap.ConcurrentMap
}

var RoundManage = RoundManageStruct{
	cmap.New(),
}

// 创建对局
func (room *RoomInfo) CreateRound() {
	roundInfo := RoundInfo{
		RoomId:     room.RoomId,
		McUserId:   room.McUserId,
		QuestionId: room.Question.Id,
	}
	roundInfo.ChatList = cmap.New()
	roundInfo.Member = cmap.New()
	for userId, member := range room.GetRoomMemberMap() {
		roundInfo.Member.Set(userId, member)
	}
	RoundManage.Set(room.RoomId, roundInfo)
}

// 聊天或提问消息
func Chat(c *Client, msg interface{}) {
	chatReq := msg.(*protobuf.ChatReq)
	if round, ok := RoundManage.GetRoundInfo(c.RoomId); ok {
		messageId, _ := app.GenerateSnowflakeID()
		isMc := false
		if c.UserDTO.UserId == round.McUserId {
			isMc = true
		}
		newMessage := protobuf.ChatMessageRes{
			Id:      messageId,
			Content: chatReq.Content,
			Answer:  AnswerStatusUnanswered,
			Aid:     c.UserDTO.UserId,
			AvaName: c.UserDTO.Username,
			AvaHead: c.UserDTO.Avatar,
			Mc:      isMc,
		}
		// 加入消息列表
		round.ChatList.Set(messageId, newMessage)
		RoundManage.Set(c.RoomId, round)
		// 往对局成员推送消息
		for userId, member := range round.GetRoundMemberMap() {
			fmt.Println("聊天消息发送："+member.AvaName, "内容："+newMessage.Content)
			if client, ok := Manager.clients[userId]; ok {
				client.Send <- map[string]interface{}{
					"protocol": ProtocolRoomPush,
					"code":     200,
					"data": &protobuf.RoomPush{
						ChangedMsg: []*protobuf.ChatMessageRes{
							&newMessage,
						},
					},
				}
			}
		}
	} else {
		// 对局不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolChatRes,
			"code": CodeRoundNotExist,
		}
	}
}

// MC回答
func Answer(c *Client, msg interface{}) {
	answerReq := msg.(*protobuf.AnswerReq)
	if round, ok := RoundManage.GetRoundInfo(c.RoomId); ok {
		// 判断是否时mc
		if round.McUserId != c.UserDTO.UserId {
			//只有mc才具备回复
			c.Send <- map[string]interface{}{
				"protocol": ProtocolAnswerRes,
				"code": CodeJustMcToReply,
			}
			return
		}
		if _, ok := AnswerStatusManage[answerReq.Answer]; ok {
			if chatMessage, ok := round.GetRoundChat(answerReq.Id); ok {
				chatMessage.Answer = answerReq.Answer
				// 更新消息
				round.ChatList.Set(answerReq.Id, chatMessage)
				RoundManage.Set(c.RoomId, round)
				// 推送更新
				for userId, _ := range round.GetRoundMemberMap() {
					fmt.Println("聊天回答结果发送：", chatMessage.Answer)
					if client, ok := Manager.clients[userId]; ok {
						client.Send <- map[string]interface{}{
							"protocol": ProtocolRoomPush,
							"code":     200,
							"data": &protobuf.RoomPush{
								ChangedMsg: []*protobuf.ChatMessageRes{
									&chatMessage,
								},
							},
						}
					}
				}
			} else {
				// 回答的消息不存在
				c.Send <- map[string]interface{}{
					"protocol": ProtocolAnswerRes,
					"code": CodeChatNotExist,
				}
				return
			}
		} else {
			// 答案类型不存在
			c.Send <- map[string]interface{}{
				"protocol": ProtocolAnswerRes,
				"code": CodeAnswerTypeWrong,
			}
			return
		}
	} else {
		// 对局不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolChatRes,
			"code": CodeRoundNotExist,
		}
	}
}

// 对局结束
func End(c *Client, msg interface{}) {
	fmt.Println("游戏结束")
	if round, ok := RoundManage.GetRoundInfo(c.RoomId); ok {
		// 判断是否时mc
		if round.McUserId != c.UserDTO.UserId {
			// 只有mc才具备公布汤底
			c.Send <- map[string]interface{}{
				"protocol": ProtocolEndRes,
				"code": CodeEndGameFailure,
			}
			return
		}
		// 保存数据
		go round.saveRoundData()
		// 修改房间数据
		room, _ := RoomManage.GetRoomInfo(c.RoomId)
		for userId, member := range room.GetRoomMemberMap() {
			member.Status = MemberStatusInRoom
			room.Member.Set(userId, member)
		}
		room.Status = RoomStatusPreparing
		// 更新房间数据
		RoomManage.Set(c.RoomId, room)
		// 推送房间数据
		for userId, _ := range room.GetRoomMemberMap() {
			if client, ok := Manager.clients[userId]; ok {
				// 推送整个 roomPush
				client.Send <- map[string]interface{}{
					"protocol": ProtocolRoomPush,
					"code":     200,
					"data": &protobuf.RoomPush{
						Status:      RoomStatusPreparing,
						SeatsChange: room.GetRoomMemberList(),
						RoomId:      c.RoomId,
						Question:    &room.Question,
					},
				}
			}
		}
	}
}

// 保存游戏对局数据
func (round *RoundInfo) saveRoundData() {
	// 记录对局数据
	roundRecord := make(map[string]interface{})
	roundRecord["RoomId"] = round.RoomId
	roundRecord["McUserId"] = round.McUserId
	roundRecord["QuestionId"] = round.QuestionId
	memberList := make(map[string]interface{})
	for userId, member := range round.GetRoundMemberMap() {
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
		"roundId": roundId, // 记录对局的mongoId
		"roomId": round.RoomId,
		"chatList": chatList,
	}
	mongo.RoundChat().InsertOne(context.TODO(), roundChat)
	// @todo 记录笔记数据
}
