package model

import (
	"time"
)

// 用户日志表
type UserQuestionLog struct {
	UserId     string    // 用户id
	QuestionId string    // 问题id
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
}
