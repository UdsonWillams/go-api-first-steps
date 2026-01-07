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

// User representa os dados do usuário extraídos do token
type User struct {
	ID       string   `json:"sub"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Username string   `json:"preferred_username"` // ou "sub" se não tiver
	Roles    []string `json:"-"`
}

const userContextKey = "user_context"

// Authenticator gerencia a verificação de tokens OIDC.
// Ele mantém uma referência ao Verifier do go-oidc para validar tokens JWT.
type Authenticator struct {
	Verifier *oidc.IDTokenVerifier
	ClientID string // ClientID usado para validar resource_access
	DevMode  bool   // Se true, bypassa a autenticação (apenas para desenvolvimento)
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
		ClientID: cfg.ClientID,
		DevMode:  false,
	}, nil
}

// NewDevAuthenticator cria um autenticador para desenvolvimento que bypassa a validação.
// NUNCA use em produção!
func NewDevAuthenticator() *Authenticator {
	return &Authenticator{
		Verifier: nil,
		DevMode:  true,
	}
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
		// DevMode: Bypassa autenticação (apenas para desenvolvimento)
		if a.DevMode {
			slog.DebugContext(c.Request.Context(), "DevMode: Autenticação bypassada")
			c.Set("user_id", "dev-user")
			c.Next()
			return
		}

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
				ResourceAccess map[string]struct {
					Roles []string `json:"roles"`
				} `json:"resource_access"`
				Name              string `json:"name"`
				Email             string `json:"email"`
				PreferredUsername string `json:"preferred_username"`
			}
			if err := idToken.Claims(&claims); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler claims"})
				return
			}

			// Busca as roles específicas do ClientID configurado
			var userRoles []string
			if clientRoles, ok := claims.ResourceAccess[a.ClientID]; ok {
				userRoles = clientRoles.Roles
			} else {
				// Fallback: Se não achar as roles do cliente, assume vazio (sem permissão)
				slog.WarnContext(c.Request.Context(), "Nenhuma role encontrada para o ClientID", "client_id", a.ClientID)
				userRoles = []string{}
			}

			// Injeta usuário Rico no Contexto
			user := &User{
				ID:       idToken.Subject,
				Name:     claims.Name,
				Email:    claims.Email,
				Username: claims.PreferredUsername,
				Roles:    userRoles,
			}
			c.Set(userContextKey, user)

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
				// Pelo menos uma role requerida deve estar presente (modo OR)
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
		} else {
			// Mesmo sem roles requeridas, precisamos parsear o usuário básico
			var claims struct {
				Name              string `json:"name"`
				Email             string `json:"email"`
				PreferredUsername string `json:"preferred_username"`
			}
			_ = idToken.Claims(&claims) // Ignora erro pois verificamos token antes

			c.Set(userContextKey, &User{
				ID:       idToken.Subject,
				Name:     claims.Name,
				Email:    claims.Email,
				Username: claims.PreferredUsername,
			})
		}

		c.Set("user_id", idToken.Subject)
		c.Next()
	}
}

// GetUser recupera o usuário autenticado do contexto
func GetUser(c *gin.Context) *User {
	val, exists := c.Get(userContextKey)
	if !exists {
		return nil
	}
	if user, ok := val.(*User); ok {
		return user
	}
	return nil
}
