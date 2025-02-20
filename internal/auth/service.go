package auth

import (
	"backend/blog/internal/user"
	"backend/blog/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func(service *AuthService) Register(username, email, password string) (*jwt.JWTData, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return &jwt.JWTData{}, errors.New("пользователь с такими данными уже зарегистрирован")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost) 

	if err != nil {
		return &jwt.JWTData{}, err
	}

	user := &user.User{
			Username: username,
			Email: email,
			Password : string(hashedPassword),
	}
	createdUser, err := service.UserRepository.Create(user)
	if err != nil {
		return &jwt.JWTData{}, err
	}

	jwtData := jwt.JWTData{
		Email: createdUser.Email,
		UserID: createdUser.ID,
	}

	return &jwtData, nil
}

func (service *AuthService) Login(email, password string) (*jwt.JWTData, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)

	var jwtData jwt.JWTData

	if existedUser == nil {
		return &jwt.JWTData{}, errors.New("не правильный email или пароль")
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return &jwt.JWTData{}, errors.New("не правильный email или пароль")
	}

	jwtData = jwt.JWTData{
		Email: existedUser.Email,
		UserID: existedUser.ID,
	}

	return &jwtData, nil
} 