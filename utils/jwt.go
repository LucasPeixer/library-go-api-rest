package utils

import (
	"fmt"
	"go-api/initializers"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT cria um novo token JWT utilizando o Id do usuário
// e sua Role e retorna o token gerado ou um erro.
func GenerateJWT(userId int, role string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Issuer:    fmt.Sprint(userId),
		Subject:   role, // Armazena a role na claim subject
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(initializers.JwtKey)
}

// ValidateJWT valida o token JWT e retorna as claims.
func ValidateJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return initializers.JwtKey, nil
	})
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
