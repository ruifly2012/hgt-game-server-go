package websocket

import "server/protobuf"

// 聊天

type ChatInfo struct {
	Id      string
	Content string
	Answer  uint32
	Aid     string
	AvaName string
	AvaHead string
	Mc      bool
}

// 聊天数据转换 protobuf ChatMessageRes 对象
func (chatInfo ChatInfo) ChangeChatToProtobuf() *protobuf.ChatMessageRes {
	return &protobuf.ChatMessageRes{
		Id:      chatInfo.Id,
		Content: chatInfo.Content,
		Answer:  chatInfo.Answer,
		Aid:     chatInfo.Aid,
		AvaName: chatInfo.AvaName,
		AvaHead: chatInfo.AvaHead,
		Mc:      chatInfo.Mc,
	}
}