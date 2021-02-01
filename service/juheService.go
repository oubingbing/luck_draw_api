package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"luck_draw/enums"
	"luck_draw/util"
	"net/http"
)

type JuHeResponse struct {
	ErrorCode int64
	Reason	string
	Result map[string]interface{}
}

const PHONE_BILL_URL = "http://op.tianjurenhe.com/ofpay/mobile/onlineorder"

func JuHePhoneBill(phone ,orderId string,bill float64) (JuHeResponse,*enums.ErrorInfo) {
	var errInfo *enums.ErrorInfo
	var juHeResponse JuHeResponse

	config,_ := util.GetConfig()
	appKey := config["JU_APP_KEY"]
	openId := config["JU_OPEN_ID"]

	singStr := fmt.Sprintf("%v%v%v%v%v",openId,appKey,phone,bill,orderId)
	h := md5.New()
	h.Write([]byte(singStr))
	cipherStr := h.Sum(nil)
	sign := hex.EncodeToString(cipherStr)

	url := PHONE_BILL_URL+fmt.Sprintf("?phoneno=%v&cardnum=%v&orderid=%v&key=%v&sign=%v",phone,bill,orderId,appKey,sign)
	httpClient := util.HttpClient{}
	err := httpClient.Get(url,nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			errInfo = &enums.ErrorInfo{enums.GiftPhoneBillSendIOErr,enums.GIFT_PHONE_BILL_SEND_IO_ERR}
		}

		util.Info(fmt.Sprintf("话费发放结果：%v,order_id:%v",string(body),orderId))

		parseErr := json.Unmarshal(body,&juHeResponse)
		if parseErr != nil {
			errInfo = &enums.ErrorInfo{enums.GiftPhoneBillParseErr,enums.GIFT_PHONE_BILL_PARSE_ERR}
		}
	})

	if err != nil {
		errInfo = &enums.ErrorInfo{enums.GiftPhoneBillRequestErr,enums.GIFT_PHONE_BILL_REQUEST_ERR}
	}

	return juHeResponse,errInfo
}
