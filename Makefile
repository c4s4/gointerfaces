NAME=gointerfaces
BUILD_DIR=build

YELLOW=\033[1m\033[93m
CYAN=\033[1m\033[96m
CLEAR=\033[0m

.PHONY: build

help:
	@echo "$(CYAN)help$(CLEAR)    Print this help page"
	@echo "$(CYAN)clean$(CLEAR)   Delete generated files"
	@echo "$(CYAN)build$(CLEAR)   Build executable binary"
	@echo "$(CYAN)release$(CLEAR) Release project"

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)

build:
	@echo "$(YELLOW)Building project$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(NAME)

release: clean build
	@echo "$(YELLOW)Releasing project$(CLEAR)"
	release
