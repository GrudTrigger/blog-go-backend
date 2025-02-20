package auth

import (
	"backend/blog/configs"
	"backend/blog/pkg/jwt"
	"backend/blog/pkg/request"
	"backend/blog/pkg/response"
	"net/http"
)

type AuthHandler struct {
	*AuthService
	config *configs.Configs
}

type AuthHandlerDeps struct {
	*AuthService
	*configs.Configs
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		config:      deps.Configs,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		jwtData, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}		
		token, err := jwt.NewJWT(handler.config.Secret).Create(jwt.JWTData{
			Email: jwtData.Email,
			UserID: jwtData.UserID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		response.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		jwtData, err := handler.AuthService.Register(body.Username, body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		token, err := jwt.NewJWT(handler.config.Secret).Create(jwt.JWTData{Email: jwtData.Email, UserID: jwtData.UserID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}
		response.Json(w, data, 200)
	}
}
