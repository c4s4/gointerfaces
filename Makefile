NAME=gointerfaces
BUILD_DIR=build
GO_VERSION=1.0.3 1.1.2 1.2.2 1.3.3 1.4.3 1.5.4 1.6.4 1.7.4 1.8beta2

YELLOW=\033[1m\033[93m
CYAN=\033[1m\033[96m
CLEAR=\033[0m

.PHONY: build

help:
	@echo "$(CYAN)help$(CLEAR)     Print this help page"
	@echo "$(CYAN)clean$(CLEAR)    Delete generated files"
	@echo "$(CYAN)build$(CLEAR)    Build executable binary"
	@echo "$(CYAN)articles$(CLEAR) Build articles"
	@echo "$(CYAN)release$(CLEAR)  Release project"

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)

build:
	@echo "$(YELLOW)Building project$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(NAME)

articles:
	@echo "$(YELLOW)Building articles$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go run gointerfaces.go $(GO_VERSION) > $(BUILD_DIR)/interfaces.md
	cp go-interfaces*.md $(BUILD_DIR)
	sed -i -e "/INTERFACES/{r $(BUILD_DIR)/interfaces.md" -e "d}" $(BUILD_DIR)/go-interfaces.md
	sed -i -e "/INTERFACES/{r $(BUILD_DIR)/interfaces.md" -e "d}" $(BUILD_DIR)/go-interfaces.en.md
	sed -i -e "s/UPDATE/$(shell date -I)/g" $(BUILD_DIR)/go-interfaces.md
	sed -i -e "s/UPDATE/$(shell date -I)/g" $(BUILD_DIR)/go-interfaces.en.md

release: clean build articles
	@echo "$(YELLOW)Releasing project$(CLEAR)"
	release
