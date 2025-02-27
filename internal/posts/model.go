package posts

import (
	"backend/blog/internal/comment"
	"backend/blog/internal/user"
	"time"

	"gorm.io/gorm"
)

// Модель поста
type Post struct {
	gorm.Model
	Title       string                 `gorm:"size:255;not null" json:"title"`             
	Content     string                 `gorm:"type:text;not null" json:"content"`          
	ImageURL    string                 `gorm:"size:512" json:"image_url"`                  
	UserID      uint                   `gorm:"not null" json:"-"`                    
	User        user.User              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"author"`
	Comments    []comment.CommentPost   `gorm:"foreignKey:PostID" json:"comments"`
	Published   bool                   `gorm:"default:false" json:"published"`             
	PublishedAt *time.Time             `json:"published_at,omitempty"` 
	View        uint                   `json:"view"`             
}