FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server

CMD ["./main"]