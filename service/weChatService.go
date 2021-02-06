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
	"log"
	"luck_draw/enums"
	"luck_draw/util"
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

type XmlData struct {
	MchAppid 		string         `xml:"mch_appid,attr"`
	Mchid     		string         `xml:"mchid,attr"`
	NonceStr     	string         `xml:"nonce_str,attr"`
	Sign     		string         `xml:"sign,attr"`
	PartnerTradeNo  string         `xml:"partner_trade_no,attr"`
	Openid     		string         `xml:"openid,attr"`
	CheckName     	string         `xml:"check_name,attr"`
	Amount     		int            `xml:"amount,attr"`
	Desc     		string         `xml:"desc,attr"`
}


type xmldas struct {
	XMLName  xml.Name       		`xml:"xml"`
	MchAppid 		string         `xml:"mch_appid,attr"`
	Mchid     		string         `xml:"mchid,attr"`
	NonceStr     	string         `xml:"nonce_str,attr"`
	Sign     		string         `xml:"sign,attr"`
	PartnerTradeNo  string         `xml:"partner_trade_no,attr"`
	Openid     		string         `xml:"openid,attr"`
	CheckName     	string         `xml:"check_name,attr"`
	Amount     		int            `xml:"amount,attr"`
	Desc     		string         `xml:"desc,attr"`
}

type xmlsource struct {
	Path  string `xml:"path,attr"`
	Param string `xml:"param,attr"`
}

type xmldestination struct {
	Path  string `xml:"path,attr"`
	Param string `xml:"param,attr"`
}

func KeySort(data map[string]interface{}) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Pay()  {

	data := make(map[string]interface{})
	data["mch_appid"] = "wxa0d7aa1607c5ac21"
	data["mchid"] = "1254223701"
	data["nonce_str"] = "qweqweqweqwe"
	data["partner_trade_no"] = "10000098201411111234567890"
	data["openid"] = "oc-wc48UGe_dVNqmAPJjIQAEw5-w"
	data["check_name"] = "NO_CHECK"
	data["amount"] = 100
	data["desc"] = "抽奖"

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
	stringSignTemp := fmt.Sprintf("%v&key=aY5fa4VPhDGWa1Qu19NuCVeBzChKIEWZ",dataStr)

	fmt.Printf("拼接key：%v\n",stringSignTemp)

	//方法一
	h := md5.New()
	h.Write([]byte(stringSignTemp)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	//fmt.Printf("%s\n", hex.EncodeToString(cipherStr)) // 输出加密结果

	md5Str := strings.ToUpper(hex.EncodeToString(cipherStr))

	fmt.Println(md5Str)
	/*bytes2:=sha256.Sum256([]byte(md5Str))//计算哈希值，返回一个长度为32的数组
	hashcode2:=hex.EncodeToString(bytes2[:])//将数组转换成切片，转换成16进制，返回字符串
	fmt.Println(hashcode2)*/

	//fmt.Println(md5str1)
	//fmt.Println(dataStr

	/*v := xmldas{}
	v.Src = xmlsource{Path: "123", Param: "456"}
	v.Dest = xmldestination{Path: "789", Param: "000"}
	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println(string(output))*/

	/*xData := XmlData{
		MchAppid:       "",
		Mchid:          "",
		NonceStr:       "",
		Sign:           "",
		PartnerTradeNo: "",
		Openid:         "",
		CheckName:      "",
		Amount:         0,
		Desc:           "",
	}

	output, err := xml.MarshalIndent(xData, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println(string(output))*/



	/*byteData,encodeErr := json.Marshal(&data)
	if encodeErr != nil {
		return
	}*/

	byteData := `
<xml>

<mch_appid>wxa0d7aa1607c5ac21</mch_appid>

<mchid>1254223701</mchid>

<nonce_str>qweqweqweqwe</nonce_str>

<partner_trade_no>10000098201411111234567890</partner_trade_no>

<openid>oc-wc48UGe_dVNqmAPJjIQAEw5-w</openid>

<check_name>NO_CHECK</check_name>

<amount>100</amount>

<desc>抽奖</desc>

<sign>BAF6451A7A3401D8B102E4C75226CC53</sign>

</xml>

`

	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"

	var wechatPayCert = "C:/Users/hasee/Downloads/WXCertUtil/cert/new/apiclient_cert.pem"
	var wechatPayKey = "C:/Users/hasee/Downloads/WXCertUtil/cert/new/apiclient_key.pem"
	//var tr *http.Transport
	// 微信提供的API证书,证书和证书密钥 .pem格式
	certs, err := tls.LoadX509KeyPair(wechatPayCert, wechatPayKey)
	if err != nil {
		log.Println("certs load err:", err)

	} else {
		// 微信支付HTTPS服务器证书的根证书  .pem格式
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				//RootCAs:      pool,
				Certificates: []tls.Certificate{certs},
			},
		}

		client := &http.Client{Transport: tr}

		req, createErr := http.NewRequest("POST", url, bytes.NewBuffer([]byte(byteData)))
		if createErr != nil {
			fmt.Printf("创建失败:%v\n",createErr)
		}

		req.Header.Set("Content-Type", "application/json")


		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		body,readErr := ioutil.ReadAll(resp.Body)
		fmt.Println(readErr)
		fmt.Println(string(body))

	}


/*	httpClient := &util.HttpClient{}
	err := httpClient.Post(url,byteData,nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		fmt.Println(readErr)
		fmt.Println(string(body))
	})

	fmt.Println(err)*/
}
