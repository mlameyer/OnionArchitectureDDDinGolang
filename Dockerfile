FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN go install -mod=mod github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN apk add --no-cache bash
CMD ["sh", "-c", "source .env.development && go run cmd/main.go"]
