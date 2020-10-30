package model

import (
	"time"
)

// 游戏主表
type Game struct {
	GameId      string    // 游戏id
	RoomId      string    // 房间id
	Password    string    // 房间密码
	QuestionId  string    // 题库id
	GroupUserId string    // 房主id
	McUserId    string    // mc用户id
	Number      uint8     // 房间人数上限
	Status      uint8     // 状态 0：未开始 1：开始中 2：已结束 3：关闭
	CreatedAt   int64     `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}
