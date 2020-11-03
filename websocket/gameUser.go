package websocket

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
)

// 用户信息
type UserInfo struct {
	// 用户id
	UserId string
	// 用户名
	Username string
	// 头像
	Avatar string
	// 房间id
	RoomId string
	// 1:闲置 2:在房间 3:准备中 4:游戏中
	Status uint32
	// 进入时间
	InsertTime string
	// 最后一次更新时间
	LastUpdateTime string
}

type UserManageStruct struct {
	cmap.ConcurrentMap
}

var UserManage = UserManageStruct{
cmap.New(),
}

// @todo 断线重连
func Reconnect(c *Client, msg interface{}) {

}

// 设置用户数据
func (user UserInfo) SaveUserInfo() {
	// @todo 判断是否存在用户数据 考虑重连问题

	UserManage.Set(user.UserId, user)
}

// 获取用户信息
func GetUserInfo(userId string) UserInfo {
	userInterface, _ := UserManage.Get(userId)
	userInfo := userInterface.(UserInfo)
	fmt.Println(userInfo)

	return userInfo
}

// 设置status
func (user UserInfo) SetStatus(status uint32)  {
	user.Status = status
	UserManage.Set(user.UserId, user)
	fmt.Println("设置status")
}

// 更新roomId
func (user UserInfo) SetRoomId(roomId string) {
	user.RoomId = roomId
	UserManage.Set(user.UserId, user)
	fmt.Println("设置roomId")
}