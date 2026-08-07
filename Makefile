.PHONY: clean build format test man install uninstall

PREFIX ?= /usr
BIN := dotclean
DOC := $(BIN).1

build:
	go build -o $(BIN) .

format:
	go fmt ./...
	go mod tidy

test: format
	go test ./... | sed '/^[[:space:]]*?/d'

man: docs/dotclean.1.md
	@command -v pandoc >/dev/null || { echo "pandoc required for make man"; exit 1; }
	pandoc docs/dotclean.1.md -s -t man -o 'docs/$(DOC)'

install: build
	mkdir -p '$(DESTDIR)$(PREFIX)/bin'
	install -m 755 '$(BIN)' '$(DESTDIR)$(PREFIX)/bin/$(BIN)'

	mkdir -p '$(DESTDIR)$(PREFIX)/share/man/man1'
	install -m 644 'docs/$(DOC)' '$(DESTDIR)$(PREFIX)/share/man/man1/$(DOC)'

uninstall:
	rm -f '$(DESTDIR)$(PREFIX)/bin/$(BIN)'
	rm -f '$(DESTDIR)$(PREFIX)/share/man/man1/$(DOC)'

clean:
	rm -f $(BIN)
