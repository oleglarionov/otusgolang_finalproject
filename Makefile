generate:
	go generate ./...

test:
	go test -race ./internal/...

run:
	docker-compose --file build/local/docker-compose.yaml up --build
