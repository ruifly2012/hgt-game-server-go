package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/orcaman/concurrent-map"
	"server/app"
	"server/protobuf"
)

// 玩家状态 0:闲置 1:在房间 2:准备中 3:游戏中
const (
	MemberStatusEmpty     = 0
	MemberStatusInRoom    = 1
	MemberStatusPreparing = 2
	MemberStatusGaming    = 3
)

// 房间状态 0：闲置 1：待定 2：准备中 3：游戏中
const (
	RoomStatusEmpty     = 0
	RoomStatusHoldOn    = 1
	RoomStatusPreparing = 2
	RoomStatusGaming    = 3
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
	// 房间状态 0：闲置 1：待定 2：准备中 3：游戏中
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
func CreateRoom(c *Client, msg interface{}) {
	createRoomReq := msg.(*protobuf.CreateRoomReq)
	roomId, _ := app.GenerateSnowflakeID()
	var Member = cmap.New()
	Member.Set(c.UserDTO.UserId, protobuf.RoomMemberSeatRes{
		Aid:     c.UserDTO.UserId,
		AvaName: c.UserDTO.Username,
		AvaHead: c.UserDTO.Avatar,
		Index:   0,
		Owner:   true,
		Status:  MemberStatusPreparing,
		Mc:      true,
		Leave:   false,
	})
	room := RoomInfo{
		RoomId:      roomId,
		Password:    createRoomReq.Password,
		OwnerUserId: c.UserDTO.UserId,
		McUserId:    c.UserDTO.UserId,
		Max:         createRoomReq.Max,
		Status:      RoomStatusPreparing,
		Member:      Member,
		Question: protobuf.QuestionRes{
			Id:       "questionId",
			Question: "我是汤面",
			Content:  "我是汤底",
		},
	}
	RoomManage.Set(roomId, room)
	TestMember(room)
	fmt.Println(c.UserDTO.Username + "创建了房间id：" + roomId)
	// 设置当前群id
	c.RoomId = roomId
	c.Send <- map[string]interface{}{
		"protocol": ProtocolCreateRoomRes,
		"code":     200,
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
func JoinRoom(c *Client, msg interface{}) {
	joinRoomReq := msg.(*protobuf.JoinRoomReq)
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(joinRoomReq.RoomId); ok {
		if room.Status == 3 {
			// 游戏已经开局
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code": CodeRoomGaming,
			}
			return
		}
		if int(room.Max) <= room.Member.Count() {
			// 房间人数上限，不能加入房间
			c.Send <- map[string]interface{}{
				"protocol": ProtocolJoinRoomRes,
				"code": CodeRoomAlreadyMax,
			}
			return
		}
		// 用户进入房间 将当前用户roomId 设置成这个
		c.RoomId = joinRoomReq.RoomId
		newUser := protobuf.RoomMemberSeatRes{
			Aid:     c.UserDTO.UserId,
			AvaName: c.UserDTO.Username,
			AvaHead: c.UserDTO.Avatar,
			Index:   0,
			Owner:   false,
			Status:  MemberStatusInRoom,
			Mc:      false,
			Leave:   false,
		}
		room.Member.Set(c.UserDTO.UserId, newUser)
		fmt.Println("用户：" + c.UserDTO.Username + " 加入房间：" + joinRoomReq.RoomId)
		TestMember(room)
		// 通知群里面所有人
		for _, info := range room.GetRoomMemberMap() {
			if client, ok := Manager.clients[info.Aid]; ok {
				if info.Aid == newUser.Aid {
					// 当前加入的用户，推送整个 roomPush
					client.Send <- map[string]interface{}{
						"protocol": ProtocolRoomPush,
						"code":     200,
						"data": &protobuf.RoomPush{
							Status:      RoomStatusPreparing,
							SeatsChange: room.GetRoomMemberList(),
							RoomId:      joinRoomReq.RoomId,
						},
					}
				} else {
					// 只推送当前用户的信息给其他用户
					client.Send <- map[string]interface{}{
						"protocol": ProtocolRoomPush,
						"code":     200,
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
			"code": CodeRoomNotExist,
		}
		return
	}
}

// 离开房间
func LeaveRoom(c *Client, msg interface{}) {
	var roomId = c.RoomId
	if roomId == "" {
		fmt.Println("当前用户没有处于房间")
	} else {
		// 找房间
		if room, ok := RoomManage.GetRoomInfo(roomId); ok {
			// 找群里面是否有这个用户
			if deleteUser, ok := room.GetRoomMember(c.UserDTO.UserId); ok {
				lastOneFlag := false
				// 判断当前成员列表是不是最后一个人
				if len(room.Member) == 1 {
					lastOneFlag = true
				}
				fmt.Println("用户：" + c.UserDTO.Username + " 离开了房间：" + roomId)
				// 离开的人的房间置为空
				c.RoomId = ""
				room.Member.Remove(c.UserDTO.UserId)
				TestMember(room)
				if lastOneFlag {
					// 删除房间数据
					fmt.Println("房间最后一个人退出，房间销毁")
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
								"code":     200,
								"data": &protobuf.RoomPush{
									SeatsChange: []*protobuf.RoomMemberSeatRes{
										&deleteUser,
									},
								},
							}
						} else {
							fmt.Println("这个人不在不推送")
						}
					}
				}
			} else {
				fmt.Println("用户不在房间里面")
			}
		} else {
			fmt.Println("房间不存在")
		}
	}
}

// 游戏准备
func Prepare(c *Client, msg interface{}) {
	prepare := msg.(*protobuf.PrepareReq)
	// 找房间
	if room, ok := RoomManage.GetRoomInfo(c.RoomId); ok {
		if info, ok := room.GetRoomMember(c.UserDTO.UserId); ok {
			// mc准备 则考虑游戏是否要开始
			if c.UserDTO.UserId == room.McUserId {
				// 判断用户是否都准备完毕
				for _, member := range room.GetRoomMemberMap() {
					if member.Status != MemberStatusPreparing {
						fmt.Println("单独给mc推送消息：" + member.AvaName + ": 玩家未准备")
						return
					}
				}
				// 所有玩家已经准备完毕 推送玩家进入游戏
				for userId, info := range room.GetRoomMemberMap() {
					// 将玩家改成游戏中状态
					info.Status = MemberStatusGaming
					room.Member.Set(userId, info)
					fmt.Println("游戏开始")
					if client, ok := Manager.clients[userId]; ok {
						client.Send <- map[string]interface{}{
							"protocol": ProtocolRoomPush,
							"code":     200,
							"data": &protobuf.RoomPush{
								Status: RoomStatusGaming,
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
				room.Member.Set(c.UserDTO.UserId, info)
				// 推送消息
				for _, info := range room.GetRoomMemberMap() {
					if client, ok := Manager.clients[info.Aid]; ok {
						if info.Aid != c.UserDTO.UserId {
							client.Send <- map[string]interface{}{
								"protocol": ProtocolRoomPush,
								"code":     200,
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
			TestMember(room)
		} else {
			fmt.Println("没有在房间里面")
		}
	} else {
		fmt.Println("没有找到房间")
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

// 测试代码
func TestMember(room RoomInfo) {
	for _, info := range room.GetRoomMemberMap() {
		jsons, _ := json.Marshal(info)
		fmt.Println(string(jsons))
	}
}
