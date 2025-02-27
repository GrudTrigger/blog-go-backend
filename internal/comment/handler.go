package comment

import (
	"backend/blog/configs"
	"backend/blog/pkg/jwt"
	"backend/blog/pkg/middleware"
	"backend/blog/pkg/request"
	"backend/blog/pkg/response"
	"net/http"
	"strconv"
)

type CommentHandlerDeps struct {
	*CommentService
	*configs.Configs
}

type CommentHandler struct {
	CommentService *CommentService
	config *configs.Configs
}

func NewCommentHandler(router *http.ServeMux, deps CommentHandlerDeps) {
	handler := &CommentHandler{
		CommentService: deps.CommentService,
		config: deps.Configs,
	}
	router.Handle("POST /comment-post/{id}", middleware.IsAuthed(handler.Create(), handler.config))
	router.Handle("PATCH /comment-post/{comment_id}", middleware.IsAuthed(handler.Update(), handler.config))
	router.Handle("DELETE /comment-post/{comment_id}", middleware.IsAuthed(handler.Delete(), handler.config))
}

func(handler *CommentHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		idPost, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userData, ok := r.Context().Value(middleware.UserContextKey).(*jwt.JWTData)
		if !ok || userData == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := request.HandleBody[CommentPostRequest](&w, r)
		if err != nil {
			return
		}

		data, err := handler.CommentService.CreatePostComment(idPost, userData, body.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, data, 201)
	}
}

func(handler *CommentHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("comment_id")
		idComment, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userData, ok := r.Context().Value(middleware.UserContextKey).(*jwt.JWTData)
		if !ok || userData == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := request.HandleBody[CommentPostUpdateRequest](&w, r)
		if err != nil {
			return
		}

		data, err := handler.CommentService.UpdateComment(idComment, body.Text, userData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, data, 200)
	}
}

func(handler *CommentHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("comment_id")
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

		data, err := handler.CommentService.Delete(id, userData.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, data, 200)
	}
}