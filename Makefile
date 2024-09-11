BIN := resume-generator
CMD_DIR := cmd/resumegenerator/

TEST_DIRS := \
	./test/database/... \
	./test/handlers/... \
	./test/generator/...

all: run

run:
	go run $(CMD_DIR)*

build:
	go build -C $(CMD_DIR) -o $(BIN)
	mv $(CMD_DIR)$(BIN) ./

test:
	go test $(TEST_DIRS)

clean:
	go clean
	rm -f $(BIN)

.PHONY: all build run clean test
