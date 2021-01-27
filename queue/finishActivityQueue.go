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
		switch item.Type {
		case model.ACTIVITY_TYPE_RED_PAK:
			go HandleReaPackage(item)
			break
		case model.ACTIVITY_TYPE_GOODS:
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

//处理话费
func HandlePhoneBill(activity model.Activity,)  {
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

	//查找gift
	gift,giftErr := service.FirstGiftById(db,activity.GiftId)
	if giftErr != nil {
		util.ErrDetail(giftErr.Code,"取出需要完结的活动奖品时发生错误",giftErr.Err.Error())
		return
	}

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
			}
		}

		for _,item := range joinLogSli {
			mp := make(map[string]interface{})
			mp["user_id"] = item.UserId
			mp["bill"] = avergeBill
			mp["join_log_id"] = item.ID
			mpStr,_ := json.Marshal(&mp)
			consume += int64(avergeBill)

			//推送到队列
			intCmd := redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_PHONE_BILL_QUEUE,string(mpStr))
			if intCmd.Err() != nil {
				util.ErrDetail(enums.ACTIVITY_PUSH_BILL_QUEUE_ERR,fmt.Sprintf("推送到话费发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,item.UserId),intCmd.Err().Error())
			}
		}
		//更新所有的
	}else{
		//拼手气,人人有份
		var avergeBill int64 = 1 //需要送的话费
		user := make(map[int]map[string]int64)
		for index,item := range joinLogSli {
			mp := make(map[string]int64)
			mp["user_id"] = item.UserId
			mp["join_log_id"] = int64(item.ID)
			mp["bill"] = avergeBill
			consume += avergeBill
			user[index] = mp
		}

		num := len(joinLogSli) //中奖人数
		leftAmount := gift.Num - float32(num)
		if leftAmount >= 1 {
			//循环扣减,直到奖金池为0
			for  {
				if leftAmount <= 0 {
					break
				}
				key := rand.Intn(num)
				//抽取一个中奖用户
				_,ok := user[key]
				if ok {
					user[key]["bill"] = user[key]["bill"]+1
					leftAmount --
					consume += 1
				}
			}
		}

		for _,v := range user {
			mpStr,_ := json.Marshal(&v)
			//推送到队列
			intCmd := redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_PHONE_BILL_QUEUE,string(mpStr))
			if intCmd.Err() != nil {
				util.ErrDetail(enums.ACTIVITY_PUSH_BILL_QUEUE_ERR,fmt.Sprintf("推送到话费发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,v["user_id"]),intCmd.Err().Error())
			}
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

//处理红包
func HandleReaPackage(activity model.Activity)  {
	redis := util.NewRedis()
	defer redis.Client.Close()
}

//处理礼品
func HandleGift(activity model.Activity)  {
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

	//查找gift
	gift,giftErr := service.FirstGiftById(db,activity.GiftId)
	if giftErr != nil {
		util.ErrDetail(giftErr.Code,"取出需要完结的活动奖品时发生错误",giftErr.Err.Error())
		return
	}

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
		user := make(map[int]map[string]int64)
		if activity.Really ==  model.ACTIVITY_REALLY_Y {
			//真送
			for index,item := range joinLogSli {
				if item.Faker == model.FAKER_N {
					mp := make(map[string]int64)
					mp["user_id"] = item.UserId
					mp["join_log_id"] = int64(item.ID)
					mp["bill"] = 1
					user[index] = mp
				}
			}
		}else{
			//假送
			for index,item := range joinLogSli {
				if item.Faker == model.FAKER_Y {
					mp := make(map[string]int64)
					mp["user_id"] = item.UserId
					mp["join_log_id"] = int64(item.ID)
					mp["bill"] = 1
					user[index] = mp
				}
			}
		}

		num := len(user) //中奖人数
		leftAmount := gift.Num
		if leftAmount >= 1 {
			//循环扣减,直到奖金池为0
			for  {
				if leftAmount <= 0 {
					break
				}
				key := rand.Intn(num)
				//抽取一个中奖用户
				v,ok := user[key]
				if ok {
					leftAmount --
					consume += 1

					mpStr,_ := json.Marshal(&v)
					//推送到队列
					intCmd := redis.Client.LPush(ctx,enums.ACTIVITY_HANDLE_GIFT_QUEUE,string(mpStr))
					if intCmd.Err() != nil {
						util.ErrDetail(enums.ACTIVITY_PUSH_GIFT_QUEUE_ERR,fmt.Sprintf("推送到物品发货队列失败,acitivity_id:%v\n,user_id:%v",activity.ID,v["user_id"]),intCmd.Err().Error())
					}
					delete(user,key)
				}
			}
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