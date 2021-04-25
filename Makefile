.PHONY: run-backend
run-backend: ## Run bencmark backend
	$(call print-target)
	docker-compose -f docker/backend/docker-compose.yaml up -d --force-recreate --remove-orphans

.PHONY: run-mysql
run-mysql: ## Run mysql
	$(call print-target)
	docker-compose -f docker/mysql/docker-compose.yaml up --force-recreate --remove-orphans

.PHONY: run-mysql-cluster
run-mysql-cluster: ## Run mysql cluster
	$(call print-target)
	docker-compose -f docker/mysql-cluster/docker-compose.yaml up --force-recreate --remove-orphans

.PHONY: run-app-mysql
run-app-mysql: ## Run all shipping services (mysql mode)
	$(call print-target)
	docker-compose -f docker/shipping/mysql.docker-compose.yaml up --force-recreate --remove-orphans -d --scale handling=3 --scale tracking=3 --scale booking=3

.PHONY: run-app-mysql-cluster
run-app-mysql-cluster: ## Run all shipping services (mysql-cluster mode)
	$(call print-target)
	docker-compose -f docker/shipping/mysql-cluster.docker-compose.yaml up --force-recreate --remove-orphans -d --scale handling=3 --scale tracking=3 --scale booking=3

.PHONY: build-app
build-app: ## build shipping services
	$(call print-target)
	docker build -t microdb/benchmark:handling -f docker/shipping/Dockerfile.handling .
	docker build -t microdb/benchmark:booking -f docker/shipping/Dockerfile.booking .
	docker build -t microdb/benchmark:tracking -f docker/shipping/Dockerfile.tracking .

.PHONY: build-mysql
build-mysql: ## build mysql database
	$(call print-target)
	docker build -t microdb/benchmark:mysql -f docker/mysql/Dockerfile .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
