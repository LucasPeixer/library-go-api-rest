package middleware

import (
	"go-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware autentica uma requisição usando JWT.
func JWTAuthMiddleware(c *gin.Context) {
	// Verifica se o header de autorização existe
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	// Define as informações do usuário no contexto
	c.Set("userId", claims.Issuer)
	c.Set("role", claims.Subject)
	c.Next()
}

// RoleRequired verifica se o usuário possui a role necessária para acessar o recurso.
func RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtém a role do usuário a partir do contexto
		userRole, existe := c.Get("role")
		if !existe {
			// Se a role não for encontrada, retorna status Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Compara a role do usuário com a role necessária
		if userRole != role {
			// Se as roles não coincidirem, retorna status Forbidden
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Se a role for válida, permite que a requisição prossiga
		c.Next()
	}
}
