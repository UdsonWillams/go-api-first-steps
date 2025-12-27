package middleware

import (
	"crypto/rsa"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// KeycloakClaims mapeia a estrutura interna do Token
type KeycloakClaims struct {
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	jwt.RegisteredClaims
}

// Auth agora recebe a Chave Pública carregada no main
func Auth(publicKey *rsa.PublicKey, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token não informado"})
			return
		}

		// Remove "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 1. Parse COM Validação de Assinatura (Seguro)
		token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(t *jwt.Token) (interface{}, error) {
			// Valida se o algoritmo é realmente RSA (evita ataques de alg: none)
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return publicKey, nil
		})

		// 2. Tratamento de erros de validação
		if err != nil || !token.Valid {
			slog.WarnContext(c.Request.Context(), "Token inválido ou expirado", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// 3. Extrai as Claims
		claims, ok := token.Claims.(*KeycloakClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Erro nas claims"})
			return
		}

		// 4. Validação de Role (RBAC)
		if requiredRole != "" {
			hasRole := false
			for _, r := range claims.RealmAccess.Roles {
				if r == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				slog.WarnContext(c.Request.Context(), "Acesso negado: Role insuficiente",
					"user_id", claims.Subject,
					"required_role", requiredRole,
					"user_roles", claims.RealmAccess.Roles)

				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Sem permissão"})
				return
			}
		}

		// Injeta ID no contexto
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
