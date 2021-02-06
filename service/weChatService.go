package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"luck_draw/enums"
	"luck_draw/util"
	"math/rand"
	"net/http"
	"sort"
	"strings"
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

	fmt.Println("thing3")
	fmt.Println(remark)

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

type WeChatPay struct {
	XMLName   		xml.Name `xml:"xml"`
	MchAppid        string   `xml:"mch_appid"`
	Mchid 			string   `xml:"mchid"`
	NonceStr  		string   `xml:"nonce_str"`
	PartnerTradeNo  string   `xml:"partner_trade_no"`
	Openid  		string   `xml:"openid"`
	CheckName  		string   `xml:"check_name"`
	Amount       	int		 `xml:"amount"`
	Sign       		string	 `xml:"sign"`
	Desc 			string 	 `xml:"desc"`
}

type WeChatPayResult struct {
	XMLName   		xml.Name `xml:"xml"`
	ReturnCode 			string 	 `xml:"return_code"`
	ReturnMsg 			string 	 `xml:"return_msg"`
	MchAppid 			string 	 `xml:"mch_appid"`
	Mchid 				string 	 `xml:"mchid"`
	DeviceInfo 			string 	 `xml:"device_info"`
	NonceStr 			string 	 `xml:"nonce_str"`
	ResultCode 			string 	 `xml:"result_code"`
	PartnerTradeNo 		string 	 `xml:"partner_trade_no"`
	PaymentNo 			string 	 `xml:"payment_no"`
	PaymentTime 		string 	 `xml:"payment_time"`
}

func KeySort(data map[string]interface{}) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
func RandChar(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i ++ {
		s.WriteByte(char[rand.Int63() % int64(len(char))])
	}
	return s.String()
}

func Pay(weChatPay WeChatPay)  {
	nano := fmt.Sprintf("%v",time.Now().UnixNano())
	config,_ := util.GetConfig()

	weChatPay.MchAppid = config["WE_CHAT_APP_ID"]
	weChatPay.Mchid = config["WE_CHAT_MCHID"]
	weChatPay.CheckName = "NO_CHECK"
	weChatPay.NonceStr = nano

	data := make(map[string]interface{})
	data["mch_appid"] 			= weChatPay.MchAppid
	data["mchid"] 				= weChatPay.Mchid
	data["nonce_str"] 			= weChatPay.NonceStr
	data["partner_trade_no"] 	= weChatPay.PartnerTradeNo
	data["openid"] 				= weChatPay.Openid
	data["check_name"] 			= weChatPay.CheckName
	data["amount"] 				= weChatPay.Amount
	data["desc"] 				= weChatPay.Desc

	keys := KeySort(data)

	dataStr := ""
	for index,key := range keys {
		if index == len(keys) - 1 {
			dataStr += fmt.Sprintf("%v=%v",key,data[key])
		}else{
			dataStr += fmt.Sprintf("%v=%v&",key,data[key])
		}
	}

	//拼接api秘钥
	stringSignTemp := dataStr+"&key="+config["WE_CHAT_PAY_API_KEY"]

	//md5
	h := md5.New()
	h.Write([]byte(stringSignTemp)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	md5Str := strings.ToUpper(hex.EncodeToString(cipherStr))

	weChatPay.Sign = md5Str
	xmlByte,err:=xml.Marshal(weChatPay)
	if err != nil {
		util.Error(fmt.Sprintf("解析xml数据出错：%v",err.Error()))
		return
	}

	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	var wechatPayCert = config["WE_CHAT_PAY_CERT"]
	var wechatPayKey = config["WE_CHAT_PAY_KEY"]

	//var tr *http.Transport
	// 微信提供的API证书,证书和证书密钥 .pem格式
	certs, err := tls.LoadX509KeyPair(wechatPayCert, wechatPayKey)
	if err != nil {
		util.Error(fmt.Sprintf("certs load err:", err.Error()))
		return
	} else {
		// 微信支付HTTPS服务器证书的根证书  .pem格式
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				//RootCAs:      pool,
				Certificates: []tls.Certificate{certs},
			},
		}

		client := &http.Client{Transport: tr}

		req, createErr := http.NewRequest("POST", url, bytes.NewBuffer(xmlByte))
		if createErr != nil {
			util.Error(fmt.Sprintf("创建失败:%v\n",createErr.Error()))
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			util.Error(fmt.Sprintf("微信支付请求失败：%v",err.Error()))
			return
		}

		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			util.Error(fmt.Sprintf("读取数据流失败：%v",readErr.Error()))
			return
		}

		util.Info(fmt.Sprintf("请求微信支付结果：%v",string(body)))

		payResult := &WeChatPayResult{}
		xml.Unmarshal(body,payResult)
		if payResult.ResultCode == "FAIL" {
			fmt.Println("支付失败")
		}

		fmt.Println(payResult)
	}
}
