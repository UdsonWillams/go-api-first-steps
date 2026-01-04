package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"go-api-first-steps/internal/config"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

// Authenticator gerencia a verificação de tokens OIDC.
// Ele mantém uma referência ao Verifier do go-oidc para validar tokens JWT.
type Authenticator struct {
	Verifier *oidc.IDTokenVerifier
}

// NewAuthenticator inicializa o Provider OIDC e configura o Verifier.
// Esta função deve ser chamada apenas uma vez na inicialização da aplicação (Singleton).
func NewAuthenticator(cfg *config.Config) (*Authenticator, error) {
	// 1. Configura o Provider (Auto Discovery das chaves)
	provider, err := oidc.NewProvider(context.Background(), cfg.KeycloakURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao inicializar OIDC Provider: %w", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	return &Authenticator{
		Verifier: verifier,
	}, nil
}

// CheckMiddleware cria um handler do Gin para validar o token JWT presente no header Authorization.
//
// Parâmetros:
//   - mode: Define a lógica de validação de roles. Pode ser "AND" (todas as roles necessárias) ou "OR" (pelo menos uma). Default: "OR".
//   - requiredRoles: Lista de roles que o usuário deve possuir para acessar o recurso.
//
// Exemplo:
//
//	auth.CheckMiddleware("OR", "admin", "manager") // Requer admin OU manager
//	auth.CheckMiddleware("AND", "admin", "finance") // Requer admin E finance
func (a *Authenticator) CheckMiddleware(mode string, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.Verifier == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Autenticação não configurada"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token não informado"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Valida o Token
		idToken, err := a.Verifier.Verify(c.Request.Context(), tokenString)
		if err != nil {
			slog.WarnContext(c.Request.Context(), "Token inválido", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// Validação de Roles
		if len(requiredRoles) > 0 {
			var claims struct {
				RealmAccess struct {
					Roles []string `json:"roles"`
				} `json:"realm_access"`
			}
			if err := idToken.Claims(&claims); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler claims"})
				return
			}

			userRoles := claims.RealmAccess.Roles
			rolesMap := make(map[string]bool)
			for _, r := range userRoles {
				rolesMap[r] = true
			}

			if strings.ToUpper(mode) == "AND" {
				// Todas as roles requeridas devem estar presentes
				for _, req := range requiredRoles {
					if !rolesMap[req] {
						c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Sem permissão (Faltam roles)"})
						return
					}
				}
			} else {
				// Modo OR (Default): Pelo menos uma role deve estar presente
				hasAtLeastOne := false
				for _, req := range requiredRoles {
					if rolesMap[req] {
						hasAtLeastOne = true
						break
					}
				}
				if !hasAtLeastOne {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Sem permissão"})
					return
				}
			}
		}

		c.Set("user_id", idToken.Subject)
		c.Next()
	}
}
