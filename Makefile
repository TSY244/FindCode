# 定义变量
BINARY_NAME := FindCodeCommand
SERVER_NAME := FindCodeServer
DOCKER_IMAGE := mergechen/findcode
GIT_URL :=  # 可通过 make run GIT_URL=<your-url> 传入
BUILD_DIR := .
DIST_DIR := dist
COVERAGE_FILE := coverage.out
COMMAND_PATH := ./cmd/command/main.go
SERVER_PATH := ./cmd/default/main.go
GO_BUILD := go build -ldflags="-s -w" -o
RUN_COMMAND_PATH := script/run_command.sh

.PHONY: build_cmd run docker-build docker-up docker-down clean build_server run_server

# 1. 编译 Go 项目 (cmd/cmd.go)
build_cmd:
	@echo "正在编译 Go 项目..."
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(BUILD_DIR)/$(BINARY_NAME) $(COMMAND_PATH)
	@echo "编译完成，生成可执行文件: $(BUILD_DIR)/$(BINARY_NAME)"

# 2. 执行 run_with_git.sh（需传入 GIT_URL 参数或环境变量）
run:
	@if [ -z "$(GIT_URL)" ]; then \
        echo "错误: 必须提供 GIT_URL 参数（如: make run GIT_URL=xxx）"; \
        exit 1; \
    fi
	@chmod +x $(RUN_COMMAND_PATH)
	@echo "正在执行 run.sh，仓库 URL: $(GIT_URL)"
	./run.sh "$(GIT_URL)"

# 3. 构建 Docker 镜像（tag: mergechen/findcode）
docker-build:
	@echo "正在构建 Docker 镜像: $(DOCKER_IMAGE)"
	docker build -t $(DOCKER_IMAGE) .
	@echo "镜像构建完成: $(DOCKER_IMAGE)"

# 4. 启动 docker-compose
docker-up:
	@if [ -z "$(GIT_URL)" ]; then \
        echo "错误: 必须提供 GIT_URL 参数（如: make run GIT_URL=xxx）"; \
        exit 1; \
    fi
	@echo "启动 docker-compose 服务..."
	GIT_URL=$(GIT_URL) docker-compose up -d
	@echo "服务已启动（端口映射: 18080:8000）"

docker-down:
	@echo "停止 docker-compose 服务..."
	docker-compose down
	@echo "服务已停止!"

# 5. 清理生成的文件
clean:
	@echo "开始清理项目..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_FILE)
	@rm -f $(BINARY_NAME) $(SERVER_NAME)
	@docker-compose down -v --remove-orphans 2>/dev/null || true
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@go clean -cache -testcache -modcache
	@echo "清理完成!"

# 6. 编译 Go server 项目 (cmd/server.go)
build_server:
	@echo "正在编译 Go server 项目..."
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(BUILD_DIR)/$(SERVER_NAME) $(SERVER_PATH)
	@echo "编译完成，生成可执行文件: $(BUILD_DIR)/$(SERVER_NAME)"

run_server:
	@echo "正在编译 Go server 项目..."
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(BUILD_DIR)/$(SERVER_NAME) $(SERVER_PATH)
	@echo "编译完成，生成可执行文件: $(BUILD_DIR)/$(SERVER_NAME)"
	@echo "正在启动 Go server 项目..."
	@$(BUILD_DIR)/$(SERVER_NAME)
