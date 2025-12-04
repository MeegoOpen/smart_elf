#!/bin/bash

# Smart Elf 构建脚本

echo "========================================"
echo "Smart Elf 独立插件构建脚本"
echo "========================================"

# 设置工作目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="${SCRIPT_DIR}/.."
cd "${PROJECT_DIR}"

echo "当前工作目录: $(pwd)"

# 设置Go环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

echo "Go版本: $(go version)"

# 创建输出目录
mkdir -p output

# 清理旧的构建文件
rm -f output/smart_elf

# 下载依赖
echo "正在下载依赖..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "依赖下载失败，请检查网络连接和go.mod文件"
    exit 1
fi

echo "正在验证依赖..."
go mod verify
if [ $? -ne 0 ]; then
    echo "依赖验证失败，请运行 'go mod tidy' 修复"
    exit 1
fi

# 运行测试
echo "正在运行测试..."
go test ./... -v
if [ $? -ne 0 ]; then
    echo "测试失败，请修复后再继续"
    exit 1
fi

# 构建项目
echo "正在构建项目..."
echo "构建目标: output/smart_elf"

go build \
    -ldflags "-s -w" \
    -o output/smart_elf \
    cmd/server/main.go

if [ $? -ne 0 ]; then
    echo "构建失败，请检查代码错误"
    exit 1
fi

# 检查构建结果
if [ -f "output/smart_elf" ]; then
    echo "构建成功!"
    echo "可执行文件: $(pwd)/output/smart_elf"
    echo "文件大小: $(du -h output/smart_elf | cut -f1)"
    
    # 复制启动脚本到输出目录
    cp scripts/start.sh output/
    chmod +x output/start.sh
    
    
    echo "========================================"
    echo "构建完成! 您可以使用以下命令启动服务:"
    echo "1. cd output && ./smart_elf"
    echo "2. ./scripts/start.sh"
    echo "========================================"
else
    echo "构建失败，未找到可执行文件"
    exit 1
fi