package token

import (
	"Auth_service/pkg/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID    string `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateAccessToken(id, role, email string) (string, error) {
	claims := Claims{
		ID:    id,
		Role:  role,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(config.Load().ACCESS_TOKEN))
	if err != nil {
		return "", err
	}

	return str, nil
}

func GenerateRefreshToken(id, role, email string) (string, error) {
	claims := Claims{
		ID:    id,
		Role:  role,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Load().REFRESH_TOKEN))
}

func ExtractClaimsRefresh(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load().REFRESH_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func ExtractClaimsAccess(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load().ACCESS_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
