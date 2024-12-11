FROM golang:1.23.3-alpine3.20 as builder

WORKDIR /app

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy project files
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

EXPOSE 8001

# Copy entrypoint.sh and make it executable
COPY ./docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["go", "run", "cmd/main.go"]
