BIN := resume-generator
CMD_DIR := cmd/resumegenerator/
TESTS_DIR := ./tests/...

all: run

run:
	go run $(CMD_DIR)*

build:
	go build -C $(CMD_DIR) -o $(BIN)
	mv $(CMD_DIR)$(BIN) ./

test:
	go test $(TESTS_DIR)

clean:
	go clean
	rm -f $(BIN)

.PHONY: all build run clean test
