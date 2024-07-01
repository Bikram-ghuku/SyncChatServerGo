MAKEQ := $(MAKE) --no-print-directory
BINARY_NAME=SyncChatServer

ifeq (, $(shell which docker-compose))
    DOCKER_COMPOSE=docker compose
else
    DOCKER_COMPOSE=docker-compose
endif


default: run

.PHONY: dev

dev:
	@echo "Starting database"
	@$(DOCKER_COMPOSE) up postgres-dev -d --wait
	@echo "Starting server"
	@bash -c "trap 'echo "";cd ../ && $(MAKEQ) dev-stop; exit 0' SIGINT SIGTERM ERR; cd backend && go run .;"

.PHONY: dev-stop

dev-stop:
	@echo ""
	@echo "Stopping server and database"
	@$(DOCKER_COMPOSE) stop postgres-dev
	@$(DOCKER_COMPOSE) down postgres-dev

.PHONY: run

run:
	@echo "Starting database"
	@$(DOCKER_COMPOSE) up postgres -d --wait
	@echo "Starting server"
	@$(DOCKER_COMPOSE) up backend-api -d --wait


.PHONY:stop

stop:
	@echo "Stopping backend server and postgres docker containers..."
	@$(DOCKER_COMPOSE) stop postgres backend-api
	@$(DOCKER_COMPOSE) down postgres backend-api

.PHONY: check_clean

check_clean:
	@echo "This will remove the database volume. This action is irreversible."
	@echo -n "Are you sure you want to proceed? [y/N] " && read ans; \
    if [ $${ans:-N} != y ] && [ $${ans:-N} != Y ]; then \
        echo "Operation canceled."; \
        exit 1; \
    fi

.PHONY: clean

clean: check_clean
	sudo rm -rf ./data
