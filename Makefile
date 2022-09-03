infra-up:
	docker compose -f ./infra/docker-compose-infra.yaml -p farmer up -d

infra-down:
	docker compose -p farmer down

sfarmer-up:
	docker compose -f ./infra/docker-compose-sfarmer.yaml -p farmer up -d
