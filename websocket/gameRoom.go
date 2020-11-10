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
	// 房间名称
	RoomName string
	// 房间密码
	Password string
	// 房主用户id
	OwnerUserId string
	// MC用户id
	McUserId string
	// 房间人数限制
	Max uint32
	// 房间状态 0：闲置 1：准备中 2：选题中 3：游戏中
	Status uint32
	// 问题
	Question QuestionInfo
	// 成员列表
	Member cmap.ConcurrentMap
}

type RoomManageStruct struct {
	cmap.ConcurrentMap
}

var RoomManage = RoomManageStruct{
	cmap.New(),
}

// 获取roomInfo
func (rm *RoomManageStruct) GetRoomInfo(roomId string) (RoomInfo, bool) {
	if roomInterface, ok := rm.Get(roomId); ok {
		return roomInterface.(RoomInfo), ok
	}

	return RoomInfo{}, false
}

// 获取所有准备中的房间
func (rm *RoomManageStruct) GetRoomPrepareList() []*protobuf.RoomPush {
	var roomList []*protobuf.RoomPush
	for item := range rm.IterBuffered() {
		room := item.Val.(RoomInfo)
		if room.Status == RoomStatusPreparing {
			roomList = append(roomList, room.ChangeRoomToProtobuf())
		}
	}

	return roomList
}

// 将 自身 roomInfo 转变成 protobuf roomPush
func (room *RoomInfo) ChangeRoomToProtobuf() *protobuf.RoomPush {
	var hasPassword bool
	if room.Password != "" {
		hasPassword = true
	}

	return &protobuf.RoomPush{
		RoomId: room.RoomId,
		RoomName: room.RoomName,
		RoomMax: room.Max,
		RoomMemberNum: uint32(room.Member.Count()),
		HasPassword: hasPassword,
	}
}

// 获取房间成员信息
func (room *RoomInfo) GetRoomMemberInfoMap() map[string]MemberInfo {
	memberList := make(map[string]MemberInfo)

	// Insert items to temporary map.
	for item := range room.Member.IterBuffered() {
		memberList[item.Key] = item.Val.(MemberInfo)
	}

	return memberList
}

// 获取房间的成员info
func (room *RoomInfo) GetRoomMemberInfo(userId string) (MemberInfo, bool) {
	if userInterface, ok := room.Member.Get(userId); ok {
		return userInterface.(MemberInfo), ok
	}

	return MemberInfo{}, false
}

