include Makefile.vars
include Makefile.docs


.SILENT:
.DEFAULT_GOAL := help

.PHONY: generate
generate: ##@application generate code from SQLC
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

.PHONY: test
test: clean ##@quality run tests with coverage
	go test ./... -v -coverprofile=coverage.txt -covermode=atomic

.PHONY: test-race
test-race: ##@quality validate race condition
	go test -race ./...

.PHONY: lint
lint: ##@quality check coding style
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	golangci-lint run

.PHONY: run
run: build ##@application run application
	./$(BINARY_NAME)

.PHONY: clean
clean: ##@application clean binary and temporary files
	@echo "🧹 Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME).zip
	@rm -f $(BINARY_NAME).tar.gz
	@rm -f coverage.txt
	@rm -rf dist/
	@echo "✅ Clean complete"

.PHONY: build
build: clean generate ##@application build application
	env CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) -ldflags $(LDFLAGS) cmd/cli/main.go

.PHONY: fmt
fmt: ##@quality run go fmt
	go fmt ./...

.PHONY: vet
vet: ##@quality Run go vet
	go vet ./...

.PHONY: imports
imports: ##@quality Run goimports
	goimports -w .

.PHONY: quality
quality:test test-race fmt vet imports lint ##@quality run all quality targets

.PHONY: install
install: build ##@application install local version
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin/..."
	sudo mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Installation complete! Run 'bragdoc --help' to get started."

.PHONY: tidy
tidy: ##@helper  ensures that the go.mod file matches the source code in the module
	go mod tidy -v

.PHONY: package
package: build ##@application creates packaged versions (zip, tar.gz) from the binary
	zip -r $(BINARY_NAME).zip $(BINARY_NAME)
	tar czf $(BINARY_NAME).tar.gz $(BINARY_NAME)

.PHONY: smoke
smoke: ##@quality run smoke tests to verify core functionality
	@echo "🔥 Running smoke tests..."
	@./smoke.sh $(TARGET_OS)

# ignore unknown commands
%:
    @:
