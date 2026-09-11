FROM golang:1.23-alpine

WORKDIR /app

COPY . .
RUN go build -o adex ./cmd/main.go


EXPOSE 8080


CMD ["./adex"]