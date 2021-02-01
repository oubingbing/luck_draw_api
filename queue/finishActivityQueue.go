package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
	"math/rand"
	"time"
)

func ScanActivity()  {
	for  {
		FinishRunningActivity()
		time.Sleep(time.Minute)
	}
}

func FinishRunningActivity()  {
	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ErrDetail(connectErr.Code,"完结活动时数据库连接错误",connectErr.Err)
		return
	}

	activity := &model.Activity{}
	data,err := activity.RunningActivity(db);
	if err != nil {
		util.ErrDetail(enums.ACTIVITY_GET_RUNNING_ERR,"完结活动时获取数据错误",err.Error())
		return
	}

	for _,item := range data {
		if float32(item.JoinNum) >= item.JoinLimitNum {
			switch item.Type {
			case model.ACTIVITY_TYPE_RED_PAK:
				go HandleReaPackage(item)
				break
			case model.ACTIVITY_TYPE_GOODS:
				go HandleGift(item)
				break
			case model.ACTIVITY_TYPE_GAME:
				go HandleGift(item)
				break
			case model.ACTIVITY_TYPE_PHONE_BILL:
				go HandlePhoneBill(item)
				break
			default:
				util.ErrDetail(enums.ACTIVITY_DEAL_NOT_HANDLE,"完结活动时获取数据错误",item.ID)
			}
		}
	}

}

