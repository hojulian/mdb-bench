name?=example-run-1

.PHONY: bench
bench: ## Run benchmark
	$(call print-target)
	go run bench/main.go --url http://localhost:8080 --freq 5000 --dur 3m --name $(name) --attack --save
	vegeta plot --title=$(name) results/$(name)/results.json > results/$(name)/plot.html

.PHONY: run-backend
run-backend: ## Run benchmark backend
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

.PHONY: run-microdb
run-microdb: ## Run microdb
	$(call print-target)
	docker-compose -f docker/microdb/docker-compose.yaml up --force-recreate --remove-orphans

.PHONY: run-app-mysql
run-app-mysql: ## Run all shipping services (mysql mode)
	$(call print-target)
	docker-compose -f docker/shipping/mysql.docker-compose.yaml up --force-recreate --remove-orphans -d --scale handling=3 --scale tracking=5 --scale booking=5

.PHONY: run-app-mysql-cluster
run-app-mysql-cluster: ## Run all shipping services (mysql-cluster mode)
	$(call print-target)
	docker-compose -f docker/shipping/mysql-cluster.docker-compose.yaml up --force-recreate --remove-orphans -d --scale handling=3 --scale tracking=5 --scale booking=5

.PHONY: run-app-microdb
run-app-microdb: ## Run all shipping services (microdb mode)
	$(call print-target)
	docker-compose -f docker/shipping/microdb.docker-compose.yaml up --force-recreate --remove-orphans -d

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

.PHONY: build-publisher
build-publisher: ## build publisher
	$(call print-target)
	docker build -t microdb/benchmark:publisher -f docker/microdb/Dockerfile.publisher .

.PHONY: build-querier
build-querier: ## build querier
	$(call print-target)
	docker build -t microdb/benchmark:querier -f docker/microdb/Dockerfile.querier .

.PHONY: build-dataorigin
build-dataorigin: ## build dataorigin
	$(call print-target)
	docker build -t microdb/benchmark:dataorigin -f docker/microdb/Dockerfile.dataorigin .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
