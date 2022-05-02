.DEFAULT_GOAL := default
default: build test

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

goimportscheck:
	@go mod download

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test: goimportscheck
	@echo "executing tests..."
	@go test -v -count 1 .

build: goimportscheck errcheck vet
	@go build .

docker-build:
	@docker build --no-cache -t iandri/graphsec:v0.0.1 .

docker-run:
	@docker run --rm -it iandri/graphsec:v0.0.1 check vm -n VM_2
