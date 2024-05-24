FROM golang:latest AS build
LABEL authors="Sonia Fatholahi"

WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Install the dependencies
RUN go mod download

COPY . .

# Disabling CGO since libc.so.6 doesn't included in alpine image
ENV CGO_ENABLED=0

# Build the application
RUN go build -o /app/payment cmd/main.go

FROM alpine:3.19 AS production

WORKDIR /app

# Copy the built files from the previous stage
COPY --from=build /app/payment payment

# Expose the program port
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/app/payment"]
