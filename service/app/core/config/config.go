package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm/logger"

	"suki/xqyun/app/core/application"
)

// Config 对应 config.yaml；密码字段不参与 YAML 解析，只从环境变量读取。
type Config struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	JWT    JWTConfig    `yaml:"jwt"`
	Log    LogConfig    `yaml:"log"`
}

type AppConfig struct {
	Env string `yaml:"env"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	// Username / Password 只从 MYSQL_USERNAME / MYSQL_PASSWORD 读取，避免写进 config.yaml。
	Username        string `yaml:"-"`
	Password        string `yaml:"-"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parse_time"`
	Loc             string `yaml:"loc"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	LogLevel        int    `yaml:"log_level"`

	// dsnOverride 由 MYSQL_DSN 设置，优先级高于分项配置。
	dsnOverride string
}

type RedisConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	// Addr 可选，设置后优先于 Host / Port。
	Addr string `yaml:"addr"`
	// Password 只从 REDIS_PASSWORD 读取。
	Password string `yaml:"-"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// JWTConfig 描述 JWT 签发方与各客户端参数。
// 现阶段使用同一个密钥，客户端通过 audience / token_type 与 TTL 区分。
type JWTConfig struct {
	// Secret 只从 JWT_SECRET 读取，至少 32 字节。
	Secret  string                     `yaml:"-"`
	Issuer  string                     `yaml:"issuer"`
	Clients map[string]JWTClientConfig `yaml:"clients"`
	Cookie  JWTCookieConfig            `yaml:"cookie"`
}

// JWTCookieConfig 描述 Refresh Token 的 HttpOnly Cookie。
type JWTCookieConfig struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	SameSite string `yaml:"same_site"`
}

// JWTClientConfig 描述单个客户端（admin / frontend 等）的令牌参数。
type JWTClientConfig struct {
	Audience   string `yaml:"audience"`
	TokenType  string `yaml:"token_type"`
	AccessTTL  string `yaml:"access_ttl"`
	RefreshTTL string `yaml:"refresh_ttl"`
}

// LogConfig 描述日志级别与文件轮转参数。
// 开发环境输出到控制台；生产环境输出到 Dir/Filename，并按大小 / 保留数量轮转。
type LogConfig struct {
	Level      string `yaml:"level"`
	Dir        string `yaml:"dir"`
	Filename   string `yaml:"filename"`
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
	Compress   bool   `yaml:"compress"`
	// Color 只影响开发环境终端输出；生产文件日志始终不带颜色。
	Color bool `yaml:"color"`
}

