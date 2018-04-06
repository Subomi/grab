GO = go
GOGET = $(GO) get -u
INSTALL_SCRIPT = ./install.sh

build: 
	sh $(INSTALL_SCRIPT)

install:
	cd ./cmd/grab
	$(GO) install -v
