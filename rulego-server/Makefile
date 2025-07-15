# RuleGo Server Makefile

# 变量定义
BINARY_NAME=rulego-server
TOOLKIT_BINARY=toolkit
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# 编译标志
LDFLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"

# 构建标签
TAGS=

# 默认目标
.PHONY: all
all: build

# 编译服务器
.PHONY: build
build:
	@echo "编译 $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go
	@echo "编译完成: bin/$(BINARY_NAME)"

# 编译工具链
.PHONY: build-toolkit
build-toolkit:
	@echo "编译 $(TOOLKIT_BINARY)..."
	@go build -o bin/$(TOOLKIT_BINARY) cmd/toolkit/main.go
	@echo "编译完成: bin/$(TOOLKIT_BINARY)"

# 编译所有二进制文件
.PHONY: build-all
build-all: build build-toolkit

# 交叉编译
.PHONY: build-linux
build-linux:
	@echo "交叉编译 Linux 版本..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 cmd/server/main.go
	@GOOS=linux GOARCH=amd64 go build -o bin/$(TOOLKIT_BINARY)-linux-amd64 cmd/toolkit/main.go

.PHONY: build-windows
build-windows:
	@echo "交叉编译 Windows 版本..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe cmd/server/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/$(TOOLKIT_BINARY)-windows-amd64.exe cmd/toolkit/main.go

.PHONY: build-darwin
build-darwin:
	@echo "交叉编译 macOS 版本..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 cmd/server/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/$(TOOLKIT_BINARY)-darwin-amd64 cmd/toolkit/main.go

# 交叉编译所有平台
.PHONY: build-cross
build-cross: build-linux build-windows build-darwin

# 运行服务器
.PHONY: run
run: build
	@echo "启动服务器..."
	@./bin/$(BINARY_NAME)

# 运行带配置文件的服务器
.PHONY: run-config
run-config: build
	@echo "启动服务器 (使用配置文件)..."
	@./bin/$(BINARY_NAME) -config configs/config.yaml

# 测试
.PHONY: test
test:
	@echo "运行测试..."
	@go test -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "运行测试覆盖率..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告生成: coverage.html"

# 代码格式化
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "代码检查..."
	@golangci-lint run

# 清理
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.out coverage.html

# 创建发布包
.PHONY: package
package: build-cross
	@echo "创建发布包..."
	@mkdir -p dist
	@cp -r configs dist/
	@cp -r docs dist/
	@cp README.md dist/
	@cp LICENSE dist/ 2>/dev/null || true
	
	# Linux 包
	@mkdir -p dist/linux-amd64
	@cp bin/$(BINARY_NAME)-linux-amd64 dist/linux-amd64/$(BINARY_NAME)
	@cp bin/$(TOOLKIT_BINARY)-linux-amd64 dist/linux-amd64/$(TOOLKIT_BINARY)
	@cp -r configs dist/linux-amd64/
	@cd dist && tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz linux-amd64/
	
	# Windows 包
	@mkdir -p dist/windows-amd64
	@cp bin/$(BINARY_NAME)-windows-amd64.exe dist/windows-amd64/$(BINARY_NAME).exe
	@cp bin/$(TOOLKIT_BINARY)-windows-amd64.exe dist/windows-amd64/$(TOOLKIT_BINARY).exe
	@cp -r configs dist/windows-amd64/
	@cd dist && zip -r $(BINARY_NAME)-$(VERSION)-windows-amd64.zip windows-amd64/
	
	# macOS 包
	@mkdir -p dist/darwin-amd64
	@cp bin/$(BINARY_NAME)-darwin-amd64 dist/darwin-amd64/$(BINARY_NAME)
	@cp bin/$(TOOLKIT_BINARY)-darwin-amd64 dist/darwin-amd64/$(TOOLKIT_BINARY)
	@cp -r configs dist/darwin-amd64/
	@cd dist && tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz darwin-amd64/
	
	@echo "发布包创建完成"

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	@go mod tidy
	@go mod download

# 更新依赖
.PHONY: update-deps
update-deps:
	@echo "更新依赖..."
	@go get -u ./...
	@go mod tidy

# 生成文档
.PHONY: docs
docs:
	@echo "生成文档..."
	@godoc -http=:6060 &
	@echo "文档服务器启动: http://localhost:6060"

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "构建 Docker 镜像..."
	@docker build -t rulego-server:$(VERSION) .
	@docker tag rulego-server:$(VERSION) rulego-server:latest

# Docker 运行
.PHONY: docker-run
docker-run:
	@echo "运行 Docker 容器..."
	@docker run -p 8080:8080 -v $(PWD)/configs:/app/configs rulego-server:latest

# 显示版本信息
.PHONY: version
version:
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git 提交: $(GIT_COMMIT)"

# 初始化项目
.PHONY: init
init:
	@echo "初始化项目..."
	@mkdir -p bin configs data logs plugins
	@go mod tidy

# 帮助信息
.PHONY: help
help:
	@echo "RuleGo Server Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  build           编译服务器"
	@echo "  build-toolkit   编译工具链"
	@echo "  build-all       编译所有二进制文件"
	@echo "  build-cross     交叉编译所有平台"
	@echo "  run             运行服务器"
	@echo "  run-config      使用配置文件运行服务器"
	@echo "  test            运行测试"
	@echo "  test-coverage   运行测试覆盖率"
	@echo "  fmt             格式化代码"
	@echo "  lint            代码检查"
	@echo "  clean           清理构建文件"
	@echo "  package         创建发布包"
	@echo "  deps            安装依赖"
	@echo "  update-deps     更新依赖"
	@echo "  docs            生成文档"
	@echo "  docker-build    构建 Docker 镜像"
	@echo "  docker-run      运行 Docker 容器"
	@echo "  version         显示版本信息"
	@echo "  init            初始化项目"
	@echo "  help            显示此帮助信息"