// Load 读取 YAML 配置，并叠加 .env 与环境变量。
// 优先级：环境变量 > .env > config.yaml > 代码默认值。
func Load(path string) (*Config, error) {
	// 找不到 .env 时忽略，部署环境通常直接注入环境变量。
	_ = godotenv.Load()

	if path == "" {
		path = os.Getenv("CONFIG_FILE")
	}
	if path == "" {
		path = "config.yaml"
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}

	cfg.setDefaults()
	if err := cfg.applyEnv(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) setDefaults() {
	if c.App.Env == "" {
		c.App.Env = "dev"
	}
	if c.Server.Addr == "" {
		c.Server.Addr = ":8888"
	}
	if c.MySQL.Port == 0 {
		c.MySQL.Port = 3306
	}
	if c.MySQL.Charset == "" {
		c.MySQL.Charset = "utf8mb4"
	}
	if c.MySQL.Loc == "" {
		c.MySQL.Loc = "Local"
	}
	if c.Redis.Port == 0 {
		c.Redis.Port = 6379
	}
	if c.JWT.Issuer == "" {
		c.JWT.Issuer = "xiaoque-admin"
	}
	if c.JWT.Cookie.Name == "" {
		c.JWT.Cookie.Name = "xqyun_admin_refresh"
	}
	if c.JWT.Cookie.Path == "" {
		c.JWT.Cookie.Path = "/api/v1/admin"
	}
	if c.JWT.Cookie.SameSite == "" {
		c.JWT.Cookie.SameSite = "Lax"
	}
	if c.JWT.Clients == nil {
		c.JWT.Clients = make(map[string]JWTClientConfig)
	}

	admin := c.JWT.Clients["admin"]
	if admin.Audience == "" {
		admin.Audience = "admin-api"
	}
	if admin.TokenType == "" {
		admin.TokenType = "admin"
	}
	if admin.AccessTTL == "" {
		admin.AccessTTL = "15m"
	}
	if admin.RefreshTTL == "" {
		admin.RefreshTTL = "168h"
	}
	c.JWT.Clients["admin"] = admin

	frontend := c.JWT.Clients["frontend"]
	if frontend.Audience == "" {
		frontend.Audience = "frontend-api"
	}
	if frontend.TokenType == "" {
		frontend.TokenType = "frontend"
	}
	if frontend.AccessTTL == "" {
		frontend.AccessTTL = "30m"
	}
	if frontend.RefreshTTL == "" {
		frontend.RefreshTTL = "720h"
	}
	c.JWT.Clients["frontend"] = frontend
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.Dir == "" {
		c.Log.Dir = "logs"
	}
	if c.Log.Filename == "" {
		c.Log.Filename = "app.log"
	}
	if c.Log.MaxSizeMB <= 0 {
		c.Log.MaxSizeMB = 100
	}
	if c.Log.MaxBackups <= 0 {
		c.Log.MaxBackups = 30
	}
	if c.Log.MaxAgeDays <= 0 {
		c.Log.MaxAgeDays = 30
	}
}

func (c *Config) applyEnv() error {
	setString(&c.App.Env, "APP_ENV")
	setString(&c.Server.Addr, "SERVER_ADDR")

	setString(&c.MySQL.Host, "MYSQL_HOST")
	if err := setInt(&c.MySQL.Port, "MYSQL_PORT"); err != nil {
		return err
	}
	setString(&c.MySQL.Database, "MYSQL_DATABASE")
	setString(&c.MySQL.Username, "MYSQL_USERNAME")
	setString(&c.MySQL.Password, "MYSQL_PASSWORD")
	setString(&c.MySQL.dsnOverride, "MYSQL_DSN")

	if err := setBool(&c.Redis.Enabled, "REDIS_ENABLED"); err != nil {
		return err
	}
	setString(&c.Redis.Host, "REDIS_HOST")
	if err := setInt(&c.Redis.Port, "REDIS_PORT"); err != nil {
		return err
	}
	setString(&c.Redis.Addr, "REDIS_ADDR")
	setString(&c.Redis.Password, "REDIS_PASSWORD")
	if err := setInt(&c.Redis.DB, "REDIS_DB"); err != nil {
		return err
	}
	if err := setInt(&c.Redis.PoolSize, "REDIS_POOL_SIZE"); err != nil {
		return err
	}

	setString(&c.JWT.Secret, "JWT_SECRET")
	setString(&c.JWT.Issuer, "JWT_ISSUER")
	setString(&c.JWT.Cookie.Name, "JWT_COOKIE_NAME")
	setString(&c.JWT.Cookie.Path, "JWT_COOKIE_PATH")
	setString(&c.JWT.Cookie.Domain, "JWT_COOKIE_DOMAIN")
	if err := setBool(&c.JWT.Cookie.Secure, "JWT_COOKIE_SECURE"); err != nil {
		return err
	}
	setString(&c.JWT.Cookie.SameSite, "JWT_COOKIE_SAME_SITE")
	for name, client := range c.JWT.Clients {
		prefix := "JWT_" + strings.ToUpper(name)
		setString(&client.Audience, prefix+"_AUDIENCE")
		setString(&client.TokenType, prefix+"_TOKEN_TYPE")
		setString(&client.AccessTTL, prefix+"_ACCESS_TTL")
		setString(&client.RefreshTTL, prefix+"_REFRESH_TTL")
		c.JWT.Clients[name] = client
	}

	setString(&c.Log.Level, "LOG_LEVEL")
	setString(&c.Log.Dir, "LOG_DIR")
	setString(&c.Log.Filename, "LOG_FILENAME")
	if err := setInt(&c.Log.MaxSizeMB, "LOG_MAX_SIZE_MB"); err != nil {
		return err
	}
	if err := setInt(&c.Log.MaxBackups, "LOG_MAX_BACKUPS"); err != nil {
		return err
	}
	if err := setInt(&c.Log.MaxAgeDays, "LOG_MAX_AGE_DAYS"); err != nil {
		return err
	}
	if err := setBool(&c.Log.Compress, "LOG_COMPRESS"); err != nil {
		return err
	}
	if err := setBool(&c.Log.Color, "LOG_COLOR"); err != nil {
		return err
	}

	return nil
}

// ApplicationConfig 将配置转换为 application 层使用的结构。
func (c *Config) ApplicationConfig() (application.Config, error) {
	lifetime, err := parseDuration(c.MySQL.ConnMaxLifetime)
	if err != nil {
		return application.Config{}, fmt.Errorf("config: mysql.conn_max_lifetime: %w", err)
	}

	// 刷新令牌存放在 Redis，关闭 Redis 会导致登录 / 刷新接口不可用，启动时直接报错。
	if !c.Redis.Enabled {
		return application.Config{}, fmt.Errorf("config: 登录会话存储在 Redis 中，请设置 redis.enabled=true")
	}

	if len(c.JWT.Secret) < 32 {
		return application.Config{}, fmt.Errorf("config: JWT_SECRET 至少需要 32 字节，请在 service/.env 中配置")
	}

	clients := make(map[string]application.JWTClientConfig, len(c.JWT.Clients))
	for name, client := range c.JWT.Clients {
		if client.Audience == "" || client.TokenType == "" {
			return application.Config{}, fmt.Errorf("config: jwt.clients.%s audience / token_type 不能为空", name)
		}
		accessTTL, err := parseDuration(client.AccessTTL)
		if err != nil {
			return application.Config{}, fmt.Errorf("config: jwt.clients.%s.access_ttl: %w", name, err)
		}
		if accessTTL <= 0 {
			return application.Config{}, fmt.Errorf("config: jwt.clients.%s.access_ttl 必须大于 0", name)
		}
		refreshTTL, err := parseDuration(client.RefreshTTL)
		if err != nil {
			return application.Config{}, fmt.Errorf("config: jwt.clients.%s.refresh_ttl: %w", name, err)
		}
		if refreshTTL <= 0 {
			return application.Config{}, fmt.Errorf("config: jwt.clients.%s.refresh_ttl 必须大于 0", name)
		}
		clients[name] = application.JWTClientConfig{
			Secret:     []byte(c.JWT.Secret),
			Audience:   client.Audience,
			TokenType:  client.TokenType,
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
		}
	}
	if _, ok := clients["admin"]; !ok {
		return application.Config{}, fmt.Errorf("config: jwt.clients.admin 未配置")
	}

	dsn := c.MySQL.dsnOverride
	if dsn == "" {
		dsn = c.MySQL.dsn()
	}

	return application.Config{
		DB: application.DBConfig{
			DSN:             dsn,
			MaxOpenConns:    c.MySQL.MaxOpenConns,
			MaxIdleConns:    c.MySQL.MaxIdleConns,
			ConnMaxLifetime: lifetime,
			LogLevel:        logger.LogLevel(c.MySQL.LogLevel),
		},
		Redis: application.RedisConfig{
			Enabled:  c.Redis.Enabled,
			Addr:     c.Redis.address(),
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
			PoolSize: c.Redis.PoolSize,
		},
		JWT: application.JWTConfig{
			Issuer:  c.JWT.Issuer,
			Clients: clients,
			Cookie: application.CookieConfig{
				Name:     c.JWT.Cookie.Name,
				Path:     c.JWT.Cookie.Path,
				Domain:   c.JWT.Cookie.Domain,
				Secure:   c.JWT.Cookie.Secure,
				SameSite: c.JWT.Cookie.SameSite,
			},
		},
	}, nil
}

func (m MySQLConfig) dsn() string {
	cfg := mysql.NewConfig()
	cfg.User = m.Username
	cfg.Passwd = m.Password
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(m.Host, strconv.Itoa(m.Port))
	cfg.DBName = m.Database
	cfg.ParseTime = m.ParseTime
	cfg.Params = map[string]string{"charset": m.Charset}

	if loc, err := time.LoadLocation(m.Loc); err == nil {
		cfg.Loc = loc
	} else {
		cfg.Loc = time.Local
	}

	return cfg.FormatDSN()
}

func (r RedisConfig) address() string {
	if r.Addr != "" {
		return r.Addr
	}
	return net.JoinHostPort(r.Host, strconv.Itoa(r.Port))
}

func parseDuration(value string) (time.Duration, error) {
	if value == "" {
		return 0, nil
	}
	return time.ParseDuration(value)
}

func setString(dst *string, key string) {
	if value, ok := os.LookupEnv(key); ok {
		*dst = value
	}
}

func setInt(dst *int, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	*dst = n
	return nil
}

func setBool(dst *bool, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("config: %s must be a boolean: %w", key, err)
	}
	*dst = b
	return nil
}
