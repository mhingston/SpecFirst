BINARY_NAME ?= specfirst
VERSION ?= dev
BUILD_DIR ?= dist

GOOS_LIST ?= darwin linux windows
GOARCH_LIST ?= amd64 arm64

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) .

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: install
install:
	go install -ldflags "-X main.version=$(VERSION)" ./...

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: dist
# Build cross-platform binaries into dist/
dist:
	@mkdir -p $(BUILD_DIR)
	@for os in $(GOOS_LIST); do \
		for arch in $(GOARCH_LIST); do \
			out=$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch; \
			if [ "$$os" = "windows" ]; then out=$$out.exe; fi; \
			echo "Building $$out"; \
			GOOS=$$os GOARCH=$$arch go build -ldflags "-X main.version=$(VERSION)" -o $$out .; \
		done; \
	done
