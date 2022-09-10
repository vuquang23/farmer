bitnami-infra-up:
	docker compose -f ./infra/docker-compose-bitnami-infra.yaml -p farmer up -d

infra-up:
	docker compose -f ./infra/docker-compose-infra.yaml -p farmer up -d

sfarmer-up:
	docker compose -f ./infra/docker-compose-sfarmer.yaml -p farmer up -d

log-infra-up:
	docker compose -f ./infra/docker-compose-log.yaml -p farmer up -d

infra-down:
	docker compose -p farmer down

sfarmer-test-collect-log:
	@mkdir -p ./infra/volumes/log
	@touch ./infra/volumes/log/sfarmer.log
	go run ./cmd/main.go sfarmer &> ./infra/volumes/log/sfarmer.log
