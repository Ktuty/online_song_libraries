# Makefile

# Загрузка переменных окружения из .env
include .env

# Переменные
MIGRATE_CMD := migrate
MIGRATE_DIR := ./schema
MIGRATE_EXT := sql
DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Цель по умолчанию
.DEFAULT_GOAL := help

# Цель для создания новой миграции
create-migration:
	@echo "Creating new migration with name: $(MIGRATION_NAME)"
	$(MIGRATE_CMD) create -ext $(MIGRATE_EXT) -dir $(MIGRATE_DIR) -seq $(MIGRATION_NAME)

# Цель для применения миграций
migrate-up:
	$(MIGRATE_CMD) -path $(MIGRATE_DIR) -database "$(DB_URL)" up

# Цель для отката миграций
migrate-down:
	$(MIGRATE_CMD) -path $(MIGRATE_DIR) -database "$(DB_URL)" down $(STEPS)

# Цель для отображения справки
help:
	@echo "Usage:"
	@echo "  make create-migration MIGRATION_NAME=<name>"
	@echo "  make migrate-up"
	@echo "  make migrate-down STEPS=<number>"
