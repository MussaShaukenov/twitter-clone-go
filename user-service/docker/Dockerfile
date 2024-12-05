FROM golang:1.23.3-alpine3.20 as builder

WORKDIR /app

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy project files
COPY ../go.mod go.sum ./
RUN go mod download

COPY .. .

EXPOSE 8002

ARG DATABASE_URL=postgres://postgres:admin@db:5432/twitter_clone?sslmode=disable

# Run migrations
#ENTRYPOINT goose -dir ./internal/migrations postgres $DATABASE_URL up

# Build the application
CMD ["go", "run", "cmd/main.go"]
