package websocket

import "server/protobuf"

// 游戏大厅

// 查询大厅数据
func RoomHall(user UserInfo, c *Client, _ interface{}) {
	c.Send <- map[string]interface{}{
		"protocol": ProtocolRoomHallRes,
		"code":     CodeSuccess,
		"data": &protobuf.RoomHallRes{
			Rooms: RoomManage.GetRoomPrepareList(),
		},
	}
}
