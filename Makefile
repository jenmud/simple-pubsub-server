PROJECT=myBetex
GO=/usr/local/go/bin/go
GOPATH=/Users/jenmud/Desktop/go-code/src/gitlab.com/jenmud/$(PROJECT)

GOINSTALL=$(GO) install
GOBUILD=$(GO) build
GORUN=$(GO) run

build: main.go
	$(GOBUILD) -o $(PROJECT) main.go

run:
	$(GORUN) *.go

clean:
	rm -rf $(PROJECT)