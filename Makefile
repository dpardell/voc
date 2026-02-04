BINARY_NAME?=voca
PREFIX?=$(HOME)/.local
bindir?=$(PREFIX)/bin
datadir?=$(PREFIX)/share/voca

LDFLAGS=-ldflags "-X 'voca/internal/dictionary.DefaultDictionaryPath=$(datadir)/dictionary.db' -X 'voca/internal/database.DefaultUserDBPath=$(datadir)/voca.db'"

.PHONY: all build clean install uninstall

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -tags "fts5" -o $(BINARY_NAME) ./cmd/voca

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)

install: build
	@echo "Installing binary to $(DESTDIR)$(bindir)..."
	mkdir -p $(DESTDIR)$(bindir)
	cp $(BINARY_NAME) $(DESTDIR)$(bindir)/
	@echo "Creating data directory at $(DESTDIR)$(datadir)..."
	mkdir -p $(DESTDIR)$(datadir)

uninstall:
	@echo "Removing from $(DESTDIR)$(bindir)..."
	rm -f $(DESTDIR)$(bindir)/$(BINARY_NAME)
	@echo "Done!"


