package main

import (
	"encoding/xml"
	"fmt"
	"luck_draw/routers"
	"luck_draw/service"
	"luck_draw/util"
	"net/http"
)

func main() {
	router:= routers.InitRouter()
	pay := service.WeChatPay{
		XMLName:        xml.Name{},
		MchAppid:       "wxa0d7aa1607c5ac21",
		Mchid:          "1254223701",
		NonceStr:       "qweqweqweqwe",
		PartnerTradeNo: "10000098201411111234567890",
		Openid:         "oc-wc48UGe_dVNqmAPJjIQAEw5-w",
		CheckName:      "NO_CHECK",
		Amount:         10,
		Sign:           "",
		Desc:           "抽奖",
	}

	service.Pay(pay)

	//go queue.Listen()
	//go queue.ScanActivity()

	server := &http.Server{
		Addr:           ":8081",
		Handler:        router,
	}

	err := server.ListenAndServe()
	if err != nil {
		util.Error(fmt.Sprintf("启动服务失败：%v\n",err))
	}
}
