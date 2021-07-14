.PHONY: list
list:
	@sh -c "$(MAKE) -p no_targets__ 2>/dev/null | \
	awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | \
	grep -v Makefile | \
	grep -v '%' | \
	grep -v '__\$$' | \
	sort -u"

.PHONY: format
format :
	goimports -w .
	gofmt -w -s .

.PHONY: lint
lint :
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.39.0 golangci-lint run

.PHONY: unit
unit:
	go test --cover ./...