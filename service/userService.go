package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/util"
	"net/http"
	"strings"
)

const WECHAT_MINI_AUTH = "https://api.weixin.qq.com/sns/jscode2session"

/**
 * 用户登录
 */
func UserLogin(db *gorm.DB, user *model.User) (string,*enums.ErrorInfo) {
	existsUser := &model.User{}
	queryErr := existsUser.FindByOpenId(db,user.OpenId)
	if queryErr != nil && !gorm.IsRecordNotFoundError(queryErr) {
		return "",&enums.ErrorInfo{Code:enums.AUTH_USER_QUERY_ERR,Err:enums.LoginQueryUserErr}
	}

	var id uint
	openId := user.OpenId
	if gorm.IsRecordNotFoundError(queryErr) {
		//新增用户
		effect,saveErr := user.Store(db)
		if saveErr != nil {
			return "",&enums.ErrorInfo{Code:enums.AUTH_USER_SAVE_ERR,Err:enums.LoginSaveUserDbErr}
		}
		if effect <= 0 {
			return "",&enums.ErrorInfo{Code:enums.AUTH_USER_SAVE_ERR,Err:enums.LoginInsertUserErr}
		}
		id = user.ID
	}else{
		//跟新用户昵称和头像
		updateErr := UserUpdate(db,existsUser.ID,user.NickName,user.AvatarUrl)
		if updateErr != nil {
			return "",updateErr
		}
		id = existsUser.ID
	}

	token,tokenErr := util.CreateToken(id,openId)
	if tokenErr != nil {
		return "",&enums.ErrorInfo{Code:enums.AUTH_USER_CREATE_JWT_ERR,Err:enums.LoginCreateTokenErr}
	}

	return token,nil

}

/**
 * 更新用户头像和昵称
 */
func UserUpdate(db *gorm.DB,id uint,nickname string,avatar string) *enums.ErrorInfo {
	user := &model.User{}
	data := make(map[string]interface{})
	data["nick_name"] = nickname
	data["avatar_url"] = avatar
	err := user.Update(db,id,data)
	if err != nil {
		util.Error(fmt.Sprintf("UserUpdate-更新用户数据失败:%v,id:%v,nickname:%v,avatar:%v",err.Error(),id,nickname,avatar))
		return &enums.ErrorInfo{Code:enums.AUTH_USER_UPDATE_ERR,Err:enums.UpdateNicknameAvatarErr}
	}

	return nil
}

func UpdatePhone(db *gorm.DB,id interface{},phone string) *enums.ErrorInfo {
	user := &model.User{}
	data := make(map[string]interface{})
	data["phone"] = phone
	err := user.Update(db,id,data)
	if err != nil {
		return &enums.ErrorInfo{enums.UserUpdatePhoneErr,enums.AUTH_USER_UPDATE_PHONE_ERR}
	}

	return nil
}

/**
 * 通过微信服务器获取用户信息
 */
func GetSessionInfo(param *enums.WxMiniLoginData) ([]byte,*enums.ErrorInfo) {
	config ,configErr := util.GetConfig()
	if configErr != nil{
		util.Info(fmt.Sprintf("获取数据失败：%v\n",configErr.Err.Error()))
		return nil,configErr
	}

	url := getSessionUrl(config,param.Code)
	var err error
	sessionData := make(map[string]string)
	client := util.HttpClient{}
	err = client.Get(url,nil, func(resp *http.Response) {
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body,&sessionData)
		if err != nil{
			util.Error(fmt.Sprintf("GetSessionInfo-json错误:%v\n,json_data:%v",err,string(body)))
		}
	})

	if err != nil {
		util.Error(fmt.Sprintf("GetSessionInfo-校验session异常:%v\n,code:%v",sessionData,enums.AUTH_REQUEST_SESSION_ERR))
		return nil,&enums.ErrorInfo{Code:enums.AUTH_REQUEST_SESSION_ERR,Err:enums.LoginRequestSessionErr}
	}

	_,mpOk := sessionData["errcode"]
	if mpOk {
		util.Error(fmt.Sprintf("GetSessionInfo-校验session异常:%v\n,url:%v",sessionData,url))
		return nil,&enums.ErrorInfo{Code:enums.AUTH_REQUEST_SESSION_RESP_ERR,Err:enums.LoginFail}
	}

	userJson,err := util.AesDecrypt(param.EncryptedData,sessionData["session_key"],param.Iv)
	util.Error(fmt.Sprintf("GetSessionInfo-解析用户json出错:%v\n,json_data:%v",err,string(userJson)))
	if err != nil{
		util.Error(fmt.Sprintf("GetSessionInfo-解析用户json出错:%v\n,json_data:%v",err,string(userJson)))
		return nil,&enums.ErrorInfo{Code:enums.AUTH_PARSE_JSON_ERR,Err:enums.LoginParseUserJsonErr}
	}

	return userJson,nil
}

func BindUser(data []byte) (*model.User,*enums.ErrorInfo) {
	userData := model.User{}
	err := json.Unmarshal(data,&userData)
	if err != nil{
		util.Error(fmt.Sprintf("GetSessionInfo-解析用户json出错:%v\n,json_data:%v",err,string(data)))
		return nil,&enums.ErrorInfo{Code:enums.AUTH_PARSE_JSON_ERR,Err:enums.LoginParseUserJsonErr}
	}

	userData.FromType = model.USER_FROM_MINI

	return &userData,nil
}

/**
 * 拼接url
 */
func getSessionUrl(config map[string]string,code string) string {
	var strBuilder strings.Builder
	strBuilder.WriteString(WECHAT_MINI_AUTH)
	strBuilder.WriteString("?appId=")
	strBuilder.WriteString(config["WX_APP_ID"])
	strBuilder.WriteString("&secret=")
	strBuilder.WriteString(config["WX_APP_SECRET"])
	strBuilder.WriteString("&js_code=")
	strBuilder.WriteString(code)
	strBuilder.WriteString("&grant_type=authorization_code")
	return strBuilder.String()
}

func FindUserById(db *gorm.DB,id int64) (*model.User,*enums.ErrorInfo) {
	user := &model.User{}
	err := user.FindById(db,id)
	if err == gorm.ErrRecordNotFound {
		return nil,&enums.ErrorInfo{enums.UserNotFound,enums.AUTH_USER_NOT_FOUND}
	}

	return user,nil
}

func GetFakerUser(db *gorm.DB) ([]model.UserIDs,error) {
	var userList []model.UserIDs
	var err error

	ctx := context.Background()
	redis := util.NewRedis()
	defer redis.Client.Close()
	result := redis.Client.Get(ctx,model.FAKER_USER_LIST)
	if len(result.Val()) > 0 {
		//解析
		err = json.Unmarshal([]byte(result.Val()),&userList)
		return userList,err
	}else{
		fUser := &model.User{}
		userList,err = fUser.FakerUsers(db)
		if err != nil {
			return nil,err
		}

		var userListByte []byte
		userListByte,err = json.Marshal(&userList)
		redis.Client.Set(ctx,model.FAKER_USER_LIST,string(userListByte),0)
	}

	return  userList,nil
}


