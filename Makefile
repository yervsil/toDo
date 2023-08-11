run:
	docker-compose up

test:
	go test -v ./internal/delivery/http ./internal/repository