FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /proxy-server ./cmd/cli

EXPOSE 8083

ENTRYPOINT ["/proxy-server", "proxy"]
CMD ["--host", "0.0.0.0", "--port", "8083"]
