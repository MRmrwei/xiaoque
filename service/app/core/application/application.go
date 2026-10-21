package application

import (
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 聚合应用启动所需的基础设施配置。
type Config struct {
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

// JWTConfig 描述 JWT 签发方与各客户端令牌参数。
type JWTConfig struct {
	Issuer  string
	Clients map[string]JWTClientConfig
	Cookie  CookieConfig
}

// CookieConfig 描述 Refresh Token 的 HttpOnly Cookie。
type CookieConfig struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	SameSite string
}

// JWTClientConfig 描述单个客户端的令牌签名与有效期参数。
type JWTClientConfig struct {
	Secret     []byte
	Audience   string
	TokenType  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// Client 返回指定客户端的 JWT 配置。
func (c JWTConfig) Client(name string) (JWTClientConfig, bool) {
	client, ok := c.Clients[name]
	return client, ok
}

// DBConfig 描述 MySQL 连接参数。
type DBConfig struct {
	DSN string
	// MaxOpenConns 最大打开连接数，<=0 时使用驱动默认值。
	MaxOpenConns int
	// MaxIdleConns 最大空闲连接数，<=0 时使用驱动默认值。
	MaxIdleConns int
	// ConnMaxLifetime 连接最长存活时间，<=0 时不限制。
	ConnMaxLifetime time.Duration
	// LogLevel GORM 日志级别，默认 logger.Warn。
	LogLevel logger.LogLevel
}

// RedisConfig 描述 Redis 连接参数。
type RedisConfig struct {
	// Enabled 为 false 时不连接 Redis，Application.Redis() 返回 nil。
	Enabled bool
	Addr    string
	// Username / Password 按 Redis 实例配置填写，未设置 ACL 时留空。
	Username string
	Password string
	DB       int
	PoolSize int
	// DialTimeout / ReadTimeout / WriteTimeout 为 0 时使用 go-redis 默认值。
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Application 定义应用级基础设施的统一访问入口。
// 业务层只依赖这个接口，便于替换实现或在测试中使用 mock。
type Application interface {
	// DB 返回 MySQL 连接（GORM）。
	DB() *gorm.DB

	// Redis 返回 Redis 客户端；未启用 Redis 时返回 nil。
	Redis() *redis.Client

	// JWT 返回管理后台令牌签名配置。
	JWT() JWTConfig

	// Close 释放数据库与 Redis 连接。
	Close() error
}
