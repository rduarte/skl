VERSION ?= dev
LDFLAGS = -s -w -X github.com/rduarte/skl/cmd.Version=$(VERSION)

.PHONY: build install clean

## Build the binary
build:
	go build -ldflags "$(LDFLAGS)" -o skl .

## Build and install to ~/.local/bin
install: build
	mkdir -p $(HOME)/.local/bin
	cp skl $(HOME)/.local/bin/skl
	@echo "✅ skl $(VERSION) instalado em $(HOME)/.local/bin/skl"

## Remove built binary
clean:
	rm -f skl
