BIN := "./bin"

test:
	go test -race -count 100 ./... -v

lint: 
	golangci-lint run

vet:
	staticcheck ./...

check: lint vet test

run-server: 
	go run ./cmd/

build_previewer:
	go build -v -o "$(BIN)/build_previewer" ./cmd/