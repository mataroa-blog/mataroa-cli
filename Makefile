.POSIX:
.SUFFIXES:

GO = go
RM = rm
INSTALL = install
GOLANGCILINT = golangci-lint
GOFLAGS =
PREFIX = /usr/local
BINDIR = bin

all: mata

mata:
	$(GO) build -o mata ./cmd/mata $(GOFLAGS)

clean:
	$(RM) -rf mata

lint:
	$(GOLANGCILINT) run

test:
	$(GO) test ./... -v

install:
	$(INSTALL) -d \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/

	$(INSTALL) -pm 0755 mata $(DESTDIR)$(PREFIX)/$(BINDIR)/

uninstall:
	$(RM) -f \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/mata

.PHONY: all mata clean install uninstall
