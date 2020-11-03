package websocket

import (
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

var UserManage = &UserManageStruct{
	cmap.New(),
}

// 设置用户数据
func (user UserInfo) SaveUserInfo() {
	// 判断是否存在用户数据 考虑重连问题
	oldUserInfo, exist := GetUserInfo(user.UserId)
	if exist && oldUserInfo.RoomId != "" {
		// 房间是否存在
		_, exist := RoomManage.GetRoomInfo(oldUserInfo.RoomId)
		if exist {
			// 需要重连
			user.RoomId = oldUserInfo.RoomId
			user.Status = oldUserInfo.Status
		}
	}
	UserManage.Set(user.UserId, user)
}

// 获取用户信息
func GetUserInfo(userId string) (UserInfo, bool) {
	userInterface, exist := UserManage.Get(userId)
	if exist {
		userInfo := userInterface.(UserInfo)
		return userInfo, true
	}

	return UserInfo{}, false
}

// 设置status
func (user UserInfo) SetStatus(status uint32) UserInfo {
	user.Status = status
	UserManage.Set(user.UserId, user)

	return user
}

// 更新roomId
func (user UserInfo) SetRoomId(roomId string) UserInfo {
	user.RoomId = roomId
	UserManage.Set(user.UserId, user)

	return user
}