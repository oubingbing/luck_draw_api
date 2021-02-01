package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
)

func HandleSendPhoneBill(inboxMessage string)  {
	db,connectErr := model.Connect()
	if connectErr != nil {
		//丢到重试
		return
	}

	redis := util.NewRedis()
	defer func() {
		redis.Client.Close()
		db.Close()
	}()

	inbox := model.InboxMessage{}
	err := json.Unmarshal([]byte(inboxMessage),&inbox)
	if err != nil {
		util.Error(fmt.Sprintf("解析发放话费数据错误：%v",inboxMessage))
		return
	}

	user,getUserErr := service.FindUserById(db,inbox.UserId)
	if getUserErr != nil {
		util.ErrDetail(getUserErr.Code,"发放话费获取用户信息失败",getUserErr.Err.Error())
		return
	}

	if len(user.Phone) <= 0 {
		util.Error(fmt.Sprintf("发放话费，用户手机号未绑定：%v",inboxMessage))
		return
	}

	response := service.JuHeResponse{}
	var responseErr *enums.ErrorInfo
	status := int8(0)
	remark := ""
	if user.Faker == int8(model.FAKER_N) {
		//真用户，发放话费

		billList := []int{1,2,5,10}
		sumBill := inbox.Bill
		if sumBill > 10 {
			//请联系客服充值
			PushPhoneBillInbox(redis,inbox,"话费金额较大，请联系客服充值，谢谢")
		}else{
			for {
				for i := 0; i <= len(billList) - 1; i++  {
					if sumBill <= float64(billList[i]) && sumBill > 0 {
						response,responseErr = service.JuHePhoneBill(user.Phone , inbox.OrderId,float64(billList[i]))
						if responseErr != nil {
							util.ErrDetail(responseErr.Code,"发放话费异常",responseErr.Err.Error())
							break
						}
						util.Info(fmt.Sprintf("话费充值完成,code:%v,reason:%v,phone:%v,金额：%v,订单号：%v",response.ErrorCode,response.Reason,user.Phone,billList[i],inbox.OrderId))
						if response.ErrorCode == int64(0) {
							remark = fmt.Sprintf("话费充值成功，请稍后确认是否充值到账，有任何问题请联系客服，谢谢")
							status = int8(model.JOIN_LOG_SEND_AWARD_SUCCESS)
							util.Info(fmt.Sprintf("话费充值成功，手机号：%v，金额：%v，订单号：%v",user.Phone,billList[i],inbox.OrderId))
						}else{
							remark = fmt.Sprintf("话费充值失败，请联系客服，谢谢")
							util.Error(fmt.Sprintf("话费充值失败，手机号：%v，金额：%v，订单号：%v",user.Phone,billList[i],inbox.OrderId))
							status = int8(model.JOIN_LOG_SEND_AWARD_FAIL)
						}

						PushPhoneBillInbox(redis,inbox,remark)

						sumBill -= float64(billList[i])
						break
					}
				}
				if sumBill <= 0 {
					break
				}
			}
		}

	}else{
		//假用户，不发放
		util.Info(fmt.Sprintf("假用户不发放话费：%v",inboxMessage))
		response.ErrorCode = 0
	}

	joinLog := &model.JoinLog{}
	update := make(map[string]interface{})
	//update["remark"] = remark
	update["status"] = status
	updateErr := joinLog.Update(db,uint(inbox.JoinLogId),update)
	if updateErr != nil {
		util.ErrDetail(enums.ACTIVITY_UPDATE_JL_ERR,"发放话费跟新join log数据库异常",updateErr.Error())
	}
}

func PushPhoneBillInbox(redis *util.MyRedis,inbox model.InboxMessage,message string)  {
	var ctx = context.Background()
	newInbox := model.InboxMessage{}
	newInbox.UserId = inbox.UserId
	newInbox.JoinLogId = inbox.JoinLogId
	newInbox.Bill = 0
	newInbox.ObjectId = inbox.ObjectId
	newInbox.ActivityName = inbox.ActivityName
	newInbox.OrderId = inbox.OrderId
	newInbox.Content = message

	mpStr,_ := json.Marshal(&newInbox)
	//推送到队列
	redis.Client.LPush(ctx,enums.INBOX_QUEUE,string(mpStr))
}