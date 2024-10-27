FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go install -mod=mod github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY .env .env
RUN apk add --no-cache bash
CMD ["sh", "-c", "source .env && go run cmd/main.go"]
