FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /registry-server ./cmd/cli

EXPOSE 8082

ENTRYPOINT ["/registry-server", "registry"]
CMD ["--port", "8082"]
