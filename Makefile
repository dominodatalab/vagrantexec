all: vet lint test

.PHONY: deps
deps:
	@go get -u golang.org/x/lint/golint

.PHONY: vet
vet:
	@echo "==== go vet ===="
	@go vet ./...

.PHONY: lint
lint:
	@echo "==== go lint ===="
	@golint -set_exit_status ./...

.PHONY: test
test:
	@echo "==== go test ===="
	@go test -v ./...
