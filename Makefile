BINARY_NAME?=voc
PREFIX?=$(HOME)/.local
bindir?=$(PREFIX)/bin
datadir?=$(PREFIX)/share/voc
KAIKKI_URL?=https://kaikki.org/frwiktionary/Fran%C3%A7ais/kaikki.org-dictionary-Fran%C3%A7ais.jsonl.gz

LDFLAGS=-ldflags "-X 'voc/internal/dictionary.DefaultDictionaryDirectory=$(datadir)' -X 'voc/internal/database.DefaultUserDBPath=$(datadir)/voc.db' -X 'voc/internal/dictionary.KaikkiURL=$(KAIKKI_URL)'"

.PHONY: all build clean install uninstall

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -tags "fts5" -o $(BINARY_NAME) ./cmd/voc

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


