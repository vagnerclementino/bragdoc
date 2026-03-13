include Makefile.vars
include Makefile.docs

.SILENT:
.DEFAULT_GOAL := help

.PHONY: test
test: clean ##@quality run tests with coverage
	go test ./... -v -coverprofile=coverage.txt -covermode=atomic

.PHONY: test-race
test-race: ##@quality validate race condition
	go test -race ./...

.PHONY: lint
lint: ##@quality check coding style
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.62.2
	golangci-lint run

.PHONY: run
run: build ##@application run application
	./$(BINARY_NAME)

.PHONY: clean
clean: ##@application clean binary and artifacts
	if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi
	rm -rf .coverdata
	rm -f $(BINARY_NAME)-test
	rm -f coverage.txt
	rm -f $(BINARY_NAME).zip
	rm -f $(BINARY_NAME).tar.gz

.PHONY: build
build: generate ##@application build application
	rm -f $(BINARY_NAME)
	env CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) -ldflags $(LDFLAGS) ./cmd/cli
	chmod +x $(BINARY_NAME)

.PHONY: fmt
fmt: ##@quality run go fmt
	go fmt ./...

.PHONY: vet
vet: ##@quality run go vet
	go vet ./...

.PHONY: imports
imports: ##@quality run goimports
	goimports -w .

.PHONY: quality
quality: test test-race fmt vet imports lint ##@quality run all quality targets

.PHONY: install
install: build ##@application install local version
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

.PHONY: tidy
tidy: ##@helper ensures that the go.mod file matches the source code
	go mod tidy -v

.PHONY: package
package: build ##@application creates packaged versions (zip, tar.gz)
	zip -r $(BINARY_NAME).zip $(BINARY_NAME)
	tar czf $(BINARY_NAME).tar.gz $(BINARY_NAME)

.PHONY: generate
generate: ##@helper generate SQLC code
	sqlc generate

.PHONY: smoke
smoke: build ##@quality run smoke tests
	./smoke.sh

.PHONY: update-golden
update-golden: ##@helper update golden files
	go test ./cmd/cli -update

.PHONY: release
release: clean ##@application create a new release (usage: make release VERSION=v1.0.0)
	@if [ -z "$(filter-out $(VERSION),0.1.0)" ] || [ "$(VERSION)" = "0.1.0" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "✅ Release $(VERSION) created! GitHub Actions will build binaries."

# ignore unknown commands
%:
	@:
