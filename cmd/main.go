package main

import (
	"backend/blog/configs"
	"backend/blog/internal/auth"
	"backend/blog/internal/posts"
	"backend/blog/internal/user"
	"backend/blog/pkg/db"
	"backend/blog/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	db := db.NewDb(config)
	router := http.NewServeMux()

	//Repository
	userRepository := user.NewUserRepository(db)
	postsRepository := posts.NewPostsRepository(db)

	//Service
	authService := auth.NewAuthService(userRepository)
	postsService := posts.NewPostsService(postsRepository)

	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
		Configs:     config,
	})
	posts.NewPostsHandler(router, posts.PostsHandlerDeps{
		PostsService: postsService,
		Configs:     config,
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
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
