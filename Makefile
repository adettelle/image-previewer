BIN := "./bin"

.PHONY:test
test:
	go test -count 1 ./... -v

test-race:
	go test -race -count 100 ./... -v

lint: 
	golangci-lint run

vet:
	staticcheck ./...

check: lint vet test

run: 
	go run ./cmd/server/

build_previewer:
	go build -v -o "$(BIN)/build_previewer" ./cmd/

up-previewer:
	docker compose up previewer -d --build

down:
	docker compose down

integration-tests:
	docker compose
