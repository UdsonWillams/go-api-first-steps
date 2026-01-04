# Concorrência em Go: O Poder dos Gophers

A concorrência é o principal diferencial do Go. Enquanto outras linguagens (como Python ou Java) adicionaram concorrência depois, o Go nasceu com ela.

## 1. Goroutines vs Threads vs AsyncIO

### O Problema das Threads (Java, C++, Python Threads)
Threads de Sistema Operacional são pesadas. Cada uma consome ~1MB de RAM. Se você criar 10.000 threads, seu servidor explode.

### O Modelo Async/Await (Python FastAPI, Node.js)
Usa uma única thread (Event Loop). É leve, mas cooperativo. Se uma função "esquecer" de ser async ou fizer um cálculo pesado (CPU bound), trava tudo.

### O Jeito Go (Goroutines)
Goroutines são "Green Threads" gerenciadas pelo Go Runtime, não pelo OS.
- **Leveza**: Iniciam com apenas **2KB** de stack.
- **Escalabilidade**: Você pode rodar **milhões** de goroutines em uma máquina normal.
- **Multiplexação**: O Go Runtime distribui essas goroutines automaticamente entre as threads do processador.

## 2. Como usar

### Iniciando uma Goroutine
Basta colocar a palavra `go` na frente de qualquer função.

```go
func enviarEmail() {
    // ... demora 2 segundos ...
}

func main() {
    // O programa NÃO espera terminar. Ele segue em frente imediatamente.
    go enviarEmail()
}
```

## 3. Channels (Canais)
"Don't communicate by sharing memory; share memory by communicating." (Rob Pike)

Channels são tubos tipados que conectam goroutines. É como elas conversam sem criar Race Conditions.

```go
func main() {
    mensagens := make(chan string)

    // Goroutine produtora
    go func() {
        mensagens <- "Olá do outro lado!" // Envia
    }()

    // Goroutine consumidora (a main espera chegar algo)
    msg := <-mensagens // Recebe
    println(msg)
}
```

## 4. WaitGroup
Quando você quer esperar várias tarefas terminarem antes de sair.

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        println("Trabalhando...", id)
    }(i)
}

wg.Wait() // Bloqueia até as 5 terminarem
println("Tudo pronto!")
```
