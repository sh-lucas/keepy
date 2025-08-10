# ---- Estágio de Build ----
# Usa a imagem oficial do Go para compilar a aplicação
FROM golang:1.23.3 AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia o código-fonte para dentro do container
COPY . .

# Compila o binário de forma estática (sem CGO) para ser compatível com a imagem scratch
# CGO_ENABLED=0 cria um binário que não depende de bibliotecas C externas
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./main ./server/.

# ---- Estágio Final ----
# Usa a imagem 'scratch', que é uma imagem vazia, como base
FROM debian:stable-slim

# Copia apenas o binário compilado do estágio de build
COPY --from=builder ./app/main ./main

# Expõe a porta que a aplicação usará (opcional, mas boa prática)
EXPOSE 80

# Comando para executar a aplicação quando o container iniciar
CMD ["./main"]