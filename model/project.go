package model

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	UserId 			int 	`gorm:"column:user_id"`
	Name 			string 	`gorm:"column:name"`
	Attachment 		string 	`gorm:"column:attachment"`
	Introduction 	string 	`gorm:"column:introduction"`
}
