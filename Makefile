MAKEQ := $(MAKE) --no-print-directory
BINARY_NAME=SyncChatServer

ifeq (, $(shell which docker-compose))
    DOCKER_COMPOSE=docker compose
else
    DOCKER_COMPOSE=docker-compose
endif


default: dev

.PHONY: dev dev-stop clean check_clean

dev:
	@echo "Starting database"
	@$(DOCKER_COMPOSE) up -d --wait
	@echo "Starting server"
	@bash -c "go run ."
	
dev-stop:
	@echo ""
	@echo "Stopping server and database"
	@$(DOCKER_COMPOSE) down


check_clean:
	@echo "This will remove the database volume. This action is irreversible."
	@echo -n "Are you sure you want to proceed? [y/N] " && read ans; \
    if [ $${ans:-N} != y ] && [ $${ans:-N} != Y ]; then \
        echo "Operation canceled."; \
        exit 1; \
    fi

clean: check_clean
	sudo rm -rf ./data
