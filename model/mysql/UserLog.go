package model

import (
	"time"
)

// 用户日志表
type UserLog struct {
	UserLogId    string    // 用户日志操作表id
	UserId       string    // 用户id
	HandleUserId int64     // 操作的用户id
	Remark       string    // 备注
	CreatedAt    int64     `xorm:"created"`
	UpdatedAt    time.Time `xorm:"updated"`
}
