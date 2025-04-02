package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManager struct {
	secretKey    string
	tokenExpires time.Duration
}

func NewJWTManager(secretKey string, tokenExpires time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:    secretKey,
		tokenExpires: tokenExpires,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTManager) GenerateToken(userID string, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpires)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
