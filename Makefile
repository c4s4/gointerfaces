NAME=gointerfaces
VERSION=1.3.0
BUILD_DIR=build

YELLOW=\033[93m
CLEAR=\033[0m

.PHONY: build

help:
	@echo "$(YELLOW)help$(CLEAR)    Print this help page"
	@echo "$(YELLOW)clean$(CLEAR)   Delete generated files"
	@echo "$(YELLOW)build$(CLEAR)   Build executable binary"
	@echo "$(YELLOW)release$(CLEAR) Release project"

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)

build:
	@echo "$(YELLOW)Building project$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(NAME)

release: clean build
	@echo "$(YELLOW)Releasing project$(CLEAR)"
	@if [ `git rev-parse --abbrev-ref HEAD` != "develop" ]; then \
		echo "You must release on branch develop"; \
		exit 1; \
	fi
	git diff --quiet --exit-code HEAD || (echo "There are uncommitted changes"; exit 1)
	git checkout master
	git merge develop
	git push
	git tag "$(VERSION)"
	git push --tag
	git checkout develop

