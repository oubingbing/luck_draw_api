package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
	"luck_draw/model"
	"luck_draw/util"
	"time"
)

func StoreAddress(db *gorm.DB,userId interface{},param *enums.AddressParam) (*model.Address,*enums.ErrorInfo) {
	uid,ok := userId.(float64)
	if !ok {
		util.ErrDetail(enums.SYSTEM_ERR,"userid断言失败",uid)
		return nil,&enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
	}

	address := &model.Address{
		UserId:uint(uid),
		Receiver:param.Receiver,
		Phone:param.Phone,
		Nation:param.Nation,
		Province:param.Province,
		City:param.City,
		District:param.District,
		DetailAddress:param.DetailAddress,
		UseType:param.UseType,
	}

	affect,err := address.Store(db)
	if err != nil {
		paramStr,_ := json.Marshal(param)
		util.ErrDetail(enums.ADDRESS_STORE_ERR,fmt.Sprintf("保存地址失败：%v",err.Error()),string(paramStr))
		return nil,&enums.ErrorInfo{enums.AddressStoreErr,enums.ADDRESS_STORE_ERR}
	}
	fmt.Println(affect)
	if affect <= 0 {
		paramStr,_ := json.Marshal(param)
		util.ErrDetail(enums.ADDRESS_STORE_AFFECT_ERR,fmt.Sprintf("保存地址失败：%v",affect),string(paramStr))
		return nil,&enums.ErrorInfo{enums.AddressStoreErr,enums.ADDRESS_STORE_AFFECT_ERR}
	}

	return address,nil
}

func UpdateAddress(db *gorm.DB,userId interface{},param *enums.AddressUpdateParam) (*model.Address,*enums.ErrorInfo) {
	address := &model.Address{}
	err := address.FindById(db,param.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,&enums.ErrorInfo{enums.AddressNotFound,enums.ADDRESS_NOT_FOUND}
		}else{
			return nil,&enums.ErrorInfo{enums.SystemErr,enums.ADDRESS_FIND_ERR}
		}
	}

	if float64(address.UserId) != userId.(float64) {
		fmt.Println(interface{}(address.UserId)==interface{}(userId))
		return nil,&enums.ErrorInfo{enums.AddressNotFound,enums.ADDRESS_NOT_FOUND}
	}

	address.Receiver 		= param.Receiver
	address.Phone 			= param.Phone
	address.Nation 			= param.Nation
	address.Province 		= param.Province
	address.City 			= param.City
	address.District 		= param.District
	address.DetailAddress 	= param.DetailAddress
	address.UseType 		= param.UseType
	updateResultErr := db.Save(address).Error
	if updateResultErr != nil {
		return nil,&enums.ErrorInfo{enums.AddressUpdateFail,enums.ADDRESS_UPDATE_ERR}
	}

	return address,nil
}

func GetAddressInfo() (map[string]interface{},*enums.ErrorInfo) {
	var err error
	var data map[string]interface{}

	key := "luck_draw_address"
	redis := util.NewRedis()
	defer redis.Client.Close()
	ctx := context.Background()
	exitKey := redis.Client.Get(ctx,key)

	if len(exitKey.Val()) <= 0 {
		fmt.Println("不走缓存")
		db,connectErr := model.Connect()
		if connectErr != nil {
			return nil,&enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
		}

		address := &model.Address{}
		data,err = address.GetAddressInfo(db)
		db.Close()
		if err != nil {
			return nil,&enums.ErrorInfo{enums.AddressListQueryERr,enums.ADDRESS_LIST_QUERY_ERR}
		}

		jsonStr,err := json.Marshal(&data)
		if err != nil {
			util.Error("地址信息转成json字符串出错")
			return nil,&enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
		}

		redis.Client.Set(ctx,key,string(jsonStr),time.Hour*24*10)
	}else{
		fmt.Println("走了缓存")
		err = json.Unmarshal([]byte(exitKey.Val()),&data)
		if err != nil {
			util.Error("地址信息json字符串转map出错")
			return nil,&enums.ErrorInfo{enums.SystemErr,enums.SYSTEM_ERR}
		}
	}

	return data,nil
}

func GetAddressPage(db *gorm.DB,userId interface{},page *model.PageParam) (*model.AddressPageSli,*enums.ErrorInfo) {
	address := &model.Address{}
	pageData,err := address.Page(db,userId,page)
	if err != nil {
		return nil,&enums.ErrorInfo{enums.AddressPageQueryERr,enums.ADDRESS_PAGE_QUERY_ERR}
	}
	return pageData,nil
}
