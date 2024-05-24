FROM golang:1.22.3-alpine AS build
LABEL authors="Sonia Fatholahi"

WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Install the dependencies
RUN go mod download

COPY . .

# Build
RUN go build -o /app/migrate migrate/migrate.go

FROM alpine:3.19 AS production

WORKDIR /app

COPY --from=build /app/migrate/migrate migrate

# Command to run the migrate binary
ENTRYPOINT ["/app/migrate"]