# go-ether-kit Makefile
# Go Library for EVM Networks Development

.PHONY: help test clean fmt lint vet coverage benchmark deps tools generate dev tag docs

# 默认目标
.DEFAULT_GOAL := help

# 项目变量
PROJECT_NAME := go-ether-kit
PACKAGE := github.com/guanzhenxing/go-ether-kit
VERSION := $(shell git describe --tags --always --dirty)

# 颜色定义
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
CYAN := \033[36m
RESET := \033[0m

## 📋 帮助信息
help: ## 显示帮助信息
	@echo "$(CYAN)$(PROJECT_NAME) - Go Library Makefile$(RESET)"
	@echo "$(YELLOW)Version: $(VERSION)$(RESET)"
	@echo ""
	@echo "$(GREEN)可用命令:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## 🔧 开发环境
deps: ## 安装项目依赖
	@echo "$(GREEN)正在安装项目依赖...$(RESET)"
	go mod tidy
	go mod download
	go mod verify

tools: ## 安装开发工具
	@echo "$(GREEN)正在安装开发工具...$(RESET)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev: deps tools ## 设置开发环境
	@echo "$(GREEN)开发环境设置完成!$(RESET)"

## 🧪 测试相关
test: ## 运行测试
	@echo "$(GREEN)正在运行测试...$(RESET)"
	go test -v ./...

test-short: ## 运行短测试
	@echo "$(GREEN)正在运行短测试...$(RESET)"
	go test -short -v ./...

test-race: ## 运行竞态检测测试
	@echo "$(GREEN)正在运行竞态检测测试...$(RESET)"
	go test -race -v ./...

coverage: ## 生成测试覆盖率报告
	@echo "$(GREEN)正在生成测试覆盖率报告...$(RESET)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)覆盖率报告已生成: coverage.html$(RESET)"

benchmark: ## 运行基准测试
	@echo "$(GREEN)正在运行基准测试...$(RESET)"
	go test -bench=. -benchmem ./...

## 🔍 代码质量
fmt: ## 格式化代码
	@echo "$(GREEN)正在格式化代码...$(RESET)"
	go fmt ./...
	goimports -w .

lint: ## 代码静态检查
	@echo "$(GREEN)正在进行代码静态检查...$(RESET)"
	golangci-lint run

vet: ## Go 代码审查
	@echo "$(GREEN)正在进行 Go 代码审查...$(RESET)"
	go vet ./...

check: fmt lint vet test ## 运行所有代码检查

## 📄 合约相关
generate: ## 生成合约绑定代码
	@echo "$(GREEN)正在生成合约绑定代码...$(RESET)"
	@if [ -f "contracts/erc20/IERC20.abi" ]; then \
		abigen --abi contracts/erc20/IERC20.abi --pkg erc20 --type IERC20 --out contracts/erc20/erc20.go; \
		echo "$(GREEN)ERC20 合约绑定已生成$(RESET)"; \
	else \
		echo "$(YELLOW)未找到 ERC20 ABI 文件，跳过生成$(RESET)"; \
	fi

## 🧹 清理相关
clean: ## 清理缓存和临时文件
	@echo "$(YELLOW)正在清理缓存文件...$(RESET)"
	rm -rf coverage.out coverage.html
	go clean -cache
	go clean -testcache
	@echo "$(GREEN)清理完成!$(RESET)"

## 📚 文档相关
docs: ## 启动文档服务器
	@echo "$(GREEN)正在启动文档服务器...$(RESET)"
	@echo "$(BLUE)访问地址: http://localhost:6060/pkg/$(PACKAGE)$(RESET)"
	godoc -http=:6060

## 🧹 维护相关
update: ## 更新所有依赖
	@echo "$(GREEN)正在更新依赖...$(RESET)"
	go get -u ./...
	go mod tidy

security: ## 运行安全检查
	@echo "$(GREEN)正在进行安全检查...$(RESET)"
	go list -json -m all | nancy sleuth 2>/dev/null || echo "$(YELLOW)Nancy 未安装，跳过安全检查$(RESET)"

## 🚀 版本管理
tag: ## 创建新标签 (使用方式: make tag VERSION=v1.0.0)
ifndef VERSION
	@echo "$(RED)错误: 请指定版本号，例如: make tag VERSION=v1.0.0$(RESET)"
	@exit 1
endif
	@echo "$(GREEN)正在创建标签 $(VERSION)...$(RESET)"
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "$(GREEN)标签 $(VERSION) 已创建并推送$(RESET)"

## 📊 项目信息
info: ## 显示项目信息
	@echo "$(CYAN)📊 项目信息$(RESET)"
	@echo "$(GREEN)项目名称:$(RESET) $(PROJECT_NAME)"
	@echo "$(GREEN)包路径:$(RESET) $(PACKAGE)"
	@echo "$(GREEN)当前版本:$(RESET) $(VERSION)"
	@echo "$(GREEN)Go 版本:$(RESET) $$(go version | awk '{print $$3}')"
	@echo "$(GREEN)代码行数:$(RESET) $$(find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1 | awk '{print $$1}')"

## 🎯 常用组合命令
all: deps fmt lint vet test ## 完整的开发流程

ci: deps fmt lint vet test coverage ## CI/CD 流程

quick: fmt test ## 快速测试
