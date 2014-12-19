GO_VERSION=1.4
NAME=gointerfaces
BUILD_DIR=build
BLOG_MD=$(NAME).md
BLOG_DIR=../blog/txt/108.md
SITE_MD=$(NAME).en.md
SITE_XML=$(NAME).en.xml
SITE_DIR=../sweetohm/pages

YELLOW=\033[93m
CLEAR=\033[0m

all: clean blog site

interfaces:
	@echo "$(YELLOW)Generate interfaces list$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go run $(NAME).go $(GO_VERSION) > $(BUILD_DIR)/interfaces.md

blog: interfaces
	@echo "$(YELLOW)Generate blog article$(CLEAR)"
	cp $(BLOG_MD) $(BUILD_DIR)
	sed -i -e "/INTERFACES/{r $(BUILD_DIR)/interfaces.md" -e "d}" $(BUILD_DIR)/$(BLOG_MD)
	sed -i -e "s/VERSION/$(GO_VERSION)/" $(BUILD_DIR)/$(BLOG_MD)
	cp $(BUILD_DIR)/$(BLOG_MD) $(BLOG_DIR)

site: interfaces
	@echo "$(YELLOW)Generate site article$(CLEAR)"
	cp $(SITE_MD) $(BUILD_DIR)
	sed -i -e "/INTERFACES/{r $(BUILD_DIR)/interfaces.md" -e "d}" $(BUILD_DIR)/$(SITE_MD)
	sed -i -e "s/VERSION/$(GO_VERSION)/" $(BUILD_DIR)/$(SITE_MD)
	md2xml -a $(BUILD_DIR)/$(SITE_MD) > $(BUILD_DIR)/$(SITE_XML)
	cp $(BUILD_DIR)/$(SITE_XML) $(SITE_DIR)

clean:
	@echo "$(YELLOW)Clean generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)