// 获取所有成员
func (room *RoomInfo) GetRoomMemberList() []*protobuf.RoomMemberSeatRes {
	var roomMembers []*protobuf.RoomMemberSeatRes

	for item := range room.Member.IterBuffered() {
		member := item.Val.(MemberInfo)
		roomMembers = append(roomMembers, member.ChangeRoomMemberToProtobuf())
	}

	return roomMembers
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
	Member.Set(user.UserId, MemberInfo{
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
		RoomName:    createRoomReq.Name,
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
	user = user.SetRoomId(roomId)
	// 设置用户状态
	user = user.SetStatus(MemberStatusPreparing)
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
		// 判断用户是否在房间里面
		newUser, memberExist := room.GetRoomMemberInfo(user.UserId)
		// 游戏已经开局 && 用户不在
		if room.Status == RoomStatusGaming && !memberExist {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code":     CodeRoomGaming,
			}
			return
		}
		// 房间人数上限，不能加入房间 && 用户不存在
		if int(room.Max) <= room.Member.Count() && !memberExist {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code":     CodeRoomAlreadyMax,
			}
			return
		}
		// 如果房间具备密码 密码不等则不能加入房间
		if room.Password != "" && joinRoomReq.Password != room.Password {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code":     CodeJoinRoomFailure,
			}
			return
		}
		// 用户不存在
		if !memberExist {
			// 用户进入房间 将当前用户roomId 设置成这个
			user = user.SetRoomId(joinRoomReq.RoomId)
			// 设置用户状态
			user = user.SetStatus(MemberStatusInRoom)
			newUser = MemberInfo{
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
		}
		fmt.Println("用户：" + user.Username + " 加入房间：" + joinRoomReq.RoomId)
		// 通知群里面所有人
		for _, info := range room.GetRoomMemberInfoMap() {
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
					// 题目列表
					var questionResList []*protobuf.QuestionRes
					if room.McUserId == newUser.Aid && room.Status == RoomStatusSelectQuestion {
						questionResList = getQuestionList(room)
					} else {
						questionResList = nil
					}
					roomPush := protobuf.RoomPush{
						SeatsChange:     room.GetRoomMemberList(),
						Status:          room.Status,
						RoomId:          joinRoomReq.RoomId,
						Msg:             chatList,
						McId:            room.McUserId,
						SelectQuestions: questionResList,
					}
					if room.Status == RoomStatusGaming {
						roomQuestion := protobuf.QuestionRes{
							Id:       room.Question.Id,
							Title:    room.Question.Title,
							Question: room.Question.Description,
						}
						if room.McUserId == newUser.Aid {
							roomQuestion.Content = room.Question.Content
						}
						roomPush.Question = &roomQuestion
					}
					// 当前加入的用户，推送整个 roomPush
					client.Send <- map[string]interface{}{
						"protocol": ProtocolJoinRoomRes,
						"code":     CodeSuccess,
						"data": &protobuf.JoinRoomRes{
							Room: &roomPush,
						},
					}
				} else {
					// 只推送当前用户的信息给其他用户
					client.Send <- map[string]interface{}{
						"protocol": ProtocolRoomPush,
						"code":     CodeSuccess,
						"data": &protobuf.RoomPush{
							SeatsChange: []*protobuf.RoomMemberSeatRes{
								newUser.ChangeRoomMemberToProtobuf(),
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
func LeaveRoom(user UserInfo, c *Client, _ interface{}) {
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
		if room.Status == RoomStatusGaming || room.Status == RoomStatusSelectQuestion {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolLeaveRoomRes,
				"code":     CodeCantLeaveCauseGaming,
			}
			return
		}
		// 找群里面是否有这个用户
		if deleteUser, ok := room.GetRoomMemberInfo(user.UserId); ok {
			lastOneFlag := false
			// 判断当前成员列表是不是最后一个人
			if len(room.Member) == 1 {
				lastOneFlag = true
			}
			fmt.Println("用户：" + user.Username + " 离开了房间：" + roomId)
			// 离开的人的房间置为空
			user = user.SetRoomId("")
			// 设置用户状态 空闲
			user = user.SetStatus(MemberStatusEmpty)
			room.Member.Remove(user.UserId)
			RoomManage.Set(room.RoomId, room)
			if lastOneFlag {
				// 房间最后一个人退出，房间销毁
				RoomManage.Remove(roomId)
			} else {
				var memberSeats = make([]*protobuf.RoomMemberSeatRes, 0)
				memberSeats = append(memberSeats, deleteUser.ChangeRoomMemberToProtobuf())
				if room.OwnerUserId == deleteUser.Aid {
					// 离开的人是房主 切换房主
					changeOwner, ok := ChangeOwner(deleteUser.Aid, roomId)
					if ok {
						memberSeats = append(memberSeats, changeOwner)
					}
				}
				if room.McUserId == deleteUser.Aid {
					// 离开的人是MC 切换MC
					changeMc, ok := ChangeMC(deleteUser.Aid, roomId)
					if ok {
						memberSeats = append(memberSeats, changeMc)
					}
				}
				// 并非最后一个人 给其他人推送当前人离开数据
				for _, info := range room.GetRoomMemberInfoMap() {
					if client, ok := Manager.clients[info.Aid]; ok {
						deleteUser.Leave = true
						client.Send <- map[string]interface{}{
							"protocol": ProtocolRoomPush,
							"code":     CodeSuccess,
							"data": &protobuf.RoomPush{
								SeatsChange: memberSeats,
							},
						}
					}
				}
			}
			// 给当前用户推送离开房间返回
			c.Send <- map[string]interface{}{
				"protocol": ProtocolLeaveRoomRes,
				"code":     CodeSuccess,
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
	prepare := msg.(*protobuf.PrepareReq)
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(user.RoomId); ok {
		// 游戏已经开始|选题中 不可以准备
		if room.Status == RoomStatusGaming || room.Status == RoomStatusSelectQuestion {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolPrepareRes,
				"code":     CodeGameStartRefusePrepare,
			}
			return
		}
		if memberInfo, ok := room.GetRoomMemberInfo(user.UserId); ok {
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
				for userId, member := range room.GetRoomMemberInfoMap() {
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
					_ = app.DB.Cols("question_id").In("user_id", userIds).Find(&questionLog)
					var unQuestionIds = make([]string, 0)
					for _, log := range questionLog {
						unQuestionIds = append(unQuestionIds, log.QuestionId)
					}
					_ = app.DB.NotIn("question_id", unQuestionIds).Limit(10).Find(&questionList)
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
				for userId, info := range room.GetRoomMemberInfoMap() {
					// 将玩家改成游戏中状态
					info.Status = MemberStatusGaming
					room.Member.Set(userId, info)
					room.Status = RoomStatusSelectQuestion
					RoomManage.Set(user.RoomId, room)
					// 设置用户状态 游戏中
					user = user.SetStatus(MemberStatusPreparing)
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
				}
			} else {
				// 你已经准备 请勿重复准备
				if memberInfo.Status == MemberStatusPreparing && prepare.Ok {
					c.Send <- map[string]interface{}{
						"protocol": ProtocolPrepareRes,
						"code":     CodeYourAlreadyPrepare,
					}
					return
				}
				// 你已经取消 请勿重复取消
				if memberInfo.Status == MemberStatusInRoom && !prepare.Ok {
					c.Send <- map[string]interface{}{
						"protocol": ProtocolPrepareRes,
						"code":     CodeYourAlreadyCancel,
					}
					return
				}
				if prepare.Ok {
					// 准备
					memberInfo.Status = MemberStatusPreparing
				} else {
					// 取消准备 在房间GetRoomMember
					memberInfo.Status = MemberStatusInRoom
				}
				// 更新当前用户信息
				room.Member.Set(user.UserId, memberInfo)
				RoomManage.Set(user.RoomId, room)
				// 设置用户状态 准备中
				user = user.SetStatus(MemberStatusPreparing)
				// 推送消息
				for _, info := range room.GetRoomMemberInfoMap() {
					if client, ok := Manager.clients[info.Aid]; ok {
						if info.Aid != user.UserId {
							client.Send <- map[string]interface{}{
								"protocol": ProtocolRoomPush,
								"code":     CodeSuccess,
								"data": &protobuf.RoomPush{
									SeatsChange: []*protobuf.RoomMemberSeatRes{
										info.ChangeRoomMemberToProtobuf(),
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
		room.Question = QuestionInfo{
			Id:          question.QuestionId,
			Title:       question.Title,
			Description: question.Description,
			Content:     question.Content,
		}
		// 游戏开始
		room.Status = RoomStatusGaming
		RoomManage.Set(user.RoomId, room)
		protobufQuestion := protobuf.QuestionRes{
			Id:       room.Question.Id,
			Title:    room.Question.Title,
			Question: room.Question.Description,
		}
		// 通知群里面所有人
		for _, info := range room.GetRoomMemberInfoMap() {
			if client, ok := Manager.clients[info.Aid]; ok {
				if client.UserId == room.McUserId {
					protobufQuestion.Content = room.Question.Content
				} else {
					protobufQuestion.Content = ""
				}
				client.Send <- map[string]interface{}{
					"protocol": ProtocolRoomPush,
					"code":     CodeSuccess,
					"data": &protobuf.RoomPush{
						Status:   RoomStatusGaming,
						RoomId:   user.RoomId,
						Question: &protobufQuestion,
					},
				}
			}
		}
		// 创建对局
		room.CreateRound()
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
func ChangeOwner(sourceOwnerUserId string, roomId string) (*protobuf.RoomMemberSeatRes, bool) {
	room, _ := RoomManage.GetRoomInfo(roomId)
	for _, info := range room.GetRoomMemberInfoMap() {
		if info.Aid != sourceOwnerUserId {
			info.Owner = true
			room.Member.Set(info.Aid, info)
			room.OwnerUserId = info.Aid
			RoomManage.Set(roomId, room)
			return info.ChangeRoomMemberToProtobuf(), true
		}
	}

	return &protobuf.RoomMemberSeatRes{}, false
}

// 换MC
func ChangeMC(sourceOwnerUserId string, roomId string) (*protobuf.RoomMemberSeatRes, bool) {
	room, _ := RoomManage.GetRoomInfo(roomId)
	for _, info := range room.GetRoomMemberInfoMap() {
		if info.Aid != sourceOwnerUserId {
			info.Mc = true
			room.Member.Set(info.Aid, info)
			room.McUserId = info.Aid
			RoomManage.Set(roomId, room)
			return info.ChangeRoomMemberToProtobuf(), true
		}
	}

	return &protobuf.RoomMemberSeatRes{}, false
}

// 测试请求
func Test(_ UserInfo, _ *Client, _ interface{}) {
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
	fmt.Println("所有用户数据")
	for item := range UserManage.IterBuffered() {
		user := item.Val.(UserInfo)
		jsons, _ := json.Marshal(user)
		fmt.Println(string(jsons))
	}
}

func getQuestionList(room RoomInfo) []*protobuf.QuestionRes {
	// 创建用户ids 长度 - mc
	var userIds = make([]string, room.Member.Count()-1)
	// 判断用户是否都准备完毕/游戏中
	for userId, member := range room.GetRoomMemberInfoMap() {
		if member.Status != MemberStatusPreparing && member.Status != MemberStatusGaming {
			// 玩家未准备
			return nil
		}
		if userId != room.McUserId {
			userIds = append(userIds, userId)
		}
	}
	// 题目列表
	var questionResList []*protobuf.QuestionRes
	// 问题列表数据
	var questionList = make([]model.Question, 0)
	// @todo
	if len(userIds) >= 0 {
		// 获取题目
		questionLog := make([]model.UserQuestionLog, 0)
		_ = app.DB.Cols("question_id").In("user_id", userIds).Find(&questionLog)
		var unQuestionIds = make([]string, 0)
		for _, log := range questionLog {
			unQuestionIds = append(unQuestionIds, log.QuestionId)
		}
		_ = app.DB.NotIn("question_id", unQuestionIds).Limit(10).Find(&questionList)
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

	return questionResList
}
