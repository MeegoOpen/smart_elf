package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Feishu   FeishuConfig   `yaml:"feishu"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type FeishuConfig struct {
	// 飞书开放平台 Open API Host
	IMOpenAPIHost string `yaml:"im_open_api_host"`
	// Meego 项目相关配置
	PluginID       string `yaml:"plugin_id"`
	PluginSecret   string `yaml:"plugin_secret"`
	ProjectAPIHost string `yaml:"project_api_host"`
	ProjectWebHost string `yaml:"project_web_host"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Debug        bool   `yaml:"debug"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 从 conf/config.yaml 加载配置
	cfgPath := filepath.Join("conf", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	log.Printf("配置加载完成: port=%d, host=%s, log_level=%s, db_debug=%v",
		config.Server.Port, config.Server.Host, config.Logger.Level, config.Database.Debug)

	return &config, nil
}

// 全局单例配置
var globalConfig *Config
var globalOnce sync.Once

// LoadGlobalConfig 以单例方式加载配置（仅首次读取文件）
func LoadGlobalConfig() (*Config, error) {
	var loadErr error
	globalOnce.Do(func() {
		var cfg *Config
		cfg, loadErr = LoadConfig()
		if loadErr == nil {
			globalConfig = cfg
		}
	})
	if globalConfig == nil {
		return nil, loadErr
	}
	return globalConfig, nil
}

// SetGlobalConfig 设置全局配置实例（已加载的配置）
func SetGlobalConfig(cfg *Config) {
	globalConfig = cfg
}

// GetConfig 获取全局配置实例（未初始化时返回nil）
func GetConfig() *Config {
	return globalConfig
}
