BINARY_NAME=voca
INSTALL_DIR=$(HOME)/.local/bin

.PHONY: all build clean install uninstall

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/voca

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Done!"

uninstall:
	@echo "Removing from $(INSTALL_DIR)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Done!"
