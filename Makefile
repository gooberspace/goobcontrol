get-deps:
	go get .
build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version=$(shell git describe --tags --always --long) -extldflags=-static" -tags timetzdata -o bin/goobcontrol_linux_amd64
build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s -X main.version=$(shell git describe --tags --always --long) -extldflags=-static" -tags timetzdata -o bin/goobcontrol_linux_arm64
build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.version=$(shell git describe --tags --always --long) -extldflags=-static" -tags timetzdata -o bin/goobcontrol_windows_amd64.exe

build-all: build-linux-amd64 build-linux-arm64 build-windows-amd64