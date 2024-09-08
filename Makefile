BIN := resume-generator
CMD_DIR := cmd/resumegenerator/

all: run

run:
	go run $(CMD_DIR)*

build:
	go build -C $(CMD_DIR) -o $(BIN)
	mv $(CMD_DIR)$(BIN) ./

clean:
	go clean
	rm -f $(BIN)

.PHONY: all build run clean
