DOCKER_COMPOSE ?= $(shell which docker-compose)
DOCKER_COMPOSE_RUN := $(DOCKER_COMPOSE) run --rm

.SILENT: help
help: ## Show this help message
	echo "usage: make [target] ..."
	echo ""
	echo "targets:"
	fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match $$'\t' | sed -e 's/: ## / - /'

start: ## Start the service
	$(DOCKER_COMPOSE) up

.PHONY: clean
clean: ## Clean up any containers and images
	rm -r vendor/
	$(DOCKER_COMPOSE) stop
	$(DOCKER_COMPOSE) rm -f
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans
