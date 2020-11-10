package websocket

import (
	"container/list"
	"fmt"
	"server/app"
	"server/protobuf"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

// 聊天回答 1:未回答 2:不相关 3:是 4:否 5:半对
const (
	AnswerStatusUnanswered = 1
	AnswerStatusNotRelate  = 2
	AnswerStatusYes        = 3
	AnswerStatusNot        = 4
	AnswerStatusHalf       = 5
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
	// 聊天id 队列
	ChatQueue *list.List
	// 成员列表
	Member cmap.ConcurrentMap
}

type RoundManageStruct struct {
	cmap.ConcurrentMap
}

var RoundManage = RoundManageStruct{
	cmap.New(),
}

// 获取对局信息
func (rm *RoundManageStruct) GetRoundInfo(roomId string) (RoundInfo, bool) {
	if roundInterface, ok := rm.Get(roomId); ok {
		return roundInterface.(RoundInfo), ok
	}

	return RoundInfo{}, false
}

// 获取对局成员信息
func (round *RoundInfo) GetRoundMemberInfoMap() map[string]MemberInfo {
	memberList := make(map[string]MemberInfo)

	// Insert items to temporary map.
	for item := range round.Member.IterBuffered() {
		memberList[item.Key] = item.Val.(MemberInfo)
	}

	return memberList
}

// 获取对局所有聊天数据
func (round *RoundInfo) GetRoundChatList() []*protobuf.ChatMessageRes {
	var chatList []*protobuf.ChatMessageRes

	for e := round.ChatQueue.Front(); e != nil; e = e.Next() {
		messageId := (e.Value).(string)
		message, ok := round.GetRoundChatInfo(messageId)
		if ok {
			chatList = append(chatList, message.ChangeChatToProtobuf())
		}
	}

	return chatList
}

// 获取对局所有聊天数据
func (round *RoundInfo) GetRoundChatMap() map[string]ChatInfo {
	chatList := make(map[string]ChatInfo)

	// Insert items to temporary map.
	for item := range round.ChatList.IterBuffered() {
		chatList[item.Key] = item.Val.(ChatInfo)
	}

	return chatList
}

// 获取对局一条聊天数据
func (round *RoundInfo) GetRoundChatInfo(messageId string) (ChatInfo, bool) {
	if messageInterface, ok := round.ChatList.Get(messageId); ok {
		return messageInterface.(ChatInfo), ok
	}

	return ChatInfo{}, false
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
	roundInfo.ChatQueue = list.New()
	for userId, member := range room.GetRoomMemberInfoMap() {
		roundInfo.Member.Set(userId, member)
	}
	RoundManage.Set(room.RoomId, roundInfo)
}

// 聊天或提问消息
func Chat(user UserInfo, c *Client, msg interface{}) {
	lastSpeakTime := time.Now().Unix()
	if lastSpeakTime-user.LastSpeakTime <= 2 {
		c.Send <- map[string]interface{}{
			"protocol": ProtocolChatRes,
			"code":     CcodeChatFastLimit,
		}
		return
	}
	chatReq := msg.(*protobuf.ChatReq)
	if round, ok := RoundManage.GetRoundInfo(user.RoomId); ok {
		messageId, _ := app.GenerateSnowflakeID()
		isMc := false
		if user.UserId == round.McUserId {
			isMc = true
		}
		newChat := ChatInfo{
			Id:      messageId,
			Content: chatReq.Content,
			Answer:  AnswerStatusUnanswered,
			Aid:     user.UserId,
			AvaName: user.Username,
			AvaHead: user.Avatar,
			Mc:      isMc,
		}
		// 加入消息列表
		round.ChatList.Set(messageId, newChat)
		round.ChatQueue.PushBack(messageId)
		RoundManage.Set(user.RoomId, round)
		// 更新用户最后一次聊天时间
		user = user.SetLastSpeakTime(lastSpeakTime)
		// 往对局成员推送消息
		for userId, member := range round.GetRoundMemberInfoMap() {
			fmt.Println("聊天消息发送："+member.AvaName, "内容："+newChat.Content)
			if client, ok := Manager.clients[userId]; ok {
				client.Send <- map[string]interface{}{
					"protocol": ProtocolRoomPush,
					"code":     200,
					"data": &protobuf.RoomPush{
						ChangedMsg: []*protobuf.ChatMessageRes{
							newChat.ChangeChatToProtobuf(),
						},
					},
				}
			}
		}
	} else {
		// 对局不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolChatRes,
			"code":     CodeRoundNotExist,
		}
	}
}

