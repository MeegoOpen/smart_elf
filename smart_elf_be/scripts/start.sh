#!/bin/bash

# Smart Elf 启动脚本

echo "========================================"
echo "Smart Elf 独立插件启动脚本"
echo "========================================"

# 设置工作目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="${SCRIPT_DIR}/.."
cd "${PROJECT_DIR}"

echo "当前工作目录: $(pwd)"

# 检查是否已安装依赖
if [ ! -f "go.sum" ]; then
    echo "正在下载依赖..."
    go mod tidy
    if [ $? -ne 0 ]; then
        echo "依赖下载失败，请检查网络连接和go.mod文件"
        exit 1
    fi
    echo "依赖下载完成"
fi


# 编译项目
echo "正在编译项目..."
go build -o output/smart_elf cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "编译失败，请检查代码错误"
    exit 1
fi
echo "编译完成"

# 创建输出目录
mkdir -p output

# 启动服务
echo "正在启动Smart Elf服务..."
echo "服务地址: http://${SERVER_HOST:-0.0.0.0}:${SERVER_PORT:-8081}"
echo "健康检查: http://${SERVER_HOST:-0.0.0.0}:${SERVER_PORT:-8081}/health"
echo "按 Ctrl+C 停止服务"
echo "========================================"

# 启动编译后的程序
./output/smart_elf