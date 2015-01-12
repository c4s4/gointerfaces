NAME=gointerfaces
VERSION=1.2.0
GO_VERSION=1.0.3 1.1.2 1.2.2 1.3.3 1.4
BUILD_DIR=build
SITE=../sweetohm/pages/

YELLOW=\033[93m
CLEAR=\033[0m

all: clean articles publish

interfaces:
	@echo "$(YELLOW)Generate interfaces list$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go run $(NAME).go $(GO_VERSION) > $(BUILD_DIR)/interfaces.md


articles: interfaces
	@echo "$(YELLOW)Generate site article$(CLEAR)"
	# build french version
	cp article.md $(BUILD_DIR)
	sed -i -e "/INTERFACES/{r $(BUILD_DIR)/interfaces.md" -e "d}" $(BUILD_DIR)/article.md
	md2xml -a $(BUILD_DIR)/article.md > $(BUILD_DIR)/article.xml
	# build english version
	cp article.en.md $(BUILD_DIR)
	sed -i -e "/INTERFACES/{r $(BUILD_DIR)/interfaces.md" -e "d}" $(BUILD_DIR)/article.en.md
	md2xml -a $(BUILD_DIR)/article.en.md > $(BUILD_DIR)/article.en.xml

publish: articles
	@echo "$(YELLOW)Publish articles on site$(CLEAR)"
	cp $(BUILD_DIR)/article.xml $(SITE)/gointerfaces.xml
	cp $(BUILD_DIR)/article.en.xml $(SITE)/gointerfaces.en.xml

tag:
	@echo "$(YELLOW)Tagging for release $(VERSION)$(CLEAR)"
	git tag "$(VERSION)"
	git push --tags

release: clean tag

clean:
	@echo "$(YELLOW)Clean generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)

html: interfaces
	@echo "$(YELLOW)Generate HTML page$(CLEAR)"
	pandoc -f markdown -t html $(BUILD_DIR)/interfaces.md > $(BUILD_DIR)/index.html

help:
	@echo "$(YELLOW)interfaces$(CLEAR) Build the list of interfaces"
	@echo "$(YELLOW)articles$(CLEAR)   Build the articles"
	@echo "$(YELLOW)publish$(CLEAR)    Build articles and copy them in the site project"
	@echo "$(YELLOW)clean$(CLEAR)      Delete generated files"
	@echo "$(YELLOW)html$(CLEAR)       Generate HTML page"
	@echo "$(YELLOW)help$(CLEAR)       Print this help page"
