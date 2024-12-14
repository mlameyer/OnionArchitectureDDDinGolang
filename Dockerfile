FROM golang:1.23-alpine
WORKDIR /app
COPY . .

ENV AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
ENV AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
ENV AWS_REGION=$AWS_REGION
ENV AWS_SECRET_NAME=$AWS_SECRET_NAME

RUN go install -mod=mod github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN apk add --no-cache bash
CMD ["sh", "-c", "source .env.development && go run cmd/main.go"]
