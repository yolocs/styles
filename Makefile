.PHONY: build fmt-check test vet run

build:
	mkdir -p bin
	go -C yolodev build -o ../bin/yolodev ./cmd/yolodev

fmt-check:
	@unformatted="$$(gofmt -l $$(find yolodev -type f -name '*.go'))"; \
	if [ -n "$$unformatted" ]; then \
		echo "Go files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

test:
	go -C yolodev test ./...

vet:
	go -C yolodev vet ./...

run: build
	./bin/yolodev theme edit themes/yolodev/placeholder.toml
