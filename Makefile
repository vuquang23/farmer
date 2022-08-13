infra-up:
	docker compose -f ./infra/docker-compose-infra.yaml up -d

infra-down:
	docker compose -f ./infra/docker-compose-infra.yaml down