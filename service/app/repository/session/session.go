// Package session 提供登录会话（刷新令牌）的存储。
// 会话存放在 Redis 中，不落 MySQL；key 带 TTL，过期后自动删除。
package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// keyPrefix 是刷新令牌 Redis key 的固定前缀。
const keyPrefix = "xqyun:refresh"

var (
	// ErrNotFound 表示刷新令牌不存在（已过期、已吊销或从未签发）。
	ErrNotFound = errors.New("刷新令牌不存在")
	// ErrRefreshTokenRotated 表示刷新令牌已被使用，不能再参与轮换。
	ErrRefreshTokenRotated = errors.New("刷新令牌已失效")
	// ErrStoreNotInitialized 表示 Redis 客户端未初始化。
	ErrStoreNotInitialized = errors.New("会话存储未初始化")
)

// Token 是一条登录会话记录。
// 明文刷新令牌只通过 HttpOnly Cookie 传递，Redis 中保存的是 SHA-256 后的哈希。
type Token struct {
	UserID    uint64  `json:"user_id"`
	TokenHash string  `json:"-"`
	TokenType string  `json:"token_type"`
	UserAgent *string `json:"user_agent,omitempty"`
	IP        *string `json:"ip,omitempty"`
	ExpiresAt int64   `json:"expires_at"`
	CreatedAt int64   `json:"created_at"`
}

// Repository 定义登录会话的存储能力。
type Repository interface {
	// Save 写入一条会话，TTL 与 ExpiresAt 对齐。
	Save(ctx context.Context, token Token) error

	// FindByHash 按令牌哈希与客户端类型读取会话。
	FindByHash(ctx context.Context, tokenHash, tokenType string) (*Token, error)

	// Rotate 原子地删除旧会话并写入新会话。
	// 旧令牌不存在（已使用或已吊销）时返回 ErrRefreshTokenRotated。
	Rotate(ctx context.Context, currentHash, tokenType string, next Token) error

	// Revoke 删除会话；不存在时视为成功。
	Revoke(ctx context.Context, tokenHash, tokenType string) error

	// RevokeByUser 删除某个用户的全部会话，用于改密 / 停用账号后强制重新登录。
	RevokeByUser(ctx context.Context, userID uint64, tokenType string) error
}

// RedisRepository 是 Repository 的 Redis 实现。
type RedisRepository struct {
	client *redis.Client
}

var _ Repository = (*RedisRepository)(nil)

// rotateScript 用 Lua 保证「检查旧令牌 → 写入新令牌 → 删除旧令牌」的原子性，
// 避免同一个刷新令牌被并发使用两次。
// 先写新令牌再删旧令牌：即使写入失败，旧会话也仍然可用，不会丢登录状态。
// ARGV[2] 是相对 TTL（毫秒），兼容不支持 PXAT 的旧版 Redis。
// KEYS[3] 是该用户的会话索引 Set：轮换时同步移除旧哈希、加入新哈希，
// 保证 RevokeByUser 能覆盖轮换后的新令牌。
var rotateScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 0 then
  return 0
