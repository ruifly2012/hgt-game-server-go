package websocket

import (
	"server/protobuf"
)

// 获取roomInfo
func (rm *RoomManageStruct) GetRoomInfo(roomId string) (RoomInfo, bool) {
	if roomInterface, ok := rm.Get(roomId); ok {
		return roomInterface.(RoomInfo), ok
	}

	return RoomInfo{}, false
}

// 获取对局信息
func (rm *RoundManageStruct) GetRoundInfo(roomId string) (RoundInfo, bool) {
	if roundInterface, ok := rm.Get(roomId); ok {
		return roundInterface.(RoundInfo), ok
	}

	return RoundInfo{}, false
}

// 获取房间的成员
func (room *RoomInfo) GetRoomMember(userId string) (protobuf.RoomMemberSeatRes, bool) {
	if userInterface, ok := room.Member.Get(userId); ok {
		return userInterface.(protobuf.RoomMemberSeatRes), ok
	}

	return protobuf.RoomMemberSeatRes{}, false
}


// 获取房间成员信息
func (room *RoomInfo) GetRoomMemberMap() map[string]protobuf.RoomMemberSeatRes {
	memberList := make(map[string]protobuf.RoomMemberSeatRes)

	// Insert items to temporary map.
	for item := range room.Member.IterBuffered() {
		memberList[item.Key] = item.Val.(protobuf.RoomMemberSeatRes)
	}

	return memberList
}

// 获取对局成员信息
func (round *RoundInfo) GetRoundMemberMap() map[string]protobuf.RoomMemberSeatRes {
	memberList := make(map[string]protobuf.RoomMemberSeatRes)

	// Insert items to temporary map.
	for item := range round.Member.IterBuffered() {
		memberList[item.Key] = item.Val.(protobuf.RoomMemberSeatRes)
	}

	return memberList
}

// 获取所有成员
func (room *RoomInfo) GetRoomMemberList() []*protobuf.RoomMemberSeatRes {
	var roomMembers []*protobuf.RoomMemberSeatRes

	for _, member := range room.GetRoomMemberMap() {
		roomMembers = append(roomMembers, &member)
	}

	return roomMembers
}

// 获取对局一条聊天数据
func (round *RoundInfo) GetRoundChat(messageId string) (protobuf.ChatMessageRes, bool) {
	if messageInterface, ok := round.ChatList.Get(messageId); ok {
		return messageInterface.(protobuf.ChatMessageRes), ok
	}

	return protobuf.ChatMessageRes{}, false
}