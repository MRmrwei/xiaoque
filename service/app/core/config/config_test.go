package config

import (
	"strings"
	"testing"
	"time"
)

const testJWTSecret = "0123456789abcdef0123456789abcdef"

func TestJWTClientDefaults(t *testing.T) {
	cfg := &Config{JWT: JWTConfig{Secret: testJWTSecret}}
	cfg.setDefaults()

	admin := cfg.JWT.Clients["admin"]
	if admin.Audience != "admin-api" || admin.TokenType != "admin" {
		t.Fatalf("admin 受众 / 类型不正确：%#v", admin)
	}
	if admin.AccessTTL != "15m" || admin.RefreshTTL != "168h" {
		t.Fatalf("admin TTL 不正确：%#v", admin)
	}

	frontend := cfg.JWT.Clients["frontend"]
	if frontend.Audience != "frontend-api" || frontend.TokenType != "frontend" {
		t.Fatalf("frontend 受众 / 类型不正确：%#v", frontend)
	}
	if frontend.AccessTTL != "30m" || frontend.RefreshTTL != "720h" {
		t.Fatalf("frontend TTL 不正确：%#v", frontend)
	}
}

func TestApplicationConfigBuildsJWTClient(t *testing.T) {
	cfg := &Config{
		Redis: RedisConfig{Enabled: true},
		JWT:   JWTConfig{Secret: testJWTSecret, Issuer: "xiaoque-admin"},
	}
	cfg.setDefaults()

	appCfg, err := cfg.ApplicationConfig()
	if err != nil {
		t.Fatalf("ApplicationConfig() 返回错误：%v", err)
	}

	admin, ok := appCfg.JWT.Client("admin")
	if !ok {
		t.Fatal("admin 客户端未生成")
	}
	if admin.Audience != "admin-api" || admin.TokenType != "admin" {
		t.Fatalf("admin 受众 / 类型不正确：%#v", admin)
	}
	if admin.AccessTTL != 15*time.Minute || admin.RefreshTTL != 168*time.Hour {
		t.Fatalf("admin TTL 不正确：%#v", admin)
	}
	if string(admin.Secret) != testJWTSecret {
		t.Fatal("admin 密钥未正确传递")
	}

	frontend, ok := appCfg.JWT.Client("frontend")
	if !ok {
		t.Fatal("frontend 客户端未生成")
	}
	if frontend.AccessTTL != 30*time.Minute || frontend.RefreshTTL != 720*time.Hour {
		t.Fatalf("frontend TTL 不正确：%#v", frontend)
	}
}

func TestApplicationConfigRejectsInvalidJWT(t *testing.T) {
	t.Run("missing admin client", func(t *testing.T) {
		cfg := &Config{Redis: RedisConfig{Enabled: true}, JWT: JWTConfig{
			Secret: testJWTSecret,
			Issuer: "xiaoque-admin",
			Clients: map[string]JWTClientConfig{
				"frontend": {
					Audience:   "frontend-api",
					TokenType:  "frontend",
					AccessTTL:  "30m",
					RefreshTTL: "720h",
				},
			},
		}}
		if _, err := cfg.ApplicationConfig(); err == nil {
			t.Fatal("缺少 admin 客户端时应返回错误")
		}
	})

	t.Run("short secret", func(t *testing.T) {
		cfg := &Config{
			Redis: RedisConfig{Enabled: true},
			JWT:   JWTConfig{Secret: strings.Repeat("s", 31)},
		}
		cfg.setDefaults()
		if _, err := cfg.ApplicationConfig(); err == nil {
			t.Fatal("密钥长度不足时应返回错误")
		}
	})

	t.Run("invalid ttl", func(t *testing.T) {
		cfg := &Config{Redis: RedisConfig{Enabled: true}, JWT: JWTConfig{
			Secret: testJWTSecret,
			Issuer: "xiaoque-admin",
			Clients: map[string]JWTClientConfig{
				"admin": {
					Audience:   "admin-api",
					TokenType:  "admin",
					AccessTTL:  "invalid",
					RefreshTTL: "168h",
				},
			},
		}}
		if _, err := cfg.ApplicationConfig(); err == nil {
			t.Fatal("非法 TTL 应返回错误")
		}
	})
}

func TestApplicationConfigRequiresRedis(t *testing.T) {
	cfg := &Config{JWT: JWTConfig{Secret: testJWTSecret}}
	cfg.setDefaults()

	if _, err := cfg.ApplicationConfig(); err == nil {
		t.Fatal("Redis 未启用时应返回错误")
	}
}

func TestJWTClientEnvOverride(t *testing.T) {
	t.Setenv("JWT_ADMIN_ACCESS_TTL", "5m")
	t.Setenv("JWT_FRONTEND_REFRESH_TTL", "168h")

	cfg := &Config{JWT: JWTConfig{Secret: testJWTSecret}}
	cfg.setDefaults()
	if err := cfg.applyEnv(); err != nil {
		t.Fatalf("applyEnv() 返回错误：%v", err)
	}

	if got := cfg.JWT.Clients["admin"].AccessTTL; got != "5m" {
		t.Fatalf("admin access_ttl = %q，期望 5m", got)
	}
	if got := cfg.JWT.Clients["frontend"].RefreshTTL; got != "168h" {
		t.Fatalf("frontend refresh_ttl = %q，期望 168h", got)
	}
}
