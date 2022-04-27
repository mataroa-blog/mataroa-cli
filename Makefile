.POSIX:
.SUFFIXES:

GO = go
RM = rm
INSTALL = install
PANDOC = pandoc
GOLANGCILINT = golangci-lint
GOFLAGS =
PREFIX = /usr/local
BINDIR = bin
MANDIR = share/man

all: mata docgen doc/mata-config.5

mata:
	$(GO) build -o mata cmd/mata/main.go $(GOFLAGS)

docgen:
	mkdir -p doc/result
	$(GO) run cmd/docgen/main.go

doc/mata-config.5: doc/mata-config.5.md
	mkdir -p doc/result
	$(PANDOC) doc/mata-config.5.md -s -t man -o doc/result/mata-config.5

clean:
	$(RM) -rf mata doc/result

lint:
	$(GOLANGCILINT) run

test:
	$(GO) test ./... -v

install:
	$(INSTALL) -d \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man5/ \

	$(INSTALL) -pm 0755 mata $(DESTDIR)$(PREFIX)/$(BINDIR)/
	$(INSTALL) -pm 0644 doc/result/*.1 $(DESTDIR)$(PREFIX)/$(MANDIR)/man1/
	$(INSTALL) -pm 0644 doc/result/*.5 $(DESTDIR)$(PREFIX)/$(MANDIR)/man5/

uninstall:
	$(RM) -f \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/mata \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/* \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man5/*

.PHONY: all mata clean install uninstall
