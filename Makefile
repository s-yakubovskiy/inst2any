.PHONY: all build clean test

# Binary name
BINARY = inst2vk

# Destination
DESTDIR = ./bin

# Gopath
GOPATH = $(shell go env GOPATH)

# Determine the version
VERSION = $(shell git describe --tags --always --dirty)

# Pass in link time variables to set version
# LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
LDFLAGS=-a -tags netgo -ldflags '-w -extldflags "-static"'

# Builds the project
build:
	go build $(LDFLAGS) -o $(DESTDIR)/$(BINARY) ./cmd/inst2vk
	# CGO_ENABLED=1 go build $(LDFLAGS) -o $(DESTDIR)/$(BINARY)_static_cgo ./cmd/inst2vk

# Installs our project: copies binaries
install:
	go install $(LDFLAGS)

# Cleans our project: deletes binaries
clean:
	if [ -f $(DESTDIR)/$(BINARY) ] ; then rm $(DESTDIR)/$(BINARY) ; fi

# Runs tests
test:
	go test ./...

# Runs vet
vet:
	go vet ./...

# Formats the code
fmt:
	go fmt ./...

docker:
	docker build -f Dockerfile . -t yharwyn/private:inst2vk 	
	docker push yharwyn/private:inst2vk


# Run all checks
all: fmt vet test build

# Help info
help:
	@echo "make build - build the project"
	@echo "make install - install the project"
	@echo "make clean - clean the project"
	@echo "make test - run tests"
	@echo "make vet - run go vet"
	@echo "make fmt - run go fmt"
	@echo "make all - run fmt, vet, test and build"
	@echo "make help - display this help information"
