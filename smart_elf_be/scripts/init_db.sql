-- Smart Elf 数据库初始化脚本

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS smart_elf DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE smart_elf;

-- 创建smart_elf表（对应AppConfig模型）
CREATE TABLE IF NOT EXISTS smart_elf (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    bot_id VARCHAR(255) NOT NULL,
    bot_secret VARCHAR(255) NOT NULL,
    bot_verification_token VARCHAR(255),
    project_key VARCHAR(255) NOT NULL,
    tenant_key VARCHAR(255),
    work_item_type_key VARCHAR(255),
    work_item_api_name VARCHAR(255),
    work_item_template_id BIGINT,
    creator_field_key VARCHAR(255),
    reply_switch BOOLEAN DEFAULT FALSE,
    create_group_switch BOOLEAN DEFAULT FALSE,
    signature VARCHAR(255),
    api_user_key VARCHAR(255),
    INDEX idx_project_key (project_key),
    INDEX idx_bot_id (bot_id)
);

-- 插入示例数据（可选）
INSERT INTO smart_elf (bot_id, bot_secret, project_key, reply_switch) 
VALUES ('test_bot_id', 'test_bot_secret', 'test_project_key', TRUE)
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;