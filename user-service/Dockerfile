FROM golang:1.23.3-alpine3.20

# Set working directory
WORKDIR /app

# Install goose (migration tool)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify


# Copy application files
COPY . .

# Expose the application port
EXPOSE 8002

# Command to start the app
CMD ["go", "run", "cmd/main.go"]
