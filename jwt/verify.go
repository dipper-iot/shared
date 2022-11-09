package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenGenerateInput struct {
	Id   string
	Role string
}

func (manager *JWTManager) Verify(accessToken string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (manager *JWTManager) Generate(user *TokenGenerateInput, claim *jwt.StandardClaims) (string, error) {
	if user == nil || claim == nil {
		return "", errors.New("Error")
	}
	claim.ExpiresAt = time.Now().Add(manager.tokenDuration).Unix()
	claims := JwtClaims{
		StandardClaims: *claim,
		ID:             user.Id,
		Role:           user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}
