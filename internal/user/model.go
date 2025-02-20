package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"size:255;not null;unique" json:"username"`
	Email        string    `gorm:"size:255;not null;unique" json:"email"`
	Password		 string    `gorm:"size:255;not null" json:"-"`
	FirstName    string    `gorm:"size:255" json:"first_name"`
	LastName     string    `gorm:"size:255" json:"last_name"`
	AvatarURL    string    `gorm:"size:512" json:"avatar_url"`
	LastLogin    time.Time `json:"last_login"`
	IsAdmin      bool      `gorm:"default:false" json:"is_admin"`
	//IsActive     bool      `gorm:"default:true" json:"is_active"`
	//Posts        []Post    `gorm:"foreignKey:AuthorID" json:"posts"`
}