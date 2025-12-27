package handlers

// CreateProductRequest representa o corpo da requisição POST
type CreateProductRequest struct {
	Name string `json:"name" binding:"required" example:"Monitor UltraWide"`
}

// UpdateProductRequest representa o corpo da requisição PUT
type UpdateProductRequest struct {
	Name string `json:"name" binding:"required" example:"Monitor UltraWide Pro"`
}

// ProductResponse representa a resposta de sucesso com dados
type ProductResponse struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"Monitor UltraWide"`
	CreatedAt string `json:"created_at" example:"2023-12-25T15:00:00Z"`
}

// MessageResponse para mensagens simples
type MessageResponse struct {
	Message string `json:"message" example:"Operação realizada com sucesso"`
}

// ErrorResponse representa erros da API
type ErrorResponse struct {
	Error string `json:"error" example:"Parâmetros inválidos"`
}
