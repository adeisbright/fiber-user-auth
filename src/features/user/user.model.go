package user

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `json:"username"`
	Email    string `json:"email" gorm:"index"`
	Password string `json:"password"`
}
