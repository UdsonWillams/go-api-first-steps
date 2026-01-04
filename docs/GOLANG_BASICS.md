# Go para Desenvolvedores Python

Um guia rápido para quem vem do Python e quer entender o que está acontecendo aqui.

## 1. Tipagem Estática vs Dinâmica
- **Python**: `x = 1`, depois `x = "texto"`. Funciona.
- **Go**: `var x int = 1`. Se tentar `x = "texto"`, nem compila.
  - Vantagem: O compilador pega bugs que em Python só estourariam em produção.

## 2. Tratamento de Erros (Try/Except vs if err != nil)
- **Python**:
  ```python
  try:
      file = open("doc.txt")
  except Exception as e:
      print(e)
  ```
- **Go**: Erros são valores normais. Não existem Exceptions.
  ```go
  file, err := os.Open("doc.txt")
  if err != nil {
      log.Println("Erro:", err)
      return
  }
  ```
  Isso força você a lidar com o erro imediatamente, tornando o código mais seguro.

## 3. Structs vs Classes
Go não tem Classes nem Herança. Go tem **Structs** e **Composição**.

```go
type User struct {
    Name string
}

// Método anexado ao struct (parece o self do Python)
func (u *User) SayHello() {
    println("Hello " + u.Name)
}
```

## 4. Interfaces (Duck Typing Tipado)
Em Go, **você não declara que implementa uma interface**. Se você tem os métodos, você implementa.
- Se anda como pato e faz quack como pato, o Go trata como pato.

## 5. Concorrência (Killer Feature)
- **Python**: Threads reais são pesadas (GIL). Asyncio é complexo.
- **Go**: Goroutines são threads super leves. Você pode iniciar 100 mil delas sem travar a máquina.
  ```go
  go func() {
      print("Rodando em paralelo!")
  }()
  ```
