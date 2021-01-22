package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"luck_draw/enums"
	"luck_draw/util"
	"net/http"
	"time"
)

func SocketNotify(userId string,code int,msg string) error {
	var errInfo *enums.ErrorInfo
	var token interface{}
	var err error

	token,errInfo = GetSocketToken(userId)
	if errInfo != nil {
		util.Info("获取token："+errInfo.Err.Error())
		return errInfo.Err
	}

	config,_ := util.GetConfig()
	url := "https://"+config["SOCKET_DOMAIN"]+"/push?token="+fmt.Sprintf("%v",token)
	data := make(map[string]interface{})
	data["message"] = map[string]interface{}{
		"message":msg,
		"code":code,
	}
	byteData,encodeErr := json.Marshal(&data)
	if encodeErr != nil {
		util.Info("解析错误："+encodeErr.Error())
		return encodeErr
	}

	httpClient := &util.HttpClient{}
	err = httpClient.Post(url,string(byteData),nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		util.Info("结果返回："+string(body))
		if readErr != nil {
			err = readErr
			util.ErrDetail(enums.SOCKET_AUTH_ERR,"socket消息推送读取数据流失败",readErr.Error())
		}

		var mp map[string]interface{}
		unErr :=json.Unmarshal(body,&mp)
		if unErr != nil {
			util.ErrDetail(enums.SOCKET_AUTH_ERR,"socket消息推送解析socket返回json数据异常",string(body))
		}

		if mp["code"].(float64) != float64(200) {
			util.ErrDetail(enums.SOCKET_AUTH_ERR,"socket消息推送解析socket返回json数据异常",string(body))
		}
	})

	return  err
}

func GetSocketToken(id string) (interface{},*enums.ErrorInfo) {
	key := enums.SOCKET_USER_TOKEN+":"+id
	redisClient := util.NewRedis()
	ctx := context.Background()
	exitKey := redisClient.Client.Get(ctx,key)
	if len(exitKey.Val()) <= 0 {
		tokenData,err := GetSocketTokenByApi("1")
		if err != nil {
			return nil,err
		}

		expire := tokenData["expire"].(float64) - float64(time.Now().Unix())
		redisClient.Client.SetEX(ctx,key,tokenData["token"],time.Duration(expire)*time.Second)
		return tokenData["token"],nil
	}else{
		return exitKey.Val(),nil
	}
}

func GetSocketTokenByApi(id string) (map[string]interface{},*enums.ErrorInfo) {
	config,_ := util.GetConfig()
	var token map[string]interface{}
	isok := true

	encrypt,err := util.AesEncryptECB([]byte(id),[]byte(config["SOCKET_SIGN_KEY"]))
	if err != nil {
		return nil,&enums.ErrorInfo{enums.SocketEncreyErr,enums.SOCKET_ENCRYPE_ERR}
	}

	decoded := base64.StdEncoding.EncodeToString(encrypt)
	data := make(map[string]interface{})
	data["sign"] = string(decoded)

	byteData,encodeErr := json.Marshal(&data)
	if encodeErr != nil {
		return nil,&enums.ErrorInfo{enums.SystemErr,enums.SOCKET_SIGN_ENCODE_ERR}
	}

	url := "https://"+config["SOCKET_DOMAIN"]+"/auth"
	httpClient := &util.HttpClient{}
	err = httpClient.Post(url,string(byteData),nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			isok = false
			util.ErrDetail(enums.SOCKET_AUTH_ERR,"读取数据流失败",readErr.Error())
		}

		var mp map[string]interface{}
		unErr :=json.Unmarshal(body,&mp)
		if unErr != nil {
			isok = false
			util.ErrDetail(enums.SOCKET_AUTH_ERR,"解析socket返回json数据异常",string(body))
		}

		if mp["code"].(float64) == float64(200) {
			token = mp["data"].(map[string]interface{})
		}else{
			isok = false
			util.ErrDetail(enums.SOCKET_AUTH_ERR,"请求socket授权返回失败",string(body))
		}
	})

	if err != nil {
		util.ErrDetail(enums.SOCKET_POST_SIGN_ERR,err.Error(),nil)
		return nil,&enums.ErrorInfo{enums.SystemErr,enums.SOCKET_POST_SIGN_ERR}
	}

	if !isok {
		return nil,&enums.ErrorInfo{enums.SystemErr,enums.SOCKET_POST_SIGN_ERR}
	}

	return token,nil
}
