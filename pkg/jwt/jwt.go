package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"lihood/conf"
	"time"
)

type StandardClaims struct {
	ID        int64 `json:"id"`
	ExpiresAt int64 `json:"create_at"`
}

func (s StandardClaims) Valid() error {
	if s.ExpiresAt < time.Now().Unix() {
		return errors.New("token is expired")
	}
	return nil
}

func GenToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StandardClaims{
		ID:        id,
		ExpiresAt: time.Now().Unix() + conf.Instance.Jwt.Expire,
	})
	return token.SignedString([]byte(conf.Instance.Jwt.SecretKey))
}

func ParseToken(token string) (int64, error) {
	claims := StandardClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Instance.Jwt.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}
