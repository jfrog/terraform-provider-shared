TEST?=./...

default: build

build: fmtcheck
	go build ./...

test:
	@echo "==> Starting unit tests"
	go test $(TEST) -v -timeout=30s -parallel=4 -count=1

fmt:
	@echo "==> Fixing source code with gofmt..."
	@gofmt -s -w ./
	(command -v ${GOBIN}/goimports &> /dev/null || go get golang.org/x/tools/cmd/goimports) && ${GOBIN}/goimports -w .

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements..."
	@sh -c "find . -name '*.go' -not -name '*vendor*' -print0 | xargs -0 gofmt -l -s"

.PHONY: build fmt test
