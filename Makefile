DOCKER_COMPOSE ?= $(shell which docker-compose)
DOCKER_COMPOSE_RUN := $(DOCKER_COMPOSE) run --rm

.SILENT: help
help: ## Show this help message
	echo "usage: make [target] ..."
	echo ""
	echo "targets:"
	fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match $$'\t' | sed -e 's/: ## / - /'

.PHONY: start
start: ## Start the system
	$(DOCKER_COMPOSE) up -d --renew-anon-volumes ## This flag avoids the npm modules bug https://github.com/docker/compose/pull/5596

.PHONY: start-platform
start-platform: ## Start the core services
	$(DOCKER_COMPOSE) up -d --renew-anon-volumes --remove-orphans filebeat logstash redis service.api-gateway service.config service.registry.device service.log

.PHONY: stop
stop: ## Stop the services
	$(DOCKER_COMPOSE) stop

.PHONY: test-service.registry.device
test-service.registry.device: ## Run tests for device registry
	$(DOCKER_COMPOSE_RUN) service.registry.device python -m unittest

.PHONY: test
test: ## Run all tests
	$(MAKE) test-service.registry.device

.PHONY: fmt
fmt: ## Format the code
	$(DOCKER_COMPOSE_RUN) service.registry.device yapf --in-place --recursive --verbose .
	$(DOCKER_COMPOSE_RUN) service.registry.device yapf --in-place --recursive --verbose /root/.local/lib/python3.5/site-packages/
	$(DOCKER_COMPOSE_RUN) service.controller.dmx yapf --in-place --recursive --verbose .

.PHONY: clean
clean: ## Clean up any containers and images
	rm -r vendor/
	$(DOCKER_COMPOSE) stop
	$(DOCKER_COMPOSE) rm -f
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans
