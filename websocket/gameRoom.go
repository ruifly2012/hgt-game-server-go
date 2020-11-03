package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/orcaman/concurrent-map"
	"server/app"
	model "server/model/mysql"
	"server/protobuf"
)

// 玩家状态 1:闲置 2:在房间 3:准备中 4:游戏中
const (
	MemberStatusEmpty     = 1
	MemberStatusInRoom    = 2
	MemberStatusPreparing = 3
	MemberStatusGaming    = 4
)

// 房间状态 1：准备中 2：选题中 3：游戏中
const (
	RoomStatusPreparing      = 1
	RoomStatusSelectQuestion = 2
	RoomStatusGaming         = 3
)

type RoomInfo struct {
	// 房间id
	RoomId string
	// 房间密码
	Password string
	// 房主用户id
	OwnerUserId string
	// MC用户id
	McUserId string
	// 房间人数限制
	Max uint32
	// 房间状态 0：闲置 1：准备中 2：选题中 3：游戏中
	Status uint8
	// 问题
	Question protobuf.QuestionRes
	// 成员列表
	Member cmap.ConcurrentMap
}

type RoomManageStruct struct {
	cmap.ConcurrentMap
}

var RoomManage = RoomManageStruct{
	cmap.New(),
}

// 创建房间
func CreateRoom(user UserInfo, c *Client, msg interface{}) {
	if user.RoomId != "" {
		// 判断当前用户是否已经处于房间
		c.Send <- map[string]interface{}{
			"protocol": ProtocolCreateRoomRes,
			"code":     CodeCreateRoomExist,
		}
		return
	}
	createRoomReq := msg.(*protobuf.CreateRoomReq)
	if createRoomReq.Max == 0 || createRoomReq.Max > 10 {
		// 房间上限人数错误
		c.Send <- map[string]interface{}{
			"protocol": ProtocolCreateRoomRes,
			"code":     CodeCreateRoomMaxIllegal,
		}
		return
	}
	roomId, _ := app.GenerateSnowflakeID()
	var Member = cmap.New()
	Member.Set(user.UserId, protobuf.RoomMemberSeatRes{
		Aid:     user.UserId,
		AvaName: user.Username,
		AvaHead: user.Avatar,
		Index:   0,
		Owner:   true,
		Status:  MemberStatusPreparing,
		Mc:      true,
		Leave:   false,
	})
	room := RoomInfo{
		RoomId:      roomId,
		Password:    createRoomReq.Password,
		OwnerUserId: user.UserId,
		McUserId:    user.UserId,
		Max:         createRoomReq.Max,
		Status:      RoomStatusPreparing,
		Member:      Member,
	}
	RoomManage.Set(roomId, room)
	fmt.Println(user.Username + "创建了房间id：" + roomId)
	// 设置当前群id
	user.SetRoomId(roomId)
	// 设置用户状态
	user.SetStatus(MemberStatusPreparing)
	c.Send <- map[string]interface{}{
		"protocol": ProtocolCreateRoomRes,
		"code":     CodeSuccess,
		"data": &protobuf.CreateRoomRes{
			Room: &protobuf.RoomPush{
				Status:      RoomStatusPreparing,
				SeatsChange: room.GetRoomMemberList(),
				RoomId:      roomId,
			},
		},
	}
}

// 加入房间
func JoinRoom(user UserInfo, c *Client, msg interface{}) {
	joinRoomReq := msg.(*protobuf.JoinRoomReq)
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(joinRoomReq.RoomId); ok {
		if room.Status == RoomStatusGaming {
			// 游戏已经开局
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code":     CodeRoomGaming,
			}
			return
		}
		if int(room.Max) <= room.Member.Count() {
			// 房间人数上限，不能加入房间
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code":     CodeRoomAlreadyMax,
			}
			return
		}
		// 用户进入房间 将当前用户roomId 设置成这个
		user.SetRoomId(joinRoomReq.RoomId)
		// 设置用户状态
		user.SetStatus(MemberStatusInRoom)
		newUser := protobuf.RoomMemberSeatRes{
			Aid:     user.UserId,
			AvaName: user.Username,
			AvaHead: user.Avatar,
			Index:   0,
			Owner:   false,
			Status:  MemberStatusInRoom,
			Mc:      false,
			Leave:   false,
		}
		room.Member.Set(user.UserId, newUser)
		fmt.Println("用户：" + user.Username + " 加入房间：" + joinRoomReq.RoomId)
		// 通知群里面所有人
		for _, info := range room.GetRoomMemberMap() {
			if client, ok := Manager.clients[info.Aid]; ok {
				if info.Aid == newUser.Aid {
					// 获取完整对局数据
					round, okRound := RoundManage.GetRoundInfo(joinRoomReq.RoomId)
					var chatList []*protobuf.ChatMessageRes
					if okRound {
						chatList = round.GetRoundChatList()
					} else {
						chatList = nil
					}
					// 当前加入的用户，推送整个 roomPush
					client.Send <- map[string]interface{}{
						"protocol": ProtocolJoinRoomRes,
						"code":     CodeSuccess,
						"data": &protobuf.JoinRoomRes{
							Room: &protobuf.RoomPush{
								SeatsChange: room.GetRoomMemberList(),
								Question:    &room.Question,
								Status:      RoomStatusPreparing,
								RoomId:      joinRoomReq.RoomId,
								Msg:         chatList,
							},
						},
					}
				} else {
					// 只推送当前用户的信息给其他用户
					client.Send <- map[string]interface{}{
						"protocol": ProtocolRoomPush,
						"code":     CodeSuccess,
						"data": &protobuf.RoomPush{
							SeatsChange: []*protobuf.RoomMemberSeatRes{
								&newUser,
							},
						},
					}
				}
			}
		}
	} else {
		// 房间不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolJoinRoomRes,
			"code":     CodeRoomNotExist,
		}
		return
	}
}

// 离开房间
func LeaveRoom(user UserInfo, c *Client, msg interface{}) {
	var roomId = user.RoomId
	if roomId == "" {
		// 当前用户没有处于房间
		c.Send <- map[string]interface{}{
			"protocol": ProtocolLeaveRoomRes,
			"code":     CodeRoomMemberNotExist,
		}
		return
	}
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(roomId); ok {
		// 找群里面是否有这个用户
		if deleteUser, ok := room.GetRoomMember(user.UserId); ok {
			lastOneFlag := false
			// 判断当前成员列表是不是最后一个人
			if len(room.Member) == 1 {
				lastOneFlag = true
			}
			fmt.Println("用户：" + user.Username + " 离开了房间：" + roomId)
			// 离开的人的房间置为空
			user.SetRoomId("")
			// 设置用户状态 空闲
			user.SetStatus(MemberStatusEmpty)
			room.Member.Remove(user.UserId)
			if lastOneFlag {
				// 房间最后一个人退出，房间销毁
				RoomManage.Remove(roomId)
			} else {
				if room.OwnerUserId == deleteUser.Aid {
					// 离开的人是房主 切换房主
					ChangeOwner(deleteUser.Aid, roomId)
				}
				if room.McUserId == deleteUser.Aid {
					// 离开的人是MC 切换MC
					ChangeMC(deleteUser.Aid, roomId)
				}
				// 并非最后一个人 给其他人推送当前人离开数据
				for _, info := range room.GetRoomMemberMap() {
					if client, ok := Manager.clients[info.Aid]; ok {
						deleteUser.Leave = true
						client.Send <- map[string]interface{}{
							"protocol": ProtocolRoomPush,
							"code":     CodeSuccess,
							"data": &protobuf.RoomPush{
								SeatsChange: []*protobuf.RoomMemberSeatRes{
									&deleteUser,
								},
							},
						}
					}
				}
			}
		} else {
			// 用户不在房间里面
			c.Send <- map[string]interface{}{
				"protocol": ProtocolLeaveRoomRes,
				"code":     CodeRoomMemberNotExist,
			}
			return
		}
	} else {
		// 房间不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolLeaveRoomRes,
			"code":     CodeRoomNotExist,
		}
		return
	}
}

// 游戏准备
func Prepare(user UserInfo, c *Client, msg interface{}) {
	fmt.Println(user.RoomId, "----------------")
	prepare := msg.(*protobuf.PrepareReq)
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(user.RoomId); ok {
		if info, ok := room.GetRoomMember(user.UserId); ok {
			// mc准备 则考虑游戏是否要开始
			if user.UserId == room.McUserId {
				var roomMemberCount = room.Member.Count()
				// @todo
				if roomMemberCount < 0 {
					// 人数太少
					c.Send <- map[string]interface{}{
						"protocol": ProtocolPrepareRes,
						"code":     CodeRoomMaxTooLittle,
					}
					return
				}
				// 创建用户ids 长度 - mc
				var userIds = make([]string, roomMemberCount-1)
				// 判断用户是否都准备完毕
				for userId, member := range room.GetRoomMemberMap() {
					if member.Status != MemberStatusPreparing {
						// 玩家未准备
						c.Send <- map[string]interface{}{
							"protocol": ProtocolPrepareRes,
							"code":     CodeRoomSomeMemberNotPrepare,
						}
						return
					}
					if userId != room.McUserId {
						userIds = append(userIds, userId)
					}
				}
				fmt.Println("游戏开始选题")
				// 题目列表
				var questionResList []*protobuf.QuestionRes
				// 问题列表数据
				var questionList = make([]model.Question, 0)
				// @todo
				if len(userIds) >= 0 {
					// 获取题目
					questionLog := make([]model.UserQuestionLog, 0)
					app.DB.Cols("question_id").In("user_id", userIds).Find(&questionLog)
					var unQuestionIds = make([]string, 0)
					for _, log := range questionLog {
						unQuestionIds = append(unQuestionIds, log.QuestionId)
					}
					app.DB.NotIn("question_id", unQuestionIds).Limit(10).Find(&questionList)
				}
				// 处理问题列表
				for _, question := range questionList {
					questionResList = append(questionResList, &protobuf.QuestionRes{
						Id:       question.QuestionId,
						Title:    question.Title,
						Question: question.Description,
						Content:  question.Content,
					})
				}
				// 所有玩家已经准备完毕 推送玩家进入游戏
				for userId, info := range room.GetRoomMemberMap() {
					// 将玩家改成游戏中状态
					info.Status = MemberStatusGaming
					room.Member.Set(userId, info)
					RoomManage.Set(user.RoomId, room)
					// 设置用户状态 游戏中
					user.SetStatus(MemberStatusPreparing)
					if client, ok := Manager.clients[userId]; ok {
						client.Send <- map[string]interface{}{
							"protocol": ProtocolRoomPush,
							"code":     CodeSuccess,
							"data": &protobuf.RoomPush{
								Status:          RoomStatusSelectQuestion,
								SelectQuestions: questionResList,
							},
						}
					}
					// 创建对局
					room.CreateRound()
				}
			} else {
				if prepare.Ok {
					// 准备
					info.Status = MemberStatusPreparing
				} else {
					// 取消准备 在房间GetRoomMember
					info.Status = MemberStatusInRoom
				}
				// 更新当前用户信息
				room.Member.Set(user.UserId, info)
				RoomManage.Set(user.RoomId, room)
				// 设置用户状态 准备中
				user.SetStatus(MemberStatusPreparing)
				// 推送消息
				for _, info := range room.GetRoomMemberMap() {
					if client, ok := Manager.clients[info.Aid]; ok {
						if info.Aid != user.UserId {
							client.Send <- map[string]interface{}{
								"protocol": ProtocolRoomPush,
								"code":     CodeSuccess,
								"data": &protobuf.RoomPush{
									SeatsChange: []*protobuf.RoomMemberSeatRes{
										&info,
									},
								},
							}
						}
					}
				}
			}
		} else {
			// 用户不在房间里面
			c.Send <- map[string]interface{}{
				"protocol": ProtocolPrepareRes,
				"code":     CodeRoomMemberNotExist,
			}
			return
		}
	} else {
		// 房间不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolPrepareRes,
			"code":     CodeRoomNotExist,
		}
		return
	}
}

