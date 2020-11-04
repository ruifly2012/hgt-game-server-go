package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xlstudio/wxbizdatacrypt"
	"golang.org/x/crypto/bcrypt"
	"server/app"
	"server/dto"
	"server/exception"
	model "server/model/mysql"
	"server/util"
)

type EncryptedDataUserInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
}

type AuthApi struct{}


// 小程序登录
func (api AuthApi) AppletLogin(c *gin.Context) {
	var code = app.PostParam(c, "code")
	var appId = app.PostParam(c, "appId")
	_, sessionKey := util.GetSessionKeyByCode(code.(string), appId.(string))
	var iv = app.PostParam(c, "iv")
	var encryptedData = app.PostParam(c, "encryptedData")
	pc := wxbizdatacrypt.WxBizDataCrypt{AppId: appId.(string), SessionKey: sessionKey}
	result, _ := pc.Decrypt(encryptedData.(string), iv.(string), true)
	var userInfo *EncryptedDataUserInfo
	json.Unmarshal([]byte(result.(string)), &userInfo)
	// 获取用户是否存在
	userBindModel := &model.UserBind{OpenId: userInfo.OpenID, Status: 1}
	hasBind, _ := app.DB.Get(userBindModel)
	var userId string
	var userModel *model.User
	if hasBind {
		// 已经存在绑定用户
		userId = userBindModel.UserId
		userModel = &model.User{UserId: userId}
		app.DB.Get(userModel)
	} else {
		userId, _ = app.GenerateSnowflakeID()
		// 新用户
		userModel = &model.User{
			UserId:   userId,
			Username: userInfo.NickName,
			Avatar:   userInfo.AvatarURL,
			Gender:   userInfo.Gender,
			Role:     0,
			Status:   0,
		}
		app.DB.InsertOne(userModel)
		// 插入user_bind
		userBindId, _ := app.GenerateSnowflakeID()
		newUserBind := model.UserBind{
			UserBindId: userBindId,
			UserId:   userId,
			Nickname: userInfo.NickName,
			AppId:    appId.(string),
			OpenId:   userInfo.OpenID,
			UnionId:  userInfo.UnionID,
			Remark:   result.(string),
			Status:   1,
		}
		app.DB.InsertOne(newUserBind)
	}
	accessToken := util.BuildToken(dto.UserDTO{
		UserId:   userId,
		Username: userModel.Username,
		Avatar:   userModel.Avatar,
		Role:     userModel.Role,
		Status:   userModel.Status,
	})

	c.Set("data", map[string]interface{}{
		"accessToken": accessToken,
		"expire":      7 * 86400,
		"userInfo": map[string]interface{}{
			"userId":   userModel.UserId,
			"username": userModel.Username,
			"avatar":   userModel.Avatar,
			"role":     userModel.Role,
			"status":   userModel.Status,
		},
	})
}

// 注册
func (api AuthApi) Register (c *gin.Context) {
	username := app.PostParam(c, "username").(string)
	// 账号 手机号/邮箱
	account := app.PostParam(c, "account").(string)
	password := app.PostParam(c, "password").(string)
	accountType := app.CheckAccountType(account)
	if accountType == 0 {
		exception.Logic("账号类型仅支持手机/邮箱")
	}
	if username == "" {
		exception.Logic("昵称不能为空")
	}
	if password == "" {
		exception.Logic("密码不能为空")
	}
	if len(password) < 6 {
		exception.Logic("密码不能低于六位数字")
	}
	var checkUser *model.User
	// 手机号
	if accountType == 1 {
		checkUser = &model.User{Mobile: account}
	} else {
		checkUser = &model.User{Email: account}
	}
	hasUser, _ := app.DB.Get(checkUser)
	if hasUser {
		exception.Logic("该账号已经存在")
	}
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	userId, _ := app.GenerateSnowflakeID()
	// 新用户
	userModel := &model.User{
		UserId:   userId,
		Username: username,
		Password: string(hashPwd),
		Role:     0,
		Status:   0,
	}
	if accountType == 1 {
		userModel.Mobile = account
	} else {
		userModel.Email = account
	}
	app.DB.InsertOne(userModel)

	c.Set("data", true)
}

// 通过账密登录
func (api AuthApi) Login(c *gin.Context) {
	// 账号 手机号/邮箱
	account := app.PostParam(c, "account").(string)
	password := app.PostParam(c, "password").(string)
	accountType := app.CheckAccountType(account)
	if accountType == 0 {
		exception.Logic("账号类型仅支持手机/邮箱")
	}
	var userModel *model.User
	// 手机号
	if accountType == 1 {
		userModel = &model.User{Mobile: account}
	} else {
		userModel = &model.User{Email: account}
	}
	hasUser, _ := app.DB.Get(userModel)
	if !hasUser {
		exception.Logic("账号不存在")
	}
	err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password)) //验证（对比）
	if err != nil {
		exception.Logic("密码错误")
	}

	accessToken := util.BuildToken(dto.UserDTO{
		UserId:   userModel.UserId,
		Username: userModel.Username,
		Avatar:   userModel.Avatar,
		Role:     userModel.Role,
		Status:   userModel.Status,
	})

	c.Set("data", map[string]interface{}{
		"accessToken": accessToken,
		"expire":      7 * 86400,
		"userInfo": map[string]interface{}{
			"userId":   userModel.UserId,
			"username": userModel.Username,
			"avatar":   userModel.Avatar,
			"role":     userModel.Role,
			"status":   userModel.Status,
		},
	})
}
