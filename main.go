package main

import (
	"fmt"
	"luck_draw/queue"
	"luck_draw/routers"
	"luck_draw/util"
	"net/http"
	"time"
)

func main() {
	router:= routers.InitRouter()

	fmt.Println(len(fmt.Sprintf("%v%v",time.Now().UnixNano(),3320)))

	/*ju,billErr := service.JuHePhoneBill("13425144866","123",float64(1))
	fmt.Println(ju)
	fmt.Println(billErr)*/

	go queue.Listen()
	go queue.ScanActivity()

	server := &http.Server{
		Addr:           ":8081",
		Handler:        router,
	}

	err := server.ListenAndServe()
	if err != nil {
		util.Error(fmt.Sprintf("启动服务失败：%v\n",err))
	}
}
