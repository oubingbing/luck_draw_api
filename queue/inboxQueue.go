package queue

import (
	"encoding/json"
	"fmt"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
	"sync"
	"time"
	redis2 "github.com/go-redis/redis/v8"
)

/**
 * 监听处理inbox
 */
func ListenInbox(wg *sync.WaitGroup)  {
	redis := util.NewRedis()
	t := time.Second * 59

	queue := enums.INBOX_QUEUE
	redis.OnQueue(wg,queue,t, func(result *redis2.StringSliceCmd, e error) {
		if len(result.Val()) > 0 {
			db,connectErr := model.Connect()
			defer db.Close()
			if connectErr != nil {
				//丢到重试
				fmt.Println("监听inbox出错")
				return
			}

			data := &model.InboxMessage{}
			err := json.Unmarshal([]byte(result.Val()[1]),data)
			if err != nil {
				util.Error(fmt.Sprintf("解析inbox数据失败:%v",result.Val()[1]))
				return
			}
			
			inbox := &model.Inbox{
				UserId:     data.UserId,
				ObjectType: 1,
				ObjectId:   data.ObjectId,
				Content:    data.Content,
				ReadAt:    nil,
			}
			saveErr := service.SaveInbox(db,inbox)
			if saveErr != nil {
				util.ErrDetail(saveErr.Code,saveErr.Err.Error(),result.Val()[1])
			}
		}
	})
}
