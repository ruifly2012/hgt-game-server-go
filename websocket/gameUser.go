package websocket

import cmap "github.com/orcaman/concurrent-map"

// 用户信息
type UserInfo struct {
	// 用户id
	Aid string
	// 用户名
	AvaName string
	// 头像
	AvaHead string
	// 房间id
	RoomId string
	// 1:闲置 2:在房间 3:准备中 4:游戏中
	Status uint32
}

type UserManageStruct struct {
	cmap.ConcurrentMap
}

var UserManage = UserManageStruct{
cmap.New(),
} 

// 断线重连
func Reconnect(c *Client, msg interface{}) {

}