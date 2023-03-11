.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	@go build -o bin/ccommits_pls