// 选题
func SelectQuestion(user UserInfo, c *Client, msg interface{}) {
	selectQuestionReq := msg.(*protobuf.SelectQuestionReq)
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(user.RoomId); ok {
		// 判断是否mc
		if room.McUserId != user.UserId {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolSelectQuestionRes,
				"code":     CodeNotRankToSelectQuestion,
			}
			return
		}
		if room.Status == RoomStatusGaming {
			// 游戏已经开局
			c.Send <- map[string]interface{}{
				"protocol": ProtocolSelectQuestionRes,
				"code":     CodeRoomGaming,
			}
			return
		}
		// 查找问题是否存在
		question := &model.Question{QuestionId: selectQuestionReq.Id}
		exist, _ := app.DB.Get(question)
		if !exist {
			// 题库不存在
			c.Send <- map[string]interface{}{
				"protocol": ProtocolSelectQuestionRes,
				"code":     CodeQuestionExist,
			}
			return
		}
		room.Question = protobuf.QuestionRes{
			Id:       question.QuestionId,
			Title:    question.Title,
			Question: question.Description,
			Content:  question.Content,
		}
		// 游戏开始
		room.Status = RoomStatusGaming
		RoomManage.Set(user.RoomId, room)
		// 通知群里面所有人
		for _, info := range room.GetRoomMemberMap() {
			if client, ok := Manager.clients[info.Aid]; ok {
				client.Send <- map[string]interface{}{
					"protocol": ProtocolRoomPush,
					"code":     CodeSuccess,
					"data": &protobuf.RoomPush{
						Status:   RoomStatusGaming,
						RoomId:   user.RoomId,
						Question: &room.Question,
					},
				}
			}
		}
	} else {
		// 房间不存在
		c.Send <- map[string]interface{}{
			"protocol": ProtocolSelectQuestionRes,
			"code":     CodeRoomNotExist,
		}
		return
	}
}

// 换房主
func ChangeOwner(sourceOwnerUserId string, roomId string) {
	room, _ := RoomManage.GetRoomInfo(roomId)
	for _, info := range room.GetRoomMemberMap() {
		if info.Aid != sourceOwnerUserId {
			room.OwnerUserId = info.Aid
			RoomManage.Set(roomId, room)
			break
		}
	}
}

// 换MC
func ChangeMC(sourceOwnerUserId string, roomId string) {
	room, _ := RoomManage.GetRoomInfo(roomId)
	for _, info := range room.GetRoomMemberMap() {
		if info.Aid != sourceOwnerUserId {
			room.McUserId = info.Aid
			RoomManage.Set(roomId, room)
			break
		}
	}
}

// 测试请求
func Test(user UserInfo, c *Client, msg interface{}) {
	fmt.Println("房间所有数据")
	for _, roomInterface := range RoomManage.Items() {
		room := roomInterface.(RoomInfo)
		jsons, _ := json.Marshal(room)
		fmt.Println(string(jsons))
	}
	fmt.Println("所有对局数据")
	for _, roundInterface := range RoundManage.Items() {
		round := roundInterface.(RoundInfo)
		jsons, _ := json.Marshal(round)
		fmt.Println(string(jsons))
	}
}
