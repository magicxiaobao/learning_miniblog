# 定义变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get
BINARY_NAME=miniblog
MAIN_PATH=cmd/miniblog
GOFILES=$(shell find . -type f -name "*.go" -not -path "./vendor/*")

# 定义颜色输出
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all build test coverage lint clean run help

all: test build

## 构建
build:
	@echo "${GREEN}==> 构建二进制文件...${RESET}"
	$(GOBUILD) -o $(BINARY_NAME) -v ./$(MAIN_PATH)

## 运行
run:
	@echo "${GREEN}==> 运行应用...${RESET}"
	./$(BINARY_NAME) -c config.yaml

## 测试
test:
	@echo "${GREEN}==> 运行单元测试...${RESET}"
	$(GOTEST) -v ./...

## 生成测试覆盖率报告
coverage:
	@echo "${GREEN}==> 生成测试覆盖率报告...${RESET}"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

## 运行lint检查
lint:
	@echo "${GREEN}==> 运行 golangci-lint...${RESET}"
	golangci-lint run

## 格式化代码
fmt:
	@echo "${GREEN}==> 格式化代码...${RESET}"
	gofmt -s -w $(GOFILES)

## 清理构建文件
clean:
	@echo "${GREEN}==> 清理...${RESET}"
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out

## 安装依赖
deps:
	@echo "${GREEN}==> 更新依赖...${RESET}"
	$(GOCMD) mod tidy

## 数据库迁移
migrate:
	@echo "${GREEN}==> 初始化数据库...${RESET}"
	@echo "创建 miniblog 数据库..."
	mysql -h127.0.0.1 -uroot -p12345678 -e "CREATE DATABASE IF NOT EXISTS miniblog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	@echo "创建表结构..."
	mysql -h127.0.0.1 -uroot -p12345678 miniblog < mysql_schema.sql

## 压力测试
benchmark:
	@echo "${GREEN}==> 运行性能测试...${RESET}"
	$(GOTEST) -bench=. -benchmem ./...

## 帮助文档
help:
	@echo ''
	@echo 'Miniblog 项目开发辅助工具'
	@echo ''
	@echo '用法:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo '目标列表:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-15s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) 