end
redis.call('SET', KEYS[2], ARGV[1], 'PX', ARGV[2])
redis.call('DEL', KEYS[1])
redis.call('SREM', KEYS[3], ARGV[3])
redis.call('SADD', KEYS[3], ARGV[4])
redis.call('PEXPIRE', KEYS[3], ARGV[2])
return 1
`)

// New 创建 Redis 会话存储；client 为 nil 时各方法返回 ErrStoreNotInitialized。
func New(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

// Save 写入会话，TTL 与 ExpiresAt 对齐。
func (r *RedisRepository) Save(ctx context.Context, token Token) error {
	if r == nil || r.client == nil {
		return ErrStoreNotInitialized
	}

	ttl := time.Until(time.Unix(token.ExpiresAt, 0))
	if ttl <= 0 {
		return fmt.Errorf("会话存储：刷新令牌过期时间无效")
	}

	payload, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("会话存储：编码会话失败：%w", err)
	}

	pipe := r.client.TxPipeline()
	pipe.Set(ctx, r.key(token.TokenType, token.TokenHash), payload, ttl)
	pipe.SAdd(ctx, r.userKey(token.TokenType, token.UserID), token.TokenHash)
	pipe.Expire(ctx, r.userKey(token.TokenType, token.UserID), ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("会话存储：保存会话失败：%w", err)
	}
	return nil
}

// FindByHash 按令牌哈希与客户端类型读取会话。
func (r *RedisRepository) FindByHash(
	ctx context.Context,
	tokenHash string,
	tokenType string,
) (*Token, error) {
	if r == nil || r.client == nil {
		return nil, ErrStoreNotInitialized
	}

	raw, err := r.client.Get(ctx, r.key(tokenType, tokenHash)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("会话存储：读取会话失败：%w", err)
	}

	var token Token
	if err := json.Unmarshal(raw, &token); err != nil {
		return nil, fmt.Errorf("会话存储：解析会话失败：%w", err)
	}
	// key 本身已经包含哈希与类型，补齐到返回值便于上层使用。
	token.TokenHash = tokenHash
	token.TokenType = tokenType

	if token.ExpiresAt <= time.Now().Unix() {
		return nil, ErrNotFound
	}
	return &token, nil
}

// Rotate 原子地删除旧会话并写入新会话。
func (r *RedisRepository) Rotate(
	ctx context.Context,
	currentHash string,
	tokenType string,
	next Token,
) error {
	if r == nil || r.client == nil {
		return ErrStoreNotInitialized
	}

	ttlMs := next.ExpiresAt*1000 - time.Now().UnixMilli()
	if ttlMs <= 0 {
		return fmt.Errorf("会话存储：新刷新令牌过期时间无效")
	}

	payload, err := json.Marshal(next)
	if err != nil {
		return fmt.Errorf("会话存储：编码会话失败：%w", err)
	}

	result, err := rotateScript.Run(
		ctx,
		r.client,
		[]string{
			r.key(tokenType, currentHash),
			r.key(tokenType, next.TokenHash),
			r.userKey(tokenType, next.UserID),
		},
		payload,
		ttlMs,
		currentHash,
		next.TokenHash,
	).Int()
	if err != nil {
		return fmt.Errorf("会话存储：轮换刷新令牌失败：%w", err)
	}
	if result == 0 {
		return ErrRefreshTokenRotated
	}
	return nil
}

// Revoke 删除会话，不存在时同样返回成功。
func (r *RedisRepository) Revoke(ctx context.Context, tokenHash, tokenType string) error {
	if r == nil || r.client == nil {
		return ErrStoreNotInitialized
	}

	if err := r.client.Del(ctx, r.key(tokenType, tokenHash)).Err(); err != nil {
		return fmt.Errorf("会话存储：删除会话失败：%w", err)
	}
	// 不清理用户索引：索引里残留的哈希只是指向已删除的 key，
	// RevokeByUser 删除不存在的 key 不会报错，索引本身有 TTL 兜底。
	return nil
}

// RevokeByUser 删除某个用户的全部会话。
// 会话索引是 Set，成员为刷新令牌哈希；索引不存在时直接删除索引 key。
func (r *RedisRepository) RevokeByUser(
	ctx context.Context,
	userID uint64,
	tokenType string,
) error {
	if r == nil || r.client == nil {
		return ErrStoreNotInitialized
	}

	setKey := r.userKey(tokenType, userID)
	hashes, err := r.client.SMembers(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("会话存储：读取用户会话失败：%w", err)
	}

	keys := make([]string, 0, len(hashes)+1)
	keys = append(keys, setKey)
	for _, hash := range hashes {
		keys = append(keys, r.key(tokenType, hash))
	}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("会话存储：吊销用户会话失败：%w", err)
	}
	return nil
}

// key 生成刷新令牌的 Redis key：xqyun:refresh:{token_type}:{token_hash}。
func (r *RedisRepository) key(tokenType, tokenHash string) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefix, strings.TrimSpace(tokenType), tokenHash)
}

// userKey 生成用户会话索引的 Redis key：xqyun:refresh:user:{token_type}:{user_id}。
// Set 成员是刷新令牌哈希，用于改密 / 停用账号时批量吊销。
func (r *RedisRepository) userKey(tokenType string, userID uint64) string {
	return fmt.Sprintf("%s:user:%s:%d", keyPrefix, strings.TrimSpace(tokenType), userID)
}
