package comment

import (
	"backend/blog/pkg/db"
	"errors"
)

type CommentRepository struct {
	Database *db.Db
}

func NewCommentRepository(db *db.Db) *CommentRepository {
	return &CommentRepository{
		Database: db,
	}
}

func (repo *CommentRepository) CreatePostComment(idPost uint64, idUser uint, text string) (string, error) {
	comment := &CommentPost{
		Text: text,
		UserID: idUser,
		PostID: uint(idPost),
	}
	result := repo.Database.Create(comment)
	if result.Error != nil {
		return "", result.Error
	}
	return "комментарий успешно создан", nil
}

func (repo *CommentRepository) UpdateComment(id uint64, text string) (string, error) {

	result := repo.Database.Model(&CommentPost{}).Where("id = ?", id).Update("text", text)
	if result.Error != nil {
		return "", result.Error
	}

	return "комментарий успешно обновлен", nil
}

func(repo *CommentRepository) Delete(id uint64) (string , error) {
	result := repo.Database.Delete(&CommentPost{}, id)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("комментарий не найден")
	}
	return "комментарий удален", nil
}

func (repo *CommentRepository) CheckedAuthorComment(id uint64, user_id uint) (bool, error) {
	var count int64
	result := repo.Database.Model(&CommentPost{}).Where("id = ? AND user_id = ?", id, user_id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
