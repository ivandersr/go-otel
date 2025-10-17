FROM golang:1.25 AS build

WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o go-weather .
CMD ["./go-weather"]
