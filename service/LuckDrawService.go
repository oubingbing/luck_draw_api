package service

import (
	"luck_draw/model"
	"time"
)

func SaveLuckDraw(param *model.LuckDrawParam) (int64,error) {
	luckDraw := &model.LuckDraw{
		Name:param.Name,
		Type:model.LUCK_DRAW_TYPE_RED_PAK,
		FROM:model.LUCK_DRAW_FROM_USER,
		LimitJoin:param.LimitJoin,
		JoinLimitNum:param.JoinLimitNum,
		ReceiveLimit:param.ReceiveLimit,
		Describe:param.Describe,
		Attachments:param.Attachments,
		ShareTitle:param.ShareTitle,
		ShareImage:param.ShareImage,
	}

	var parseErr error
	luckDraw.StartAt,parseErr = time.Parse("2006-01-02 15:04:05",param.StartAt)
	if parseErr != nil {

	}

	luckDraw.EndAt,parseErr = time.Parse("2006-01-02 15:04:05",param.EndAt)
	if parseErr != nil {

	}

	luckDraw.RunAt,parseErr = time.Parse("2006-01-02 15:04:05",param.RunAt)
	if parseErr != nil {

	}

	effect,err := luckDraw.Store()
	return effect,err
}