package posts

import (
	"backend/blog/configs"
	"backend/blog/pkg/jwt"
	"backend/blog/pkg/middleware"
	"backend/blog/pkg/request"
	"backend/blog/pkg/response"
	"net/http"
	"strconv"
)

type PostsHandlerDeps struct {
	*PostsService
	*configs.Configs
}

type PostsHandler struct {
	PostsService *PostsService
	config *configs.Configs
}

func NewPostsHandler(router *http.ServeMux, deps PostsHandlerDeps) {
	handler := &PostsHandler{
		PostsService: deps.PostsService,
		config: deps.Configs,
	}
	router.Handle("POST /post", middleware.IsAuthed(handler.Create(), handler.config))
	router.HandleFunc("GET /posts-all", handler.GetAllPosts())
	router.HandleFunc("GET /posts/{id}", handler.GetPostById())
	router.Handle("PATCH /posts/{id}", middleware.IsAuthed(handler.UploadPost(), deps.Configs))
	router.Handle("DELETE /posts/{id}", middleware.IsAuthed(handler.Delete(), deps.Configs))
}

func(handler *PostsHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[PostCreateRequest](&w, r)
		if err != nil {
			return
		}

		userData, ok := r.Context().Value(middleware.UserContextKey).(*jwt.JWTData)
		if !ok || userData == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		result,err := handler.PostsService.Create(body, userData.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, result, 201)
	}
}

func (handler *PostsHandler) GetAllPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := handler.PostsService.GetAllPosts()
		if err != nil {
			http.Error(w, "ошибка при получении постов", http.StatusBadRequest)
			return
		}
		response.Json(w, posts, 200)
	}
}

func (handler *PostsHandler) GetPostById() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := handler.PostsService.GetPostById(id)
		if err != nil {
			http.Error(w, "ошибка при получении поста", http.StatusBadRequest)
			return
		}
		response.Json(w, data, 200)
	}
}

func (handler *PostsHandler) UploadPost() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[PostUploadRequest](&w, r)
		if err != nil {
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userData, ok := r.Context().Value(middleware.UserContextKey).(*jwt.JWTData)
		if !ok || userData == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		uploadPost, err := handler.PostsService.UploadPost(body, id, userData.UserID)
		if err != nil {
			http.Error(w, "ошибка при редактировании поста", http.StatusBadRequest)
			return
		}
		response.Json(w, uploadPost, 201)
	}
}

func (handler *PostsHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userData, ok := r.Context().Value(middleware.UserContextKey).(*jwt.JWTData)
		if !ok || userData == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		message, err := handler.PostsService.Delete(id, userData.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, message, 200)
	}
}