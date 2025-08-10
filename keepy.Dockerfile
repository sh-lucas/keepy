# keepy.Dockerfile

# Estágio 1: Compilação da aplicação Go
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Otimização de cache de dependências
COPY go.mod ./
RUN go mod download

# Copia e compila
COPY . .
RUN CGO_ENABLED=0 go install .

# Estágio 2: Imagem final, mínima
FROM debian:12-slim

# Copia o binário para um diretório padrão no PATH
COPY --from=builder /go/bin/keepy /usr/local/bin/keepy

# Agora o ENTRYPOINT pode chamar o comando diretamente, sem o caminho
ENTRYPOINT ["keepy"]