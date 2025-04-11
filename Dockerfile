FROM golang:1.24.2-alpine3.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Install Swaggo
RUN go install github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc4

COPY . .

# Run swag init in your project directory to generate Swagger documentation
RUN swag init -v3.1 -o docs -g main.go --parseDependency --parseInternal

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
