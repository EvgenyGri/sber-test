export GO111MODULE=on

GOLANGCI_BIN:=golangci-lint

build:
	$(info #Build...)
	cd $(CURDIR)/cmd && go build -o ../bin/sber-test

lint:
	$(info #Lint...)
	$(GOLANGCI_BIN) run --config=.golangci-lint.yaml ./...

deps:
	$(info #Install dependencies...)
	go mod download


test:
	$(info #Running tests...)
	go test ./...


