PREFIX ?= /usr

BIN := dotclean
GO_SRC := $(shell find . -type f -name '*.go' ! -name '*_test.go')

DOC := docs/$(BIN).1
MAN_MD := $(DOC).md

clean:
	rm -f '$(BIN)'
	go clean -testcache

build: $(BIN)

$(BIN): $(GO_SRC) go.mod go.sum
	go build -o '$(BIN)' .

man: $(DOC)

$(DOC): $(MAN_MD)
	@command -v pandoc >/dev/null || { echo "pandoc required for make man"; exit 1; }
	pandoc '$(MAN_MD)' -s -t man -o '$(DOC)'

install: $(BIN) $(DOC)
	mkdir -p '$(DESTDIR)$(PREFIX)/bin'
	install -m 755 '$(BIN)' '$(DESTDIR)$(PREFIX)/bin/$(BIN)'

	mkdir -p '$(DESTDIR)$(PREFIX)/share/man/man1'
	install -m 644 '$(DOC)' '$(DESTDIR)$(PREFIX)/share/man/man1/$(notdir $(DOC))'

uninstall:
	rm -f '$(DESTDIR)$(PREFIX)/bin/$(BIN)'
	rm -f '$(DESTDIR)$(PREFIX)/share/man/man1/$(DOC)'

.PHONY: clean build man install uninstall
