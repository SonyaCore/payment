# Makefile

GO := go

# Paths
MIGRATE_PATH := migrate/migrate.go
CMD_PATH := cmd/main.go

# Targets
.PHONY: migrate run docker

# Migration
migrate:
	$(GO) run $(MIGRATE_PATH)

# Run
run:
	$(GO) run $(CMD_PATH)

# Build
	$(GO) build $(CMD_PATH) -ldflags="-s -w"

# Docker run
docker:
	docker-compose up -d
