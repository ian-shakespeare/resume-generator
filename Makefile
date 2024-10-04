CMD_DIR := ./cmd/

MAIN_FILE := resumegen.go
BIN := resumegen

TEST_DIRS := \
	./test/generator/... \
	./test/cli/... \
	./test/assert/...

all: run

run:
	@go run $(CMD_DIR)$(MAIN_FILE)

build:
	@go build $(CMD_DIR)$(MAIN_FILE)
	@chmod +x $(BIN)

test:
	@go test $(TEST_DIRS)

clean:
	@go clean
	@rm -f $(BIN)

.PHONY: all run build test clean
