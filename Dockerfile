FROM golang:1.21-alpine

LABEL name = "banners-api"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /banners-api cmd/app/main.go

EXPOSE 8080

CMD ["/banners-api"]