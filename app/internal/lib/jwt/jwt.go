package jwt

import (
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user model.User, app model.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["app_id"] = app.ID
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["login"] = user.Login

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
