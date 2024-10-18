ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app


# Copia os arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copia todo o código fonte
COPY . .

# Compila o aplicativo - corrigindo o caminho e o output
RUN go build -v -o /usr/local/bin/app ./cmd/api

# Segunda etapa - imagem final
FROM debian:bookworm

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# Copia apenas o binário compilado
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

# Executa o aplicativo
CMD ["/usr/local/bin/app"]