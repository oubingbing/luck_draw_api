package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
)

/**
 * 微信小程序授权登录
 */
func Login(ctx *gin.Context)  {
	loginType,ok := ctx.GetQuery("type")
	if !ok {
		util.ResponseJson(ctx,enums.AUTH_LOGIN_TYPE_ERR,enums.LoginTypeErr.Error(),nil)
		return
	}

	if loginType != "wechat" {

	}

	var loginData enums.WxMiniLoginData
	if err := ctx.ShouldBind(&loginData); err != nil {
		util.ResponseJson(ctx,enums.AUTH_PARAMS_ERROR,err.Error(),nil)
		return
	}


	data,errInfo := service.GetSessionInfo(&loginData)
	if errInfo != nil {
		util.ResponseJson(ctx,errInfo.Code,errInfo.Err.Error(),nil)
		return
	}

	userInfo,err := service.BindUser(data)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	token,err := service.UserLogin(db,userInfo)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",token)
	return
}

func CheckLogin(ctx *gin.Context)  {
	util.ResponseJson(ctx,enums.SUCCESS,"",nil)
	return
}

func GetUserInfo(ctx *gin.Context)  {
	uid,_ := ctx.Get("user_id")
	userId,cok := uid.(float64)
	if !cok {
		util.ErrDetail(enums.Auth_TRANS_UID_ERR,fmt.Sprintf("ID转化异常，用户user_id:%v",uid),cok)
		util.ResponseJson(ctx,enums.Auth_TRANS_UID_ERR,enums.UserIdTransErr.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	user,err := service.FindUserById(db,int64(userId))
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	decoded, isBase64 := base64.StdEncoding.DecodeString(user.NickName)
	if isBase64 == nil {
		user.NickName = string(decoded)
	}

	data := make(map[string]interface{})
	data["id"] 			= user.ID
	data["nickname"] 	= user.NickName
	data["gender"] 		= user.Gender
	data["avatar"] 		= user.AvatarUrl
	data["phone"] 		= user.Phone

	util.ResponseJson(ctx,enums.SUCCESS,"",data)
	return
}

func GetUserPhone(ctx *gin.Context)  {
	uid,_ := ctx.Get("user_id")
	var loginData enums.WxMiniLoginData
	if err := ctx.ShouldBind(&loginData); err != nil {
		util.ResponseJson(ctx,enums.AUTH_PARAMS_ERROR,err.Error(),nil)
		return
	}

	userInfo,errInfo := service.GetSessionInfo(&loginData)
	if errInfo != nil {
		util.ResponseJson(ctx,errInfo.Code,errInfo.Err.Error(),nil)
		return
	}

	util.Info(fmt.Sprintf("用户:%v,手机号信息：%v",uid,string(userInfo)))

	phoneData := &enums.UserPhone{}
	parseErr := json.Unmarshal(userInfo,phoneData)
	if parseErr != nil {
		util.ResponseJson(ctx,enums.SYSTEM_ERR,"手机号信息解析失败",nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	err := service.UpdatePhone(db,uid,phoneData.PhoneNumber)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",phoneData.PhoneNumber)
	return
}
