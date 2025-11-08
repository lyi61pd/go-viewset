package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 应用配置
type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string `json:"type"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	Charset      string `json:"charset"`
	ParseTime    bool   `json:"parseTime"`
	Loc          string `json:"loc"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"`
	Mode string `json:"mode"`
}

// GetDSN 生成数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.Charset,
		d.ParseTime,
		d.Loc,
	)
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 如果没有指定路径，使用默认路径
	if configPath == "" {
		configPath = "config.json"
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}
