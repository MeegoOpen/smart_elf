#!/bin/bash

# Smart Elf 数据库初始化和启动脚本（使用SQLite）

echo "Smart Elf 启动脚本（读取 conf/config.yaml 配置）"
echo "数据库与飞书配置由 conf/config.yaml 提供，应用启动时自动迁移表结构（GORM AutoMigrate）。"

# 删除旧的数据库文件（可选，仅用于开发环境）
if [ -f ./smart_elf.db ]; then
    echo "检测到旧的数据库文件，正在删除..."
    rm -f ./smart_elf.db
fi

echo "准备启动应用程序..."

# 启动应用程序
go run cmd/server/main.go
