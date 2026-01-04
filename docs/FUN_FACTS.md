# Fun Facts sobre Go 🐹

Curiosidades para você contar no cafézinho.

## 1. A Origem
O Go foi criado no **Google** em 2007, por três lendas da computação que estavam cansadas da complexidade do C++:
- **Ken Thompson**: Criador do UNIX (sim, o sistema operacional inteiro), da linguagem B (pai do C) e da codificação UTF-8.
- **Rob Pike**: Criador do Plan 9 e UTF-8.
- **Robert Griesemer**: Criador da JVM HotSpot.

Dizem que eles desenharam o Go enquanto esperavam um _build_ gigante de C++ compilar (demorava 45 minutos!). O compilador do Go é um dos mais rápidos do mundo por causa desse trauma.

## 2. O Mascote: Gopher
O mascote não tem nome, é apenas "The Go Gopher". Ele foi desenhado por **Renee French** (esposa do Rob Pike).
Ele é baseado em desenhos que ela fez para uma camiseta da rádio WFMU anos antes.

O Gopher é icônico e a comunidade cria milhares de versões dele (Gopher Batman, Gopher Jedi, etc).

## 3. Sem Generics? (Por 10 anos!)
O Go ficou famoso (e infame) por não ter Generics (tipo `List<T>`) por mais de uma década. Os criadores diziam que "adicionava complexidade demais". Finalmente, na versão 1.18 (2022), eles cederam e implementaram.

## 4. Formatador Oficial
O Go foi a primeira linguagem popular a impor um estilo de código único via ferramenta: `gofmt`.
Não existe discussão sobre "onde colocar a chave" em Go. O `gofmt` decide e pronto. Isso acabou com 90% das discussões inúteis em Code Reviews.

## 5. Keywords
O Go só tem **25 palavras reservadas** (keywords). Para comparação:
- C++: ~90
- Java: ~50
- Python: ~35

É uma linguagem desenhada para ser lida facilmente.
