package websocket

import "server/protobuf"

// 成员

type MemberInfo struct {
	Aid     string
	AvaName string
	AvaHead string
	Index   uint32
	Owner   bool
	Status  uint32
	Mc      bool
	Leave   bool
}

// 房间成员转 protobuf member对象
func (memberInfo MemberInfo) ChangeRoomMemberToProtobuf() *protobuf.RoomMemberSeatRes {
	return &protobuf.RoomMemberSeatRes{
		Aid:     memberInfo.Aid,
		AvaName: memberInfo.AvaName,
		AvaHead: memberInfo.AvaHead,
		Index:   memberInfo.Index,
		Owner:   memberInfo.Owner,
		Status:  memberInfo.Status,
		Mc:      memberInfo.Mc,
		Leave:   memberInfo.Leave,
	}
}