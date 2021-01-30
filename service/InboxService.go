package service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/util"
	"time"
)

func SaveInbox(db *gorm.DB,inbox *model.Inbox) *enums.ErrorInfo {
	effect,err := inbox.Store(db)
	if err != nil {
		str,_ := json.Marshal(inbox)
		util.ErrDetail(enums.INBOX_CREATE_FAIL,fmt.Sprintf("消息保存失败：%v",err.Error()),string(str))
		return &enums.ErrorInfo{enums.InboxStoreErr,enums.INBOX_CREATE_FAIL}
	}

	if effect <= 0 {
		str,_ := json.Marshal(inbox)
		util.ErrDetail(enums.INBOX_CREATE_FAIL,fmt.Sprintf("消息保存失败"),string(str))
		return &enums.ErrorInfo{enums.InboxStoreErr,enums.INBOX_CREATE_FAIL}
	}

	return nil
}

func ReadInbox(db *gorm.DB,id interface{}) *enums.ErrorInfo {
	inbox := &model.Inbox{}
	data := make(map[string]interface{})
	data["read_at"] = time.Now().Format(enums.DATE_FORMAT)
	err := inbox.Update(db,id,data)
	if err != nil {
		return &enums.ErrorInfo{enums.InboxUpdateReadErr,enums.INBOX_UPDATE_READ_FAIL}
	}

	return nil
}

func GetInboxList(db *gorm.DB,userId interface{},page *model.PageParam) (model.InboxPageSli,*enums.ErrorInfo) {
	inbox := &model.Inbox{}
	list,err := inbox.Page(db,userId,page)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,nil
		}else{
			util.ErrDetail(enums.INBOX_PAGE_QUERY_FAIL,"查询消息盒子失败",enums.InboxPageQueryErr.Error())
			return nil,&enums.ErrorInfo{enums.InboxPageQueryErr,enums.INBOX_PAGE_QUERY_FAIL}
		}
	}

	config,_ := util.GetConfig()
	domain := config["COS_DOMAIN"]

	var ParseErr *enums.ErrorInfo
	for index,_ := range list {
		list[index].AttachmentsSli,ParseErr = AppendDomain(domain,list[index].Attachments)
		if ParseErr != nil {
			return nil,ParseErr
		}
		list[index].Attachments = ""
	}

	return list,nil
}

func CountInboxUnRead(db *gorm.DB,userId interface{}) (int,*enums.ErrorInfo) {
	inbox := &model.Inbox{}
	count,err := inbox.CountUnRead(db,userId)
	if err != nil {
		return 0,&enums.ErrorInfo{enums.InboxCountQueryErr,enums.INBOX_COUNT_QUERY_FAIL}
	}

	return count,nil
}
