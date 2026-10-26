# terminu-harness — build / test / check
# 跑 make（或 make help）列出所有指令。

BINARY := terminu-harness

.DEFAULT_GOAL := help

.PHONY: help
help: ## 列出所有指令
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-8s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## 建 MCP server 到 .local/bin（plugin 的 scripts/mcp.sh 會優先用它，比原始碼舊時不用）
	go build -o .local/bin/$(BINARY) ./cmd/$(BINARY)

.PHONY: test
test: ## 跑測試；設 TERMINU_HARNESS__TDP_REPO 指到 terminu 的 clone 時也驗真的清單
	go test ./...

.PHONY: check
check: ## gofmt、go vet、測試、claude plugin validate
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	go vet ./...
	go test ./...
	claude plugin validate .

.PHONY: clean
clean: ## 刪掉 .local/bin
	rm -rf .local/bin
