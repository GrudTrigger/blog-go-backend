package main

import (
	"backend/blog/configs"
	"backend/blog/internal/auth"
	"backend/blog/internal/comment"
	"backend/blog/internal/posts"
	"backend/blog/internal/user"
	"backend/blog/pkg/db"
	"backend/blog/pkg/middleware"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var ctx  = context.Background()

func main() {
	config := configs.LoadConfig()
	db := db.NewDb(config)
	router := http.NewServeMux()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	// Проверка соединения с Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
			log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	fmt.Println("Подключение к Redis успешно!")
	//Repository
	userRepository := user.NewUserRepository(db)
	postsRepository := posts.NewPostsRepository(db, rdb)
	commentRepository := comment.NewCommentRepository(db)

	//Service
	authService := auth.NewAuthService(userRepository)
	postsService := posts.NewPostsService(postsRepository)
	commentService := comment.NewCommentService(commentRepository)

	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
		Configs:     config,
	})
	posts.NewPostsHandler(router, posts.PostsHandlerDeps{
		PostsService: postsService,
		Configs:     config,
	})
	comment.NewCommentHandler(router, comment.CommentHandlerDeps{
		CommentService: commentService,
		Configs: config,
	})

	//Middlewares
	stack := middleware.Chain(
		middleware.Logging,
		middleware.CORS,
	)

	server := http.Server{
		Addr:    ":8082",
		Handler: stack(router),
	}

	fmt.Println("Сервер запущен на порту 8082")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
