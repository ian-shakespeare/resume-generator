CMD_DIR := ./cmd/

WEB_DIR := $(CMD_DIR)resumegenweb/
WEB_BIN := resumegenweb

CLI_DIR := $(CMD_DIR)resumegencli/
CLI_BIN := resumegen

BUILD_CMD := go build
DIR_FLAG := -C
OUT_FLAG := -o
CHMOD := chmod +x

TAILWIND_CMD := npx tailwindcss -i ./static/main.css -o ./static/output.css

TEST_DIRS := \
	./test/database/... \
	./test/generator/... \
	./test/cli/...

all: run

run:
	@cd web; $(TAILWIND_CMD)
	go run $(WEB_DIR)resumegenweb.go

web:
	@cd web; $(TAILWIND_CMD)
	$(BUILD_CMD) $(DIR_FLAG) $(WEB_DIR) $(OUT_FLAG) $(WEB_BIN)
	@mv $(WEB_DIR)$(WEB_BIN) ./$(WEB_BIN)
	@$(CHMOD) $(WEB_BIN)

cli:
	$(BUILD_CMD) $(DIR_FLAG) $(CLI_DIR) $(OUT_FLAG) $(CLI_BIN)
	@mv $(CLI_DIR)$(CLI_BIN) ./$(CLI_BIN)
	@$(CHMOD) $(CLI_BIN)

test:
	go test $(TEST_DIRS)

clean:
	go clean
	@rm -f $(WEB_BIN) $(CLI_BIN)

.PHONY: all run web cli test clean
