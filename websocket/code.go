package websocket

const (
	CodeSuccess = 200
	CodeError = 500
	// ---------- 创建房间相关
	// 房间不合法
	CodeCreateRoomNameIllegal = 20001
	// 房间人数设置不合法
	CodeCreateRoomMaxIllegal = 20002
	// 房间创建失败
	CodeCreateRoomFailure = 20003
	// 已经创建房间 退出即可
	CodeCreateRoomExist = 20004

	// ---------- 房间相关
	// 房间不存在
	CodeRoomNotExist = 20101
	// 当前房间不在游戏中
	CodeRoomNotGaming =20102
	// 当前房间在游戏中
	CodeRoomGaming = 20103
	// 房间成员不存在
	CodeRoomMemberNotExist = 20104
	// 当前房间不在选题中
	CodeRoomNotSelectQuestion = 20105

	// ---------- 加入房间相关
	// 加入房间失败
	CodeJoinRoomFailure = 20200
	// 房间已经满人
	CodeRoomAlreadyMax = 20201
	// 已经加入房间
	CodeAlreadyExistRoom = 20202

	// ---------- 房间推送相关
	// 房间推送失败
	CodeRoomPush = 20300

	// ---------- 离开房间相关
	// 离开房间失败
	CodeLeaveRoomFailure = 20400
	// 游戏中 不能离开房间
	CodeCantLeaveCauseGaming = 20401

	// ---------- 准备相关
	// 准备失败
	CodeRoomPrepareFailure = 20500
	// 还有玩家没有准备
	CodeRoomSomeMemberNotPrepare = 20501
	// 人数太少 最少三个人
	CodeRoomMaxTooLittle = 20502
	// 游戏已经开始 不可以准备
	CodeGameStartRefusePrepare = 20503
	// 你已经准备 请勿重复准备
	CodeYourAlreadyPrepare = 20504
	// 你已经取消 请勿重复取消
	CodeYourAlreadyCancel = 20505

	// ---------- 踢人相关
	// 踢人失败
	CodeKickMemberFailure = 20600
	// 踢人参数错误
	CodeKickParamError = 20601
	// 此人不存在 无法踢人
	CodeKickMemberNotExist = 20602

	// ---------- 交换位置相关
	// 交换位置失败
	CodeExchangePositionFailure = 20700
	// 当前座位已经有人
	CodePositionBusy = 20701

	// ---------- 结束游戏相关
	// 结束游戏失败
	CodeEndGameFailure = 20800
	// 游戏不在游戏中不能结束
	CodeNotGamingCantEnd = 20801

	// ---------- 聊天相关
	// 聊天被限制
	CodeChatLimited = 20900
	// 聊天内容不合法
	CodeChatContentIllegal = 20901
	// 场次记录不存在
	CodeRoundNotExist = 20902
	// 聊天记录不存在
	CodeChatNotExist = 20903
	// 说话太快了
	CcodeChatFastLimit = 20904

	// ---------- 成员相关
	// 成员不是闲置状态
	CodeMemberBusy = 21000
	// 成员不是MC
	CodeMemberNotMC = 21001
	// 成员不是房主
	CodeMemberNotOwner = 21002

	// ----------- 回答相关
	// 答案类型不存在
	CodeAnswerTypeWrong = 21100
	// 不是mc 不具备权限
	CodeJustMcToReply = 21101

	// ----------- 题目相关
	// 题目不存在
	CodeQuestionExist = 21200
	// 不是mc 没有权限选题
	CodeNotRankToSelectQuestion = 21201
)
