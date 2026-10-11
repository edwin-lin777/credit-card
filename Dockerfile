
# Stage 1: Build our Go application
FROM golang:1.27-alpine AS builder

WORKDIR /src

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile Go into a Linux executable
RUN CGO_ENABLED=0 GOOS=linux go build -o /study ./cmd/study


# Stage 2: Create the final lightweight image
FROM scratch

WORKDIR /app

# Copy the compiled application
COPY --from=builder /study /app/study

# Copy credit card YAML files
COPY --from=builder /src/cards /app/cards

# Run the application
ENTRYPOINT ["/app/study"]
