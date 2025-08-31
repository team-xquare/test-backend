package jwt

import (
	"time"

	"github.com/team-xquare/deployment-platform/internal/pkg/config"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uint, email string) (accessToken, refreshToken string, err error) {
	accessExpiry, _ := time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	refreshExpiry, _ := time.ParseDuration(config.AppConfig.JWTRefreshExpiry)

	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", "", errors.Internal("Failed to generate access token")
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", "", errors.Internal("Failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, errors.Unauthorized("Invalid token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.Unauthorized("Invalid token")
}
