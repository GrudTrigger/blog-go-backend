package posts

import (
	"backend/blog/pkg/db"
	"errors"

	"gorm.io/gorm/clause"
)

type PostsRepository struct {
	Database *db.Db
}

func NewPostsRepository(db *db.Db) *PostsRepository {
	return &PostsRepository{
		Database: db,
	}
}

func (repo *PostsRepository) Create(body *PostCreateRequest, user_id uint) (string, error) {
	post := &Post{
		Title: body.Title,
		Content: body.Content,
		ImageURL: body.ImageURL,
		Published: body.Published,
		UserID: user_id,
	}
	result := repo.Database.Create(post)
	if result.Error != nil {
		return "", result.Error
	}
	return "пост успешно создан", nil
}

func (repo *PostsRepository) GetAllPosts() (*[]Post, error) {
	var posts []Post
	result := repo.Database.Preload("User").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return &posts, nil
}

func (repo *PostsRepository) GetPostById(id uint64) (*Post, error) {
	var post Post

	result:= repo.Database.First(&post, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (repo *PostsRepository) CheckedAuthor(user_id uint, id uint64) (bool, error) {
	var count int64

	result := repo.Database.Model(&Post{}).Where("id = ? AND user_id = ?", id, user_id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (repo *PostsRepository) UpdatePost(post *Post) (*Post, error) {
	var returningPost Post

	result := repo.Database.
		Clauses(clause.Returning{}).
		Model(&Post{}).
		Where("id = ?", post.ID).
		Updates(post)

	if result.Error != nil {
		return nil, result.Error
	}

	err := repo.Database.
		Preload("User").
		First(&returningPost, post.ID).Error

	if err != nil {
		return nil, err
	}

	return &returningPost, nil
}

func (repo *PostsRepository) Delete(id uint64) (string, error) {
	result := repo.Database.Delete(&Post{}, id)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("пост не найден")
	}
	return "пост успешно удален", nil
}	