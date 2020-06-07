package model

import "github.com/jinzhu/gorm"

type Chat struct {
	*gorm.Model
	username string
	msg string
}
