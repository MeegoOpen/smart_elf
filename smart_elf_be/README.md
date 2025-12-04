# Smart Elf 独立插件

这是一个基于 Gin 框架的 Smart Elf 独立插件项目，支持飞书机器人与 飞书项目 系统的集成，可以自动创建工单并提供群组管理功能。

## 功能特性

- 飞书机器人事件监听与处理
- 自动在 飞书项目 系统中创建工作项（工单）
- 配置管理（更新、查询配置）
- 签名生成和验证机制
- 群组自动创建与关联
- 自动反馈工单链接

## 项目结构

```
smart_elf_standalone/
├── cmd/
│   └── server/           # 应用入口
├── internal/             # 内部包
│   ├── handler/          # HTTP处理器
│   ├── model/            # 数据模型
│   └── service/          # 业务逻辑
├── pkg/                  # 可导出的包
│   ├── config/           # 配置管理
│   └── database/         # 数据库连接
├── conf/                 # 配置文件
├── scripts/              # 脚本文件
├── go.mod                # Go模块定义
└── README.md             # 项目说明
```

## 依赖管理

主要依赖：
- github.com/gin-gonic/gin - Web框架
- gorm.io/gorm - ORM框架
- gorm.io/driver/mysql - MySQL驱动
- github.com/larksuite/oapi-sdk-go/v3 - 飞书SDK
- github.com/larksuite/project-oapi-sdk-golang - 飞书项目SDK
- github.com/rs/zerolog - 日志库

## 安装与运行

1. 安装依赖：
```bash
go mod tidy
```

2. 配置数据库和飞书机器人信息（在conf/config.yaml中）

3. 运行服务：
```bash
go run cmd/server/main.go
```
