package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"luck_draw/enums"
)

type Address struct {
	gorm.Model
	UserId  		uint		`gorm:"column:user_id"`
	Receiver 		string		`gorm:"column:receiver"`
	Phone 			string		`gorm:"column:phone"`
	Nation 			string		`gorm:"column:nation"`
	Province 		string		`gorm:"column:province"`
	City 			string		`gorm:"column:city"`
	District 		string		`gorm:"column:district"`
	DetailAddress 	string		`gorm:"column:detail_address"`
	UseType 		int8		`gorm:"column:use_type"`
}

type City struct {
	ID        		uint 		`gorm:"primary_key"`
	Name  			string		`gorm:"column:name"`
	CityId	 		string		`gorm:"column:city_id"`
	ProvinceId 		string		`gorm:"column:province_id"`
}

type Province struct {
	ID        		uint 		`gorm:"primary_key"`
	Name  			string		`gorm:"column:name"`
	ProvinceId 		string		`gorm:"column:province_id"`
}

type Country struct {
	ID        		uint 		`gorm:"primary_key"`
	Name  			string		`gorm:"column:name"`
	CountryId 		string		`gorm:"column:country_id"`
	CityId	 		string		`gorm:"column:city_id"`
}

type AddressPageSli []enums.AddressPage

func (Address) TableName() string  {
	return "address"
}

func (address *Address)Store(db *gorm.DB) (int64,error) {
	createResult := db.Create(address)
	return createResult.RowsAffected,createResult.Error
}

func (address *Address) FindById(db *gorm.DB,id interface{}) error {
	err := db.Table(address.TableName()).Where("id = ?",id).First(address).Error
	return err
}

func (address *Address) GetAddressInfo(db *gorm.DB) (map[string]interface{},error) {
	var country []Country
	var city []City
	var province []Province
	var err error

	err = db.Table("country").Find(&country).Error
	if err != nil{
		return nil,err
	}

	err = db.Table("city").Find(&city).Error
	if err != nil{
		return nil,err
	}

	err = db.Table("province").Find(&province).Error
	if err != nil{
		return nil,err
	}

	data := make(map[string]interface{})
	data["country"] = country
	data["city"] = city
	data["province"] = province

	return data,err
}

func (address *Address)Page(db *gorm.DB,userId interface{},page *PageParam) (*AddressPageSli,error) {
	var pageData AddressPageSli
	err :=  Page(db,address.TableName(),page).
		Where("user_id = ?",userId).
		Select("id,user_id,receiver,phone,province,city,district,use_type,detail_address").
		Order(fmt.Sprintf("%v %v",page.OrderBY,page.Sort)).
		Find(&pageData).Error
	if err != nil {
		return nil,err
	}

	return &pageData,nil
}