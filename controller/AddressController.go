package controller

import (
	"github.com/gin-gonic/gin"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/service"
	"luck_draw/util"
	"regexp"
)

func CreateAddress(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")

	var param enums.AddressParam
	errInfo := &enums.ErrorInfo{}
	if errInfo.Err = ctx.ShouldBind(&param); errInfo.Err != nil {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,errInfo.Err.Error(),nil)
		return
	}

	reg := `^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\d{8}$`
	rgx := regexp.MustCompile(reg)
	if phoneCheck := rgx.MatchString(param.Phone); !phoneCheck {
		util.ResponseJson(ctx,enums.ADDRESS_FORMAT_ERR,enums.AddressPhoneErr.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	address,err := service.StoreAddress(db,uid,&param)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",address)
	return
}

func UpdateAddress(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	var param enums.AddressUpdateParam
	errInfo := &enums.ErrorInfo{}
	if errInfo.Err = ctx.ShouldBind(&param); errInfo.Err != nil {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,errInfo.Err.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	address,err := service.UpdateAddress(db,uid,&param)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",address)
	return
}

func GetAddressList(ctx *gin.Context)  {
	data,err := service.GetAddressInfo()
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",data)
	return
}

func GetAddressPage(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	var param model.PageParam
	errInfo := &enums.ErrorInfo{}
	if errInfo.Err = ctx.ShouldBind(&param); errInfo.Err != nil {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,errInfo.Err.Error(),nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	page,err := service.GetAddressPage(db,uid,&param)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",page)
	return
}

func GetAddressDetail(ctx *gin.Context)  {
	id,ok := ctx.GetQuery("id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,"id不能为空",nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	address,err := service.AddressDetail(db,id)
	if err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	var mp map[string]interface{} = map[string]interface{}{
		"id":address.ID,
		"receiver":address.Receiver,
		"phone":address.Phone,
		"province":address.Province,
		"city":address.City,
		"district":address.District,
		"detail_address":address.DetailAddress,
		"useType":address.UseType,
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",mp)
	return
}

func DeleteAddress(ctx *gin.Context)  {
	uid,_:= ctx.Get("user_id")
	id,ok := util.Input(ctx,"id")
	if !ok {
		util.ResponseJson(ctx,enums.ACTIVITY_PARAM_ERR,"id不能为空",nil)
		return
	}

	db,connectErr := model.Connect()
	defer db.Close()
	if connectErr != nil {
		util.ResponseJson(ctx,connectErr.Code,connectErr.Err.Error(),nil)
		return
	}

	if err := service.DeleteAddress(db,uid,id);err != nil {
		util.ResponseJson(ctx,err.Code,err.Err.Error(),nil)
		return
	}

	util.ResponseJson(ctx,enums.SUCCESS,"",nil)
	return
}