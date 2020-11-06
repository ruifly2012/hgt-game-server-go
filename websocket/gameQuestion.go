package websocket

import "server/protobuf"

// 问题
type QuestionInfo struct {
	Id          string
	Title       string
	Description string
	Content     string
}

// 将问题info转 protobuf
func (questionInfo QuestionInfo) ChangeChatToProtobuf() *protobuf.QuestionRes {
	return &protobuf.QuestionRes{
		Id:       questionInfo.Id,
		Title:    questionInfo.Title,
		Question: questionInfo.Description,
		Content:  questionInfo.Content,
	}
}