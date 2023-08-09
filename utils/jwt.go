package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成 Token
func GenerateToken(id uint, username string) (string, error) {
	// 实例化一个 claims 对象
	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("jwt.expire_duration")) * time.Second).Unix(),
			Issuer:    viper.GetString("jwt.issuer"),
			Subject:   viper.GetString("jwt.subject"),
		},
	}

	// 生成 token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(viper.GetString("jwt.key")))
}

// ParseToken 解析 token
func ParseToken(token string) (*Claims, error) {
	// 解析 token
	tokenClaims, err := jwt.ParseWithClaims(token, new(Claims), func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.key")), nil
	})
	if err != nil {
		return nil, err
	}

	// 断言获取 claims
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
