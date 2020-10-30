package model

import (
	"time"
)

// 游戏笔记表
type GameNotes struct {
	GameNotesId string    // 笔记id
	GameId      string    // 游戏id
	GameMsgId   string    // 游戏消息id
	UserId      string    // 用户id
	Status      uint8     // 状态 1：正常 2：删除
	CreatedAt   int64     `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}
