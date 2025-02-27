package posts

import (
	"backend/blog/pkg/db"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm/clause"
)

type PostsRepository struct {
	Database *db.Db
	cache *redis.Client
}

func NewPostsRepository(db *db.Db, cache *redis.Client) *PostsRepository {
	return &PostsRepository{
		Database: db,
		cache: cache,
	}
}

func (repo *PostsRepository) Create(body *PostCreateRequest, user_id uint, ctxRedis context.Context) (string, error) {
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
	repo.cache.Del(ctxRedis, "all_posts")
	return "пост успешно создан", nil
}

func (repo *PostsRepository) GetAllPosts(ctx context.Context) (*[]Post, error) {
	cachedPosts,err := repo.cache.Get(ctx, "all_posts").Result()
	if err == redis.Nil {
		//Если в кеше пусто, берем посты из бд
		var posts []Post
		result := repo.Database.Preload("User").Preload("Comments").Find(&posts)
		if result.Error != nil {
			return nil, result.Error
		}
		// Сохраняем в Redis (JSON)
		jsonData, err := json.Marshal(posts)
		if err != nil {
			return nil, err
		}
		repo.cache.Set(ctx, "all_posts", jsonData, time.Minute *  5)
		return &posts, nil
	}
		// Если данные есть в кеше, десериализуем
		var posts []Post
		json.Unmarshal([]byte(cachedPosts), &posts)
		return &posts, nil
}

func (repo *PostsRepository) GetPostById(id uint64) (*Post, error) {
	var post Post

	result:= repo.Database.Preload("User").Preload("Comments").First(&post, "id = ?", id)
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

func (repo *PostsRepository) UpdatePost(post *Post, ctxRedis context.Context) (*Post, error) {
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
	repo.cache.Del(ctxRedis, "all_posts")
	return &returningPost, nil
}

func (repo *PostsRepository) Delete(id uint64, ctxRedis context.Context) (string, error) {
	result := repo.Database.Delete(&Post{}, id)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("пост не найден")
	}
	repo.cache.Del(ctxRedis, "all_posts")
	return "пост успешно удален", nil
}	