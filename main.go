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
	orderNo := service.RandChar(32)
	pay := service.WeChatPay{
		XMLName:        xml.Name{},
		PartnerTradeNo: orderNo,
		Openid:         "oc-wc48UGe_dVNqmAPJjIQAEw5-w",
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
