.PHONY: run-mysql
run-mysql:
	$(call print-target)
	docker-compose -f docker/mysql/docker-compose.yaml up --force-recreate --remove-orphans

.PHONY: run-shipping
run-shipping:
	$(call print-target)
	docker-compose -f docker/shipping/docker-compose.yaml up --force-recreate --remove-orphans

.PHONY: build-shipping
build-shipping: ## build shipping services
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
