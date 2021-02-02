package queue

import (
	"encoding/json"
	"fmt"
	"luck_draw/service"
	"luck_draw/util"
)

func HandleWxNotify(data string)  {
	mp := make(map[string]string)
	err := json.Unmarshal([]byte(data),&mp)
	if err != nil {
		util.Error(fmt.Sprintf("解析微信消息通知错误：%v",data))
		return
	}

	if mp["type"] == "d" {
		//通知抽奖结果
		service.WxNotifyDraw(mp["id"],mp["openid"] ,mp["activityName"],mp["result"],mp["time"],mp["giftName"],mp["remark"])
	}else{
		//发奖励
		service.WxNotifyAward(mp["id"],mp["openid"] ,mp["giftName"],mp["activityName"],mp["time"],mp["remark"])
	}
}
