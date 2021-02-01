package queue

import (
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
	"sync"
	"time"
)

func AttemptJoin(db *gorm.DB,id interface{}) (string,int) {
	finish := 0
	msg := "参加失败，请重试"
	var userId int64

	defer func() {
		db.Close()
		notifyErr := service.SocketNotify(string(userId),finish,msg)
		if notifyErr != nil {
			util.Error(notifyErr.Error())
		}
		util.Info(fmt.Sprintf("%v,%v,%v",finish,msg,userId))
	}()

	//tx := db.Begin()
	joinLog := &model.JoinLog{}
	userId = joinLog.UserId
	err := joinLog.LockById(db,id)
	if err == gorm.ErrRecordNotFound {
		//tx.Rollback()
		finish = enums.ACTIVITY_DEAL_QUEUE_NOT_FOUND
		util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_NOT_FOUND,enums.ActivityQueueNotFound.Error(),id)
		return msg,finish
	}


	if joinLog.Status != model.JOIN_LOG_STATUS_QUEUE {
		finish = enums.ACTIVITY_STATUS_NOT_RUNNING
		//tx.Rollback()
		return msg,finish
	}

	activity := &model.Activity{}
	err = activity.FirstById(db,joinLog.ActivityId)
	if err == gorm.ErrRecordNotFound {
		//tx.Rollback()
		finish = enums.ACTIVITY_DEAL_QUEUE_A_NOT_FOUND
		util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_A_NOT_FOUND,enums.ActivityQueueANotFound.Error(),id)
		return msg,finish
	}

	//Faker join
	if int(activity.Really) == model.ACTIVITY_REALLY_N {
		fakerUserErr := service.JoinFakerUser(db,activity,userId)
		if fakerUserErr != nil {
			msg = fakerUserErr.Err.Error()
			finish = fakerUserErr.Code
			return msg,finish
		}
	}

	//重新查一次
	activity = &model.Activity{}
	err = activity.FirstById(db,joinLog.ActivityId)
	if err == gorm.ErrRecordNotFound {
		//tx.Rollback()
		finish = enums.ACTIVITY_DEAL_QUEUE_A_NOT_FOUND
		util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_A_NOT_FOUND,enums.ActivityQueueANotFound.Error(),id)
		return msg,finish
	}

	if float32(activity.JoinNum) >= activity.JoinLimitNum {
		data := make(map[string]interface{})
		data["remark"] = "人数已满"
		data["status"] = model.JOIN_LOG_STATUS_FAIL
		err := joinLog.Update(db,joinLog.ID,data)
		msg = "人数已满，下次抓紧机会啦"
		finish = enums.ACTIVITY_MEMBER_ENOUTH
		if err != nil {
			//tx.Rollback()
			util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_UPDATE_LOG_ERR,enums.ActivityJoinLogUpdateFailErr.Error(),id)
			return msg,finish
		}
	}

	data := make(map[string]interface{})
	data["remark"] = "加入成功"
	data["status"] = model.JOIN_LOG_STATUS_SUCCESS
	data["joined_at"] = time.Now().Format("2006-01-02 15:04:05")
	err = joinLog.Update(db,joinLog.ID,data)
	if err != nil {
		//tx.Rollback()
		finish = enums.ACTIVITY_DEAL_QUEUE_UPDATE_LOG_ERR
		util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_UPDATE_LOG_ERR,enums.ActivityJoinLogUpdateFailErr.Error(),id)
		return msg,finish
	}

	activityData := make(map[string]interface{})
	activityData["join_num"] = activity.JoinNum+1
	err = activity.Update(db,activity.ID,activityData)
	if err != nil {
		//tx.Rollback()
		finish = enums.ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR
		util.ErrDetail(enums.ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR,enums.ActivityUpdateJoinNumFailErr.Error(),id)
		return msg,finish
	}

	//tx.Commit()
	msg = "加入成功"

	return msg,finish
}

/**
 * 监听redis
 */
func ListenAttemptJoin(wg *sync.WaitGroup)  {
	redis := util.NewRedis()
	t := time.Second * 59

	queue := enums.ACTIVITY_QUEUE
	redis.OnQueue(wg,queue,t, func(result *redis2.StringSliceCmd, e error) {
		if len(result.Val()) > 0 {
			db,connectErr := model.Connect()
			defer db.Close()
			if connectErr != nil {
				//丢到重试
				return
			}

			util.Info(fmt.Sprintf("取出加入活动的log ID：%v",result.Val()[1]));
			AttemptJoin(db,result.Val()[1])
		}
	})
}

/**
 * 监听发送话费队列
 */
func ListenPhoneBill(wg *sync.WaitGroup)  {
	redis := util.NewRedis()
	t := time.Second * 59

	queue := enums.ACTIVITY_HANDLE_PHONE_BILL_QUEUE
	redis.OnQueue(wg,queue,t, func(result *redis2.StringSliceCmd, e error) {
		if len(result.Val()) > 0 {
			util.Info(fmt.Sprintf("取出需要发送话费的数据：%v",result.Val()[1]));
			HandleSendPhoneBill(result.Val()[1])
		}
	})
}

func Listen()  {
	var wg sync.WaitGroup
	wg.Add(1)
	go ListenInbox(&wg)
	go ListenAttemptJoin(&wg)
	go ListenPhoneBill(&wg)
	wg.Wait()
	//程序退出，需要通知开发人员
}
