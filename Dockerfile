# syntax=docker/dockerfile:1

FROM golang:1.25.4

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed

# Copy start script
COPY start.sh .
RUN chmod +x start.sh

EXPOSE 8080

# Run
CMD ["./start.sh"]