all: test

format:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test -count=1 ./...

validate: sort-import format vet lint

.PHONY: test format vet lint validate
