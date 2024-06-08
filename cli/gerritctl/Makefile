.DEFAULT_GOAL := build
VERSION=0.0.1
SOURCE_GO=gerritctl.go
EXECUTABLE=gerritctl
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64

BIN_ALL_PLATFORMS=windows linux darwin

windows:
	rm -f $(EXECUTABLE)
	env GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)"  $(SOURCE_GO)

linux:
	rm -f $(EXECUTABLE)
	env GOOS=linux GOARCH=amd64 go build -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"  $(SOURCE_GO)

darwin:
	rm -f $(EXECUTABLE)
	env GOOS=darwin GOARCH=amd64 go build -v -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)" $(SOURCE_GO)

all: build

clean:
	rm -f $(BIN_ALL_PLATFORMS)

build: windows linux darwin ## Build binaries
	@echo version: $(VERSION)
