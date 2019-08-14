 # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOFMT=$(GOCMD) fmt
    BINARY_NAME=kvcache
    DEFAULT_PORT=11211

build:
	$(GOBUILD) -o $(BINARY_NAME)

test:
	$(GOTEST) -v ./broker ./cache ./verb_worker ./server

fmt:
	$(GOFMT) ./...

run:
	$(GOBUILD) -o $(BINARY_NAME)
	./$(BINARY_NAME) $(DEFAULT_PORT)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
