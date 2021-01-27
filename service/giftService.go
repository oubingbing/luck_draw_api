package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
)

var giftNotFound error = errors.New("礼品不存在")
var giftSaveErr error = errors.New("数据异常，保存失败")

func SaveGift(db *gorm.DB,userId int,giftParam *enums.GiftParam) (int64,*enums.ErrorInfo) {
	gift := &model.Gift{
		Name:giftParam.Name,
		Num:giftParam.Num,
		UserId:userId,
		Type:giftParam.Type,
		FROM:giftParam.FROM,
		STATUS:giftParam.STATUS,
		Des:giftParam.Des,
		Attachments:giftParam.Attachments,
	}

	effect,err := gift.Store(db)
	if err != nil {
		return effect,&enums.ErrorInfo{giftSaveErr,enums.GIFT_SAVE_ERR}
	}

	return effect,nil
}

func FirstGiftById(db *gorm.DB,id int64) (*enums.GiftDetail,*enums.ErrorInfo) {
	gift := &model.Gift{}
	detail,err := gift.First(db,id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,&enums.ErrorInfo{giftNotFound,enums.GIFT_NOT_FOUND}
		}
		return nil,&enums.ErrorInfo{err,enums.GIFT_FIRST_ERR}
	}

	return detail,nil
}