package comment

import (
	"backend/blog/internal/user"

	"gorm.io/gorm"
)

type CommentPost struct {
	gorm.Model
	Text    string    `gorm:"type:text;not null" json:"text"`
	UserID  uint      `gorm:"not null" json:"-"`
	User    user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"author_comment"`
	PostID  uint      `gorm:"not null" json:"-"`
}