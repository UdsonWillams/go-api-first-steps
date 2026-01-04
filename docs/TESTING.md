# Testes e Benchmarks em Go

Go possui um framework de testes integrado robusto. Não precisa instalar pytest ou unitttest, já vem na linguagem.

## 1. Testes Unitários (`_test.go`)
Arquivos de teste devem terminar com `_test.go` e as funções devem começar com `TestXxx`.

### Table-Driven Tests (Padrão Go)
Em Go, evitamos escrever 10 funções de teste repetidas. Criamos uma tabela de cenários:

```go
func TestSoma(t *testing.T) {
    tests := []struct {
        nome     string
        a, b     int
        esperado int
    }{
        {"Positivos", 2, 2, 4},
        {"Zeros", 0, 0, 0},
        {"Negativos", -1, -1, -2},
    }

    for _, tt := range tests {
        t.Run(tt.nome, func(t *testing.T) {
            resultado := tt.a + tt.b
            if resultado != tt.esperado {
                t.Errorf("Esperado %d, veio %d", tt.esperado, resultado)
            }
        })
    }
}
```

## 2. Rodando Testes
- Rodar todos: `go test ./...`
- Rodar com detalhes: `go test -v ./...`
- Cobertura: `go test -cover ./...`

## 3. Benchmarks (Performance)
Go tem ferramenta nativa para medir performance. Funções começam com `BenchmarkXxx`.

```go
func BenchmarkSoma(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Soma(2, 2)
    }
}
```

### Rodando Benchmark
Execute no terminal:
```bash
go test -bench=. ./...
```

Ele vai rodar a função milhões de vezes e te dizer exatamente quantos **nanosegundos** cada operação leva. É excelente para otimizar código crítico.
