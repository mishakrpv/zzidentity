SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

TAG_NAME := $(shell git describe --abbrev=0 --tags --exact-match)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))

BIN_NAME := zzidentity
CODENAME ?= dunno

DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

# Default build target
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

#? dist: Create the "dist" directory
dist:
	mkdir -p dist

.PHONY: binary
#? binary: Build the binary
binary: dist
	@echo SHA: $(VERSION) $(CODENAME) $(DATE)
	CGO_ENABLED=0 GOGC=off GOOS=${GOOS} GOARCH=${GOARCH} go build -installsuffix nocgo -o ./dist/${GOOS}/${GOARCH}/$(BIN_NAME) ./cmd/$(BIN_NAME)

.PHONY: debug
#? debug: Run Delve
debug:
	dlv debug cmd/$(BIN_NAME)/$(BIN_NAME).go

.PHONY: run
#? run: Run the application
run:
	export CONFIG_FILE="$(shell pwd)/settings.yaml"; \
	export DOTENV_FILE="$(shell pwd)/.env"; \
	go run cmd/$(BIN_NAME)/$(BIN_NAME).go

.PHONY: lint
#? lint: Run golangci-lint
lint:
	golangci-lint run

.PHONY: fmt
#? fmt: Format the Code
fmt:
	gofmt -s -l -w $(SRCS)

.PHONY: help
#? help: Get more info on make commands
help: Makefile
	@echo "Choose a command run:"
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sort | sed -e 's/^/ /'
