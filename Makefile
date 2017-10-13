GO = go

PACKAGES = ./cmd/... ./pkg/...

SOURCES = $(shell find cmd pkg -type f -name "*.go")

TARGET = ./cmd/casbin-demo/server

all: build test

.PHONY: clean
clean:
	-rm -vrf $(TARGET) docs

BINDATA_TARGET = ./pkg/conf/model.go
bindata:
	go-bindata -o $(BINDATA_TARGET) -pkg conf config

.PHONY: imports
imports:
	goimports -w -l $(SOURCES)

.PHONY: build
build: $(TARGET)

$(TARGET): $(SOURCES)
	$(GO) build -o $(TARGET) ./cmd/casbin-demo

.PHONY: test
test:
	-echo 'hello test'
