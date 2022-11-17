FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build ./...
RUN go build -o /glasnik

EXPOSE 3000

CMD ["/glasnik"]