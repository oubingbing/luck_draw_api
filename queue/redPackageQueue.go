package queue

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
	"time"
)

func HandleRedPackage(inboxMessage string)  {
	db,connectErr := model.Connect()
	if connectErr != nil {
		//丢到重试
		return
	}

	curTime := time.Now().Format(enums.DATE_FORMAT)
	var ctx = context.Background()
	redis := util.NewRedis()
	defer func() {
		redis.Client.Close()
		db.Close()
	}()

	inbox := model.InboxMessage{}
	err := json.Unmarshal([]byte(inboxMessage),&inbox)
	if err != nil {
		util.Error(fmt.Sprintf("解析发放红包数据错误：%v",inboxMessage))
		return
	}

	user,getUserErr := service.FindUserById(db,inbox.UserId)
	if getUserErr != nil {
		util.ErrDetail(getUserErr.Code,"发放红包获取用户信息失败",getUserErr.Err.Error())
		return
	}

	status := int8(0)
	remark := ""
	if user.Faker == int8(model.FAKER_N) {
		//真用户，发放红包
		orderId := fmt.Sprintf("%v",inbox.OrderId)
		orderNo := service.RandChar(32)
		pay := service.WeChatPay{
			XMLName:        xml.Name{},
			PartnerTradeNo: orderNo,
			Openid:         user.OpenId,
			Amount:         int(inbox.Bill)*100,
			Sign:           "",
			Desc:           "金抽抽活动抽奖红包",
		}
		payResult := service.Pay(pay)
		payResultTxt := ""
		if payResult == "FAIL" {
			payResultTxt = "红包发放失败"
			remark = fmt.Sprintf("您抽到的%v红包发放失败，请联系客服，谢谢",inbox.Bill)
			util.Error(fmt.Sprintf("红包充值失败，手机号：%v，金额：%v，订单号：%v",user.Phone,inbox.Bill,orderId))
			status = int8(model.JOIN_LOG_SEND_AWARD_FAIL)
		}else{
			payResultTxt = "红包发放成功"
			remark = fmt.Sprintf("您的%v元红包发放成功啦，请留意微信消息",inbox.Bill)
			status = int8(model.JOIN_LOG_SEND_AWARD_SUCCESS)
			util.Info(fmt.Sprintf("红包充值成功，手机号：%v，金额：%v，订单号：%v",user.Phone,inbox.Bill,orderId))
		}

		PushPhoneBillInbox(redis,inbox,remark)

		mp := make(map[string]string)
		mp["type"] = "s"
		mp["id"] = fmt.Sprintf("%v",inbox.ObjectId)
		mp["openid"] = user.OpenId
		mp["activityName"] = inbox.ActivityName
		mp["time"] = curTime
		mp["giftName"] = fmt.Sprintf("%v元红包",inbox.Bill)
		mp["remark"] = payResultTxt
		mpStr,_ := json.Marshal(&mp)
		redis.Client.LPush(ctx,enums.WX_NOTIFY_QUEUE,string(mpStr))
	}else{
		//假用户，不发放
		util.Info(fmt.Sprintf("假用户不发放红包：%v",inboxMessage))
		response.ErrorCode = 0
	}

	joinLog := &model.JoinLog{}
	update := make(map[string]interface{})
	//update["remark"] = remark
	update["status"] = status
	updateErr := joinLog.Update(db,uint(inbox.JoinLogId),update)
	if updateErr != nil {
		util.ErrDetail(enums.ACTIVITY_UPDATE_JL_ERR,"发放红包跟新join log数据库异常",updateErr.Error())
	}
}
