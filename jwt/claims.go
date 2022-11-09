package jwt

import "github.com/dgrijalva/jwt-go"

type JwtClaims struct {
	jwt.StandardClaims
	ID   string `json:"id"`
	Role string `json:"role"`
}