//处理话费
func HandlePhoneBill(activity model.Activity)  {
	if float32(activity.JoinNum) < activity.JoinLimitNum {
		return
	}

	db,connectErr := model.Connect()
	if connectErr != nil {
		util.ErrDetail(connectErr.Code,"完结活动时数据库连接错误",connectErr.Err)
		return
	}
	
	redis := util.NewRedis()
	defer func() {
		redis.Client.Close()
		db.Close()
	}()

	joinLog := &model.JoinLog{}
	joinLogSli,err := joinLog.GetJoinLogByActivityId(db,activity.ID)
	if err != nil {
		util.ErrDetail(enums.ACTIVITY_JOIN_LOG_QUERY_ERR,"取出需要完结的活动参加记录时发生错误",err.Error())
		return
	}

	//将活动变更为已完成
	ac := &model.Activity{}
	updateActivity := make(map[string]interface{})
	updateActivity["status"] = model.ACTIVITY_STATSUS_FINISH
	updateAcErr := ac.Update(db,activity.ID,updateActivity)
	if updateAcErr != nil {
		util.ErrDetail(enums.ACTIVITY_FINDISH_DB_ERR,"活动变更为已完成数据库出错",activity.ID)
		return
	}

	//查找gift
	gift,giftErr := service.FirstGiftById(db,activity.GiftId)
	if giftErr != nil {
		util.ErrDetail(giftErr.Code,"取出需要完结的活动奖品时发生错误",giftErr.Err.Error())
		return
	}

	var ctx = context.Background()
	var consume int64 = 0
	if activity.DrawType == model.ACTIVITY_DRAW_TYPE_AVERAGE {
		//平均，人人有份
		averge := gift.Num / float32(activity.JoinNum)
		//话费的区间
		bill := []int{1,2,5,10}
		avergeBill := 1 //需要送的话费
		for _,item := range bill {
			if averge <= float32(item) {
				avergeBill = item
				break
			}
		}

		for _,item := range joinLogSli {
			inbox := model.InboxMessage{}
			inbox.UserId = item.UserId
			inbox.Bill = float64(avergeBill)
			inbox.JoinLogId = int64(item.ID)
			inbox.ObjectId = item.ActivityId
			inbox.ActivityName = activity.Name
			inbox.OrderId = item.OrderId
			inbox.Content = fmt.Sprintf("恭喜您，在活动%v中获得%v元话费，稍后将会充值到您的账户中，请留意手机短信消息",activity.Name,avergeBill)
			mpStr,_ := json.Marshal(&inbox)
			consume += int64(avergeBill)

			//推送到队列
			intCmd := redis.Client.LPush(ctx,enums.INBOX_QUEUE,string(mpStr))
			intCmd = redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_PHONE_BILL_QUEUE,string(mpStr))
			if intCmd.Err() != nil {
				util.ErrDetail(enums.ACTIVITY_PUSH_BILL_QUEUE_ERR,fmt.Sprintf("推送到话费发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,item.UserId),intCmd.Err().Error())
			}

			joinLog := &model.JoinLog{}
			update := make(map[string]interface{})
			update["remark"] = inbox.Content
			update["status"] = model.JOIN_LOG_STATUS_WIN
			update["num"]    = inbox.Bill
			updateErr := joinLog.Update(db,uint(inbox.JoinLogId),update)
			if updateErr != nil {
				util.ErrDetail(enums.ACTIVITY_UPDATE_JL_ERR,"跟新用户中奖join log数据库异常",updateErr.Error())
			}
		}
		//更新所有的
	}else{
		//拼手气,人人有份
		var avergeBill int64 = 1 //需要送的话费
		user := make(map[int]*model.InboxMessage)
		for index,item := range joinLogSli {
			inbox := model.InboxMessage{}
			inbox.UserId = item.UserId
			inbox.JoinLogId = int64(item.ID)
			inbox.Bill = float64(avergeBill)
			inbox.ObjectId = item.ActivityId
			inbox.ActivityName = activity.Name
			inbox.OrderId = item.OrderId
			inbox.Content = ""
			consume += avergeBill
			user[index] = &inbox
		}

		num := len(joinLogSli) //中奖人数
		leftAmount := gift.Num - float32(num)
		if leftAmount >= 1 {
			//循环扣减,直到奖金池为0
			seed := 1
			for  {
				if leftAmount <= 0 {
					break
				}
				rand.Seed(time.Now().UnixNano()+int64(seed))
				key := rand.Intn(num)
				//抽取一个中奖用户
				_,ok := user[key]
				if ok {
					user[key].Bill += 1
					leftAmount --
					consume += 1
				}
			}
		}

		for _,v := range user {
			v.Content = fmt.Sprintf("恭喜您，在%v中获得%v元话费，稍后将会充值到您的账户中，请留意手机短信消息",activity.Name,v.Bill)
			mpStr,_ := json.Marshal(&v)
			//推送到队列
			intCmd := redis.Client.LPush(ctx,enums.INBOX_QUEUE,string(mpStr))
			intCmd = redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_PHONE_BILL_QUEUE,string(mpStr))
			if intCmd.Err() != nil {
				util.ErrDetail(enums.ACTIVITY_PUSH_BILL_QUEUE_ERR,fmt.Sprintf("推送到话费发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,v.UserId),intCmd.Err().Error())
			}

			joinLog := &model.JoinLog{}
			update := make(map[string]interface{})
			update["remark"] = v.Content
			update["status"] = model.JOIN_LOG_STATUS_WIN
			update["num"]    = v.Bill
			updateErr := joinLog.Update(db,uint(v.JoinLogId),update)
			if updateErr != nil {
				util.ErrDetail(enums.ACTIVITY_UPDATE_JL_ERR,"跟新用户中奖join log数据库异常",updateErr.Error())
			}
		}

	}

	//还有一种是真用户只有1元，假用户可以中更多

	//更新活动实际消耗奖品数量
	updateConsueme := make(map[string]interface{})
	updateConsueme["consume"] = consume
	updateConsumeErr := ac.Update(db,activity.ID,updateConsueme)
	if updateConsumeErr != nil {
		util.ErrDetail(enums.ACTIVITY_UPDATE_CONSUME_DB_ERR,"活动更新实际消耗奖品数量出错",activity.ID)
		return
	}
}

//处理红包
func HandleReaPackage(activity model.Activity)  {
	if float32(activity.JoinNum) < activity.JoinLimitNum {
		return
	}

	redis := util.NewRedis()
	defer redis.Client.Close()
}

//处理礼品
func HandleGift(activity model.Activity)  {
	if float32(activity.JoinNum) < activity.JoinLimitNum {
		return
	}

	db,connectErr := model.Connect()
	if connectErr != nil {
		util.ErrDetail(connectErr.Code,"完结活动时数据库连接错误",connectErr.Err)
		return
	}

	redis := util.NewRedis()
	defer func() {
		redis.Client.Close()
		db.Close()
	}()

	joinLog := &model.JoinLog{}
	joinLogSli,err := joinLog.GetJoinLogByActivityId(db,activity.ID)
	if err != nil {
		util.ErrDetail(enums.ACTIVITY_JOIN_LOG_QUERY_ERR,"取出需要完结的活动参加记录时发生错误",err.Error())
		return
	}

	//查找gift
	gift,giftErr := service.FirstGiftById(db,activity.GiftId)
	if giftErr != nil {
		util.ErrDetail(giftErr.Code,"取出需要完结的活动奖品时发生错误",giftErr.Err.Error())
		return
	}

	//将活动变更为已完成
	ac := &model.Activity{}
	updateActivity := make(map[string]interface{})
	updateActivity["status"] = model.ACTIVITY_STATSUS_FINISH
	updateAcErr := ac.Update(db,activity.ID,updateActivity)
	if updateAcErr != nil {
		util.ErrDetail(enums.ACTIVITY_FINDISH_DB_ERR,"活动变更为已完成数据库出错",activity.ID)
		return
	}

	var winId []int64

	//更新未中奖的
	joinLogNot := &model.JoinLog{}
	update := make(map[string]interface{})
	update["remark"] = "很遗憾，您与大奖擦肩而过，请参加其他活动争取把大奖领回家吧，加油！"
	update["status"] = model.JOIN_LOG_STATUS_LOSE
	updateActivity["num"] = float64(0)
	updateErr := joinLogNot.UpdateNotWin(db,activity.ID,winId,update)
	if updateErr != nil {
		util.ErrDetail(enums.ACTIVITY_UPDATE_JL_ERR,"更新用户未中奖join log数据库异常",updateErr.Error())
	}

	var ctx = context.Background()
	var consume int64 = 0
	if activity.DrawType == model.ACTIVITY_DRAW_TYPE_AVERAGE {
		//人人有份
		/*for _,item := range joinLogSli {
			mp := make(map[string]interface{})
			mp["user_id"] = item.UserId
			mp["bill"] = 1
			mp["join_log_id"] = item.ID
			mpStr,_ := json.Marshal(&mp)
			consume += 1

			//推送到队列
			intCmd := redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_PHONE_BILL_QUEUE,string(mpStr))
			if intCmd.Err() != nil {
				util.ErrDetail(enums.ACTIVITY_PUSH_BILL_QUEUE_ERR,fmt.Sprintf("推送到话费发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,item.UserId),intCmd.Err().Error())
			}
		}*/
	}else{
		//拼手气
		user := make(map[int]*model.InboxMessage)
		var loseUser []*model.InboxMessage
		if activity.Really ==  model.ACTIVITY_REALLY_Y {
			//真送
			for index,item := range joinLogSli {
				if item.Faker == model.FAKER_N {
					inbox := model.InboxMessage{}
					inbox.UserId = item.UserId
					inbox.JoinLogId = int64(item.ID)
					inbox.Bill = 1
					inbox.ObjectId = item.ActivityId
					inbox.ActivityName = activity.Name
					inbox.Content = "很遗憾，您与大奖擦肩而过，请参加其他活动争取把大奖领回家吧，加油！"
					user[index] = &inbox
				}
			}
		}else{
			//假送
			for index,item := range joinLogSli {
				if item.Faker == model.FAKER_Y {
					inbox := model.InboxMessage{}
					inbox.UserId = item.UserId
					inbox.JoinLogId = int64(item.ID)
					inbox.Bill = 1
					inbox.ObjectId = item.ActivityId
					inbox.ActivityName = activity.Name
					inbox.Content = "很遗憾，您与大奖擦肩而过，请参加其他活动争取把大奖领回家吧，加油！"
					user[index] = &inbox
				}else{
					inbox := model.InboxMessage{}
					inbox.UserId = item.UserId
					inbox.JoinLogId = int64(item.ID)
					inbox.Bill = 1
					inbox.ObjectId = item.ActivityId
					inbox.ActivityName = activity.Name
					inbox.Content = "很遗憾，您与大奖擦肩而过，请参加其他活动争取把大奖领回家吧，加油！"
					loseUser = append(loseUser,&inbox)
				}
			}
		}

		num := len(user) //中奖人数
		leftAmount := activity.ReceiveLimit
		fmt.Println(leftAmount)
		i := 1
		if leftAmount >= 1 {
			//循环扣减,直到奖金池为0
			for  {
				if leftAmount <= 0 {
					break
				}
				rand.Seed(time.Now().UnixNano()+int64(i))
				i++
				key := rand.Intn(num)
				//抽取一个中奖用户
				v,ok := user[key]
				if ok {
					leftAmount --
					consume += 1

					user[key].Content = fmt.Sprintf("恭喜您在活动 %v 中获得 %v X1，请确保您已填写了收货地址，我们会往您的默认收货地址寄送奖品，谢谢。",activity.Name,gift.Name)
					util.Info(fmt.Sprintf("中奖用户：%v,中奖活动：%v，中奖内容：%v",user[key].UserId,user[key].ObjectId,user[key].Content))

					joinLog := &model.JoinLog{}
					update := make(map[string]interface{})
					update["remark"] = user[key].Content
					update["status"] = model.JOIN_LOG_STATUS_WIN
					update["num"] 	 = 1
					updateErr := joinLog.Update(db,uint(user[key].JoinLogId),update)
					if updateErr != nil {
						util.ErrDetail(enums.ACTIVITY_UPDATE_JL_ERR,"跟新用户中奖join log数据库异常",updateErr.Error())
					}

					winId = append(winId,user[key].JoinLogId)

					//通知中奖的
					mpStr,_ := json.Marshal(user[key])
					intCmd := redis.Client.LPush(ctx,enums.INBOX_QUEUE,string(mpStr))
					//推送到队列
					intCmd = redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_GIFT_QUEUE,string(mpStr))
					if intCmd.Err() != nil {
						util.ErrDetail(enums.ACTIVITY_PUSH_GIFT_QUEUE_ERR,fmt.Sprintf("推送到物品发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,v.UserId),intCmd.Err().Error())
					}
					delete(user,key)
				}
			}
		}

		//通知未中奖的
		for i,_ := range loseUser {
			mpStr,_ := json.Marshal(loseUser[i])
			redis.Client.LPush(ctx,enums.INBOX_QUEUE,string(mpStr))
		}
	}

	//更新活动实际消耗奖品数量
	updateConsueme := make(map[string]interface{})
	updateConsueme["consume"] = consume
	updateConsumeErr := ac.Update(db,activity.ID,updateConsueme)
	if updateConsumeErr != nil {
		util.ErrDetail(enums.ACTIVITY_UPDATE_CONSUME_DB_ERR,"活动更新实际消耗奖品数量出错",activity.ID)
		return
	}

}