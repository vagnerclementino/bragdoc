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
	echo $(GOPATH)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.50.1
	golangci-lint run

.PHONY: run
run: build ##@application run application
	./$(BINARY_NAME)

.PHONY: clean
clean: ##@application clean binary
	if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi
	rm -f $(BINARY_NAME).zip
	rm -f $(BINARY_NAME).tar.gz

.PHONY: build
build: clean ##@application build application
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) -ldflags $(LDFLAGS) cmd/cli/main.go

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
	mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

.PHONY: tidy
tidy: ##@helper  ensures that the go.mod file matches the source code in the module
	go mod tidy -v

.PHONY: package
package: build ##@application creates packaged versions (zip, tar.gz) from the binary
	zip -r $(BINARY_NAME).zip $(BINARY_NAME)
	tar czf $(BINARY_NAME).tar.gz $(BINARY_NAME)

# ignore unknown commands
%:
    @:
