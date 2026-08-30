.PHONY: build test vet run

build:
	mkdir -p bin
	go -C yolodev build -o ../bin/yolodev ./cmd/yolodev

test:
	go -C yolodev test ./...

vet:
	go -C yolodev vet ./...

run: build
	./bin/yolodev theme edit themes/yolodev/placeholder.toml
