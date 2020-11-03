package websocket

import "server/protobuf"

// 游戏

// 读取数据
func Load(user UserInfo, c *Client, msg interface{}) {
	if user.RoomId != "" {
		// 房间是否存在
		room, exist := RoomManage.GetRoomInfo(user.RoomId)
		if exist {
			c.Send <- map[string]interface{}{
				"protocol": ProtocolLoadRes,
				"code": CodeSuccess,
				"data": &protobuf.LoadRes{
					Reconnect: true,
					RoomId: room.RoomId,
					Password: room.Password,
				},
			}
			return
		}
	}
	c.Send <- map[string]interface{}{
		"protocol": ProtocolLoadRes,
		"code": CodeSuccess,
		"data": &protobuf.LoadRes{},
	}
}
