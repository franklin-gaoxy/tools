# 可执行文件名
BINARY_NAME := hostinfo

# 源代码目录（main.go所在目录）
SRC_DIR := src

# 安装路径
INSTALL_PATH := /usr/bin

# 默认目标：构建项目
all: build

# 构建项目（进入src目录编译）
build:
	@echo "进入 $(SRC_DIR) 目录构建 $(BINARY_NAME)..."
	cd $(SRC_DIR) && go build -o ../$(BINARY_NAME) main.go
	@echo "构建完成：$(BINARY_NAME)"

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -f $(BINARY_NAME)
	@echo "清理完成"

# 安装程序到系统目录
install: build
	@echo "将 $(BINARY_NAME) 安装到 $(INSTALL_PATH)..."
	# 检查权限，必要时使用sudo
	if [ ! -w $(INSTALL_PATH) ]; then \
		echo "需要管理员权限，尝试使用sudo安装..."; \
		sudo cp $(BINARY_NAME) $(INSTALL_PATH)/; \
	else \
		cp $(BINARY_NAME) $(INSTALL_PATH)/; \
	fi
	@echo "安装完成，可直接使用 '$(BINARY_NAME)' 命令"

# 卸载程序
uninstall:
	@echo "卸载 $(BINARY_NAME)..."
	if [ ! -w $(INSTALL_PATH) ]; then \
		sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "卸载完成"

# 运行程序
run: build
	@echo "运行 $(BINARY_NAME)..."
	./$(BINARY_NAME)

# 声明伪目标
.PHONY: all build clean install uninstall run

