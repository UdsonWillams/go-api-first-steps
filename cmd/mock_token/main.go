package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// 1. Gera um par de chaves RSA de 2048 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 2. Prepara a Chave Pública para o .env
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	pemStr := string(pemBlock)
	lines := ""
	for _, line := range splitLines(pemStr) {
		if line != "-----BEGIN PUBLIC KEY-----" && line != "-----END PUBLIC KEY-----" && line != "" {
			lines += line
		}
	}

	fmt.Println("=== 1. COPIE ISTO PARA SEU .ENV (KEYCLOAK_PUBLIC_KEY) ===")
	fmt.Println(lines)
	fmt.Println("=========================================================")
	fmt.Println("")

	// 3. Gera um Token JWT
	claims := jwt.MapClaims{
		"sub": "usuario-teste-id-123",
		"realm_access": map[string]interface{}{
			"roles": []string{"admin", "manager"},
		},
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	// --- CORREÇÃO AQUI ---
	// Removemos o 'privateKey' desta função
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// A chave privada entra APENAS aqui, na hora de assinar
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== 2. USE ESTE TOKEN NO SWAGGER OU POSTMAN ===")
	fmt.Printf("Bearer %s\n", tokenString)
	fmt.Println("===============================================")
}

// splitLines ajuda a limpar a string PEM
func splitLines(s string) []string {
	var lines []string
	var currentLine string
	for _, char := range s {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(char)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}
