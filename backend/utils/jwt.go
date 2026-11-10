package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"blog-api/constant"
)

// Claims 自定义 JWT 载荷
type Claims struct {
	UserID   uint          `json:"user_id"`
	Username string        `json:"username"`
	Role     constant.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 token，同时返回过期时间（供登录接口返回 expires_at）
func GenerateToken(secret []byte, userID uint, username string, role constant.Role, expire time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(expire)
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// ParseToken 解析并校验 token
func ParseToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
