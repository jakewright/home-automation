DOCKER_COMPOSE ?= $(shell which docker-compose)
DOCKER_COMPOSE_RUN := $(DOCKER_COMPOSE) run --rm

.SILENT: help
help: ## Show this help message
	echo "usage: make [target] ..."
	echo ""
	echo "targets:"
	fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match $$'\t' | sed -e 's/: ## / - /'


.PHONY: fmt
fmt: ## Format the code
	$(DOCKER_COMPOSE_RUN) service.controller.dmx yapf --in-place --recursive --verbose .

.PHONY: clean
clean: ## Clean up any containers and images
	rm -r vendor/
	$(DOCKER_COMPOSE) stop
	$(DOCKER_COMPOSE) rm -f
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

.PHONY: loc
loc: ## Count lines of code
	cloc --exclude-dir=vendor,dist,.idea,.vscode --not-match-f="package-lock.json" .
