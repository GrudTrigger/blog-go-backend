package posts

import (
	"errors"

	"gorm.io/gorm"
)


type PostsService struct {
	PostsRepository *PostsRepository
}

func NewPostsService(postsRepository *PostsRepository) *PostsService {
	return &PostsService{
		PostsRepository: postsRepository,
	}
}

func (service *PostsService) Create(body *PostCreateRequest, user_id uint) (string, error) {
	res, err := service.PostsRepository.Create(body, user_id)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (service *PostsService) GetAllPosts() (*[]Post, error) {
	posts, err := service.PostsRepository.GetAllPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func(service *PostsService) GetPostById(id uint64) (*Post, error) {
	post, err := service.PostsRepository.GetPostById(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (service *PostsService) UploadPost(body *PostUploadRequest, id uint64, user_id uint) (*Post, error) {
	isAuthor, err := service.PostsRepository.CheckedAuthor(user_id, id)
	if err != nil {
		return nil, err
	}
	if(!isAuthor) {
		return nil, errors.New("вы не являетесь автором поста")
	}

	updatePost, err := service.PostsRepository.UpdatePost(&Post{
		Model: gorm.Model{ID: uint(id)},
		Title: body.Title,
		Content: body.Content,
		ImageURL: body.ImageURL,
	})

	if err != nil {
		return nil, err
	}

	return updatePost, nil
}

func (service *PostsService) Delete(id uint64, user_id uint) (string, error) {
	isAuthor, err := service.PostsRepository.CheckedAuthor(user_id, id)
	if err != nil {
		return "", err
	}
	if(!isAuthor) {
		return "", errors.New("вы не являетесь автором поста")
	}

	message, err := service.PostsRepository.Delete(id)
	if err != nil {
		return "", err
	}

	return message, nil
}