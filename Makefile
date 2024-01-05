# 单元测试
.PHONY: test
test:
	go test -race -cover ./...

.PHONY: lint
lint: golint golangci-lint

golangci-lint:
	golangci-lint run -c .github/linters/.golangci.yml --issues-exit-code 1 ./...

golint:
	golint -set_exit_status ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config .github/linters/.golangci.yml --fix

.PHONY: fmt
fmt:
	@goimports -l -w .

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) fmt
	@$(MAKE) tidy


