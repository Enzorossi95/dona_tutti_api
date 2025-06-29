# Choose whatever you want, version >= 1.16
FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

COPY go.mod go.sum ./
RUN go mod download

RUN swag init

EXPOSE 9999

CMD ["air", "-c", ".air.toml"]