package comment

import "backend/blog/pkg/jwt"

type CommentService struct {
	*CommentRepository
}

func NewCommentService(commentRepository *CommentRepository) *CommentService {
	return &CommentService{
		commentRepository,
	}
}

func (service *CommentService) CreatePostComment(idPost uint64, userData *jwt.JWTData, text string) (string, error) {
	response, err := service.CommentRepository.CreatePostComment(idPost, userData.UserID, text)
	if err != nil {
		return "", err
	}
	return response, nil
}

func(service *CommentService) UpdateComment(idComment uint64, text string, userData *jwt.JWTData) (string, error) {

	author, err := service.CommentRepository.CheckedAuthorComment(idComment, userData.UserID)
	if err != nil {
		return "", err
	}

	if !author {
		return "вы не являетесь автором комментария", nil
	}

	result, err := service.CommentRepository.UpdateComment(idComment, text)
	if err != nil {
		return "", err
	}

	return result, nil
}

func(service *CommentService) Delete(id uint64, user_id uint) (string, error) {
	author, err := service.CommentRepository.CheckedAuthorComment(id, user_id)
	if err != nil {
		return "", err
	}
	if !author {
		return "вы не являетесь автором комментария", nil
	}
	result, err := service.CommentRepository.Delete(id)
	if err != nil {
		return "", err
	}
	return result, nil
}