.POSIX:
.SUFFIXES:

GO = go
RM = rm
INSTALL = install
SCDOC = scdoc
GOLANGCILINT = golangci-lint
GOFLAGS =
PREFIX = /usr/local
BINDIR = bin
MANDIR = share/man

all: mata doc/mata.1 doc/mata-config.5

mata:
	$(GO) build $(GOFLAGS)

doc/mata-config.5: doc/mata-config.5.scd
	$(SCDOC) <doc/mata-config.5.scd >doc/mata-config.5

doc/mata.1: doc/mata.1.scd
	$(SCDOC) <doc/mata.1.scd >doc/mata.1

clean:
	$(RM) -f mata doc/mata.1 doc/mata-config.5

lint:
	$(GOLANGCILINT) run

test:
	$(GO) test ./...

install:
	$(INSTALL) -d \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/

	$(INSTALL) -pm 0755 mata $(DESTDIR)$(PREFIX)/$(BINDIR)/
	$(INSTALL) -pm 0644 doc/mata.1 $(DESTDIR)$(PREFIX)/$(MANDIR)/man1/
	$(INSTALL) -pm 0644 doc/mata-config.5 $(DESTDIR)$(PREFIX)/$(MANDIR)/man1/

uninstall:
	$(RM) -f \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/mata \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/mata.1

.PHONY: all mata clean install uninstall
