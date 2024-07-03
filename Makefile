# Docker settings
DOCKER_CMD=docker run --rm -ti -w /build/pomo -v $$PWD:/build/pomo
DOCKER_IMAGE=pomo-build

# Version retrieval from git tags
VERSION ?= $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
	VERSION := UNKNOWN
endif

# Linker flags for embedding version info
LDFLAGS=\
	-X github.com/joaorufino/pomo/pkg/internal/version.Version=$(VERSION)

# Default target
.PHONY: \
	all \
	test \
	docs \
	pomo-build \
	readme \
	release \
	release-linux \
	release-darwin \
	clean

# All target to run default build and tests
all: default test

# Default build target
default: bin/pomo

# Build the main binary
bin/pomo: 
	cd cmd/pomo && \
	go build -ldflags '${LDFLAGS}' -o ../../$@

# Run tests and vet
test:
	go test ./...
	go vet ./...

# Build Docker image for build environment
pomo-build:
	docker build -t $(DOCKER_IMAGE) .

# Linux build targets
bin/pomo-linux: bin/pomo-$(VERSION)-linux-amd64

bin/pomo-$(VERSION)-linux-amd64: bin
	$(DOCKER_CMD) --env GOOS=linux --env GOARCH=amd64 $(DOCKER_IMAGE) go build -ldflags "${LDFLAGS}" -o $@

bin/pomo-$(VERSION)-linux-amd64.md5: bin/pomo-$(VERSION)-linux-amd64
	md5sum bin/pomo-$(VERSION)-linux-amd64 | sed -e 's/bin\///' > $@

# macOS build targets
bin/pomo-darwin: bin/pomo-$(VERSION)-darwin-amd64

bin/pomo-$(VERSION)-darwin-amd64: bin
	# Cross-compile for Darwin (macOS)
	$(DOCKER_CMD) --env GOOS=darwin --env GOARCH=amd64 --env CC=x86_64-apple-darwin15-cc --env CGO_ENABLED=1 $(DOCKER_IMAGE) go build -ldflags "${LDFLAGS}" -o $@

bin/pomo-$(VERSION)-darwin-amd64.md5: bin/pomo-$(VERSION)-darwin-amd64
	md5sum bin/pomo-$(VERSION)-darwin-amd64 | sed -e 's/bin\///' > $@

# Release targets
release-linux: bin/pomo-$(VERSION)-linux-amd64 bin/pomo-$(VERSION)-linux-amd64.md5

release-darwin: bin/pomo-$(VERSION)-darwin-amd64 bin/pomo-$(VERSION)-darwin-amd64.md5

release: release-linux release-darwin

# Documentation targets
docs: www/data/readme.json
	cd www && cp ../install.sh static/ && hugo -d ../docs

www/data/readme.json: www/data README.md
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > $@

www/data bin:
	mkdir -p $@

# Clean up build artifacts
clean:
	rm -rf bin/pomo* www/data/readme.json docs

# Utility target to create necessary directories
bin:
	mkdir -p bin

