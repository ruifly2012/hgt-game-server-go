package websocket

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
	"server/app"
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
			} else {
				fmt.Println("消息推送 这个人已经掉线")
			}
		}
	} else {
		fmt.Println("发送聊天时 对局不存在")
	}
}

// MC回答
func Answer(c *Client, msg interface{}) {
	answerReq := msg.(*protobuf.AnswerReq)
	if round, ok := RoundManage.GetRoundInfo(c.RoomId); ok {
		// 判断是否时mc
		if round.McUserId != c.UserDTO.UserId {
			fmt.Println("很抱歉，只有mc才具备回复")
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
					} else {
						fmt.Println("消息推送 这个人已经掉线")
					}
				}
			} else {
				fmt.Println("mc 回答的消息不存在")
			}
		} else {
			fmt.Println("mc 回答的结果超出预期")
		}
	} else {
		fmt.Println("mc 回答时对局不存在")
	}
}

// 对局结束
func End(c *Client, msg interface{}) {
	fmt.Println("游戏结束")
	if round, ok := RoundManage.GetRoundInfo(c.RoomId); ok {
		// 判断是否时mc
		if round.McUserId != c.UserDTO.UserId {
			fmt.Println("很抱歉，只有mc才具备公布汤底")
			return
		}
		// @todo 保存数据
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
			} else {
				fmt.Println("用户可能掉线/离开")
			}
		}
	}
}

// 保存游戏对局数据
func (round *RoundInfo) saveRoundData() {

}
