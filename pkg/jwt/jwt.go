package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct{
	Email string
	UserID uint
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {

	 t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"user_id": data.UserID,
	 })

	 s, err := t.SignedString([]byte(j.Secret))

	 if err != nil {
		return "", err
	 }

	 return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}

	email := t.Claims.(jwt.MapClaims)["email"]
	id := t.Claims.(jwt.MapClaims)["user_id"]
	return t.Valid, &JWTData{Email: email.(string), UserID: uint(id.(float64))}
}