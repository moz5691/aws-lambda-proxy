.PHONY: build clean deploy

GOCMD=go 
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
SERVICE=puppies
BINARY_NAME=main
CLIENT_BINARY_NAME=client

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o build/$(SERVICE) cmd/$(SERVICE)server/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gen:
	protoc --proto_path=$(GOPATH)/src:. --twirp_out=. --go_out=. rpc/$(SERVICE)/service.proto

dev-build:
	$(GOBUILD) -o build/$(BINARY_NAME) -v cmd/$(SERVICE)server/main.go

dev-run:
	$(GOBUILD) -o build/$(BINARY_NAME) -v cmd/$(SERVICE)server/main.go
	build/$(BINARY_NAME)

dev-crun:
	$(GOBUILD) -o build/$(CLIENT_BINARY_NAME) -v cmd/$(SERVICE)client/main.go
	build/$(CLIENT_BINARY_NAME)

test:
	$(GOTEST) -v ./cmd/$(SERVICE)services