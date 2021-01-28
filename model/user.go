package model

import "github.com/jinzhu/gorm"

const (
	USER_FROM_MINI			= 1		//小程序
	USER_FRMO_OFFICIAL		= 2		//公众号
)

type User struct {
	gorm.Model
	NickName 		string		`gorm:"column:nick_name"`		//昵称
	AvatarUrl		string		`gorm:"column:avatar_url"`		//头像
	Gender			int8		`gorm:"column:gender"`			//性别
	OpenId			string		`gorm:"column:open_id"`			//openid
	UnionId			string		`gorm:"column:union_id"`
	City			string		`gorm:"column:city"`
	Country			string		`gorm:"column:country"`
	Language		string		`gorm:"column:language"`
	Province		string		`gorm:"column:province"`
	FromType		int8		`gorm:"column:from_type"`		//用户来源,1=小程序，2=h5公众号
	Phone			string		`gorm:"column:phone"`
	Faker			int8		`gorm:"column:faker"`
}

func (User) TableName() string  {
	return "wechat_user"
}

func (user *User)Store(db *gorm.DB) (int64,error) {
	result := db.Create(user)
	return result.RowsAffected,result.Error
}

func (user *User)FindByOpenId(db *gorm.DB,openId string) error {
	err := db.Table(user.TableName()).Where("deleted_at is null").Where("open_id = ?",openId).First(user).Error
	return err
}

func (user *User)Update(db *gorm.DB,id uint,data map[string]interface{}) error {
	err := db.Table(user.TableName()).Where("deleted_at is null").Where("id = ?",id).Updates(data).Error
	return err
}

func (user *User)FindById(db *gorm.DB,id int64) error {
	err := db.Table(user.TableName()).Where("deleted_at is null").Where("id = ?",id).First(user).Error
	return err
}