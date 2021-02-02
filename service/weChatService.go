package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"luck_draw/enums"
	"luck_draw/util"
	"net/http"
	"time"
)

const ACCESS_TOKEN_URL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential"
const WX_NOTIFY_URL = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token="

func WxNotify(data map[string]interface{})  {
	token,getTokenErr := GetWxAccessToken()
	if getTokenErr != nil {
		util.ErrDetail(getTokenErr.Code,getTokenErr.Err.Error(),"")
		return
	}

	url := WX_NOTIFY_URL+token
	byteData,encodeErr := json.Marshal(&data)
	if encodeErr != nil {
		return
	}

	httpClient := &util.HttpClient{}
	err := httpClient.Post(url,string(byteData),nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			util.ErrDetail(enums.WX_NOTIFY_IO_ERR,"模板消息读取数据流失败",readErr.Error())
		}
		util.Info(fmt.Sprintf("微信消息通知结果：%v",string(body)))
	})

	if err != nil {
		util.Error(fmt.Sprintf("微信消息通知异常：%v",string(byteData)))
	}
}

//通知奖品发放成功
func WxNotifyAward(id ,openid ,giftName,activityName,time,remark string)  {
	data := make(map[string]interface{})
	data["touser"] = openid
	data["template_id"] = enums.WX_TEMPLATE_SEND_SUCCESS
	data["page"] = fmt.Sprintf("/pages/home/index/index?path=pages/home/detail/detail&id=%v",id)
	data["data"] = map[string]map[string]string{
		"thing5":map[string]string{"value":giftName},
		"thing6":map[string]string{"value":activityName},
		"time3":map[string]string{"value":time},
		"thing7":map[string]string{"value":remark},
	}
	data["miniprogram_state"] = "formal"
	data["lang"] = "zh_CN"
	WxNotify(data)
}

//通知抽奖结果
func WxNotifyDraw(id ,openid ,activityName,result,time,giftName,remark string)  {
	data := make(map[string]interface{})
	data["touser"] = openid
	data["template_id"] = enums.WX_TEMPLATE_DRAW_FINISH
	data["page"] = fmt.Sprintf("/pages/home/index/index?path=pages/home/detail/detail&id=%v",id)
	data["data"] = map[string]map[string]string{
		"thing4":map[string]string{"value":activityName},
		"phrase5":map[string]string{"value":result},
		"date6":map[string]string{"value":time},
		"thing8":map[string]string{"value":giftName},
		"thing3":map[string]string{"value":remark},
	}
	data["miniprogram_state"] = "formal"
	data["lang"] = "zh_CN"
	WxNotify(data)
}

func GetWxAccessToken() (string,*enums.ErrorInfo) {
	var token interface{}
	var err error
	var ctx = context.Background()
	redis := util.NewRedis()
	cmd := redis.Client.Get(ctx,enums.WX_ACCESS_TOKEN_CACHE_KEY)
	if cmd.Err() == nil && len(cmd.Val()) > 0 {
		return cmd.Val(),nil
	}

	token,err = RequestWxAccessToken()
	if err != nil {
		//请求Access token失败
		util.ErrDetail(enums.AUTH_PARSE_ACCESS_REQUEST_ERR,"请求Access token失败",err.Error())
		return "",&enums.ErrorInfo{enums.WxAccessTokenRequestErr,enums.AUTH_PARSE_ACCESS_REQUEST_ERR}
	}

	redis.Client.SetEX(ctx,enums.WX_ACCESS_TOKEN_CACHE_KEY,token,time.Hour*2)
	return token.(string),nil
}

func RequestWxAccessToken() (string,error) {
	var err error
	config,_ := util.GetConfig()
	appId := config["WX_APP_ID"]
	appSecret := config["WX_APP_SECRET"]
	url := ACCESS_TOKEN_URL+"&appid="+appId+"&secret="+appSecret

	client := util.HttpClient{}
	data := make(map[string]interface{})
	err = client.Get(url,nil, func(resp *http.Response) {
		body, readIoErr := ioutil.ReadAll(resp.Body)
		if readIoErr != nil {
			err = readIoErr
			util.ErrDetail(enums.AUTH_WX_ACCESSTOKEN_READ_IO_ERR,"获取Access token读取数据流错误",readIoErr.Error())
			return
		}

		err = json.Unmarshal(body,&data)
		if err != nil {
			util.ErrDetail(enums.AUTH_PARSE_ACCESS_TOKEN_ERR,"解析Access token失败",err)
			return
		}
	})

	if err != nil {
		util.ErrDetail(enums.AUTH_PARSE_ACCESS_REQUEST_ERR,"请求Access token失败",err.Error())
		return "",err
	}

	return data["access_token"].(string),err
}