// MC回答
func Answer(user UserInfo, c *Client, msg interface{}) {
	answerReq := msg.(*protobuf.AnswerReq)
	if round, ok := RoundManage.GetRoundInfo(user.RoomId); ok {
		// 判断是否时mc
		if round.McUserId != user.UserId {
			//只有mc才具备回复
			c.Send <- map[string]interface{}{
				"protocol": ProtocolAnswerRes,
				"code":     CodeJustMcToReply,
			}
			return
		}
		if _, ok := AnswerStatusManage[answerReq.Answer]; ok {
			if chatMessage, ok := round.GetRoundChatInfo(answerReq.Id); ok {
				chatMessage.Answer = answerReq.Answer
				// 更新消息
				round.ChatList.Set(answerReq.Id, chatMessage)
				RoundManage.Set(user.RoomId, round)
				// 推送更新
				for userId := range round.GetRoundMemberInfoMap() {
					fmt.Println("聊天回答结果发送：", chatMessage.Answer)
					if client, ok := Manager.clients[userId]; ok {
						client.Send <- map[string]interface{}{
							"protocol": ProtocolRoomPush,
							"code":     200,
							"data": &protobuf.RoomPush{
								ChangedMsg: []*protobuf.ChatMessageRes{
									chatMessage.ChangeChatToProtobuf(),
								},
							},
						}
					}
				}
			} else {
				// 回答的消息不存在
				c.Send <- map[string]interface{}{
					"protocol": ProtocolAnswerRes,
					"code":     CodeChatNotExist,
				}
				return
			}
		} else {
			// 答案类型不存在
			c.Send <- map[string]interface{}{
				"protocol": ProtocolAnswerRes,
				"code":     CodeAnswerTypeWrong,
			}
			return
		}
	} else {
		// 对局不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolAnswerRes,
			"code":     CodeRoundNotExist,
		}
	}
}

// 对局结束
func End(user UserInfo, c *Client, _ interface{}) {
	fmt.Println("游戏结束")
	if round, ok := RoundManage.GetRoundInfo(user.RoomId); ok {
		// 判断是否时mc
		if round.McUserId != user.UserId {
			// 只有mc才具备公布汤底
			c.Send <- map[string]interface{}{
				"protocol": ProtocolEndRes,
				"code":     CodeEndGameFailure,
			}
			return
		}
		// 修改房间数据
		room, _ := RoomManage.GetRoomInfo(user.RoomId)
		// 游戏中才能结束游戏
		if room.Status != RoomStatusGaming {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolEndRes,
				"code":     CodeNotGamingCantEnd,
			}
			return
		}
		// 保存数据
		go round.SaveRoundData()
		for userId, member := range room.GetRoomMemberInfoMap() {
			if userId == room.McUserId {
				member.Status = MemberStatusPreparing
			} else {
				member.Status = MemberStatusInRoom
			}
			room.Member.Set(userId, member)
		}
		room.Status = RoomStatusPreparing
		// 更新房间数据
		RoomManage.Set(user.RoomId, room)
		// 推送房间数据
		for userId := range room.GetRoomMemberInfoMap() {
			if client, ok := Manager.clients[userId]; ok {
				// 推送整个 roomPush
				client.Send <- map[string]interface{}{
					"protocol": ProtocolRoomPush,
					"code":     200,
					"data": &protobuf.RoomPush{
						Status:      RoomStatusPreparing,
						SeatsChange: room.GetRoomMemberList(),
						RoomId:      user.RoomId,
						Question:    room.Question.ChangeChatToProtobuf(),
					},
				}
			}
		}
	}
}
