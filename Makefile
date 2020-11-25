# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=songs
GOPATH=./cmd/$(BINARY_NAME)
BINPATH=./bin

build: 
	mkdir -p $(BINPATH)
	$(GOBUILD) -o $(BINPATH)/$(BINARY_NAME) $(GOPATH) 


test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -rf $(BINPATH)

run: build
	$(BINPATH)/$(BINARY_NAME)

deps:
	$(GOGET) github.com/markbates/goth
	$(GOGET) github.com/markbates/pop
    
all: test build

