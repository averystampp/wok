FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./

EXPOSE 8080

RUN go build ./cmd/main.go
CMD ["./main"]