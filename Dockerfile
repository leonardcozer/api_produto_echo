# Stage 1: Build
FROM golang:1.21-alpine AS builder

# set workdir and copy all source first so go get/tidy can update go.mod in-context
WORKDIR /app

# Copiar o código-fonte
COPY . .

# Baixar dependências (instala git para permitir buscar módulos git-hosted)
RUN apk add --no-cache git && \
	go get go.mongodb.org/mongo-driver && \
	go mod tidy && \
	go mod download

# Compilar a aplicação (build do pacote em cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-go-arquitetura ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o binário do builder
COPY --from=builder /app/api-go-arquitetura .

# Expor porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./api-go-arquitetura"]
