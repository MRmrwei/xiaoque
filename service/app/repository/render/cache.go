// Package render 提供官网前台渲染数据（页面 / 导航 / 站点配置）的 Redis 缓存。
// 缓存按分片存储：页面按 locale + path，导航与站点配置按 locale。
package render

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"

	"suki/xqyun/app/dto"
)

// keyPrefix 与登录会话的 xqyun:refresh:* 区分开。
const keyPrefix = "xqyun:render"

// supportedLocales 是导航缓存覆盖的语言，导航数据中同时保存中英文文案。
var supportedLocales = []string{"zh-CN", "en"}

// PageKey 是页面缓存分片的键。
type PageKey struct {
	Locale string
	Path   string
}

// Cache 是前台渲染数据的分片缓存。
// client 为 nil 或 ttl <= 0 时所有方法都不读写 Redis，按未命中处理。
type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

// New 创建渲染缓存。
func New(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{client: client, ttl: ttl}
}

// Enabled 表示缓存是否生效。
func (c *Cache) Enabled() bool {
	return c != nil && c.client != nil && c.ttl > 0
}

// Page 读取页面分片。
func (c *Cache) Page(ctx context.Context, locale string, path string) (dto.PublishedPage, bool) {
	var page dto.PublishedPage
	if !c.get(ctx, pageKey(locale, path), &page) {
		return dto.PublishedPage{}, false
	}
	return page, true
}

// StorePage 写入页面分片。
func (c *Cache) StorePage(
	ctx context.Context,
	locale string,
	path string,
	page dto.PublishedPage,
) {
	c.set(ctx, pageKey(locale, path), page)
}

// InvalidatePage 删除页面分片。
func (c *Cache) InvalidatePage(ctx context.Context, keys ...PageKey) error {
	if len(keys) == 0 {
		return nil
	}
	values := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		value := pageKey(key.Locale, key.Path)
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	return c.delete(ctx, values...)
}

// Nav 读取导航分片。
func (c *Cache) Nav(ctx context.Context, locale string) (dto.SiteNav, bool) {
	var nav dto.SiteNav
	if !c.get(ctx, navKey(locale), &nav) {
		return dto.SiteNav{}, false
	}
	return nav, true
}

// StoreNav 写入导航分片。
func (c *Cache) StoreNav(ctx context.Context, locale string, nav dto.SiteNav) {
	c.set(ctx, navKey(locale), nav)
}

// InvalidateNav 删除所有语言的导航分片，菜单保存后调用。
func (c *Cache) InvalidateNav(ctx context.Context) error {
	keys := make([]string, 0, len(supportedLocales))
	for _, locale := range supportedLocales {
		keys = append(keys, navKey(locale))
	}
	return c.delete(ctx, keys...)
}

// Site 读取站点配置分片。
func (c *Cache) Site(ctx context.Context, locale string) (map[string]any, bool) {
	var site map[string]any
	if !c.get(ctx, siteKey(locale), &site) {
		return nil, false
	}
	return site, true
}

// StoreSite 写入站点配置分片。
func (c *Cache) StoreSite(ctx context.Context, locale string, site map[string]any) {
	c.set(ctx, siteKey(locale), site)
}

// InvalidateSite 删除指定语言的站点配置分片。
func (c *Cache) InvalidateSite(ctx context.Context, locales ...string) error {
	if len(locales) == 0 {
		return nil
	}
	keys := make([]string, 0, len(locales))
	seen := make(map[string]struct{}, len(locales))
	for _, locale := range locales {
		value := siteKey(locale)
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		keys = append(keys, value)
	}
	return c.delete(ctx, keys...)
}

// get 读取并解析缓存；未命中、Redis 异常或内容损坏都按未命中处理。
func (c *Cache) get(ctx context.Context, key string, target any) bool {
	if !c.Enabled() {
		return false
	}

	raw, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false
	}
	if err != nil {
		// 缓存故障不影响主流程，回源数据库即可。
		hlog.CtxWarnf(ctx, "render cache: get %s: %v", key, err)
		return false
	}
	if err := json.Unmarshal(raw, target); err != nil {
		hlog.CtxWarnf(ctx, "render cache: decode %s: %v", key, err)
		_ = c.client.Del(ctx, key).Err()
		return false
	}
	return true
}

// set 写入缓存；序列化或 Redis 异常只记日志，不影响请求。
func (c *Cache) set(ctx context.Context, key string, value any) {
	if !c.Enabled() {
		return
	}

	raw, err := json.Marshal(value)
	if err != nil {
		hlog.CtxWarnf(ctx, "render cache: encode %s: %v", key, err)
		return
	}
	if err := c.client.Set(ctx, key, raw, c.ttl).Err(); err != nil {
		hlog.CtxWarnf(ctx, "render cache: set %s: %v", key, err)
	}
}

// delete 删除缓存，返回错误交给调用方记日志。
func (c *Cache) delete(ctx context.Context, keys ...string) error {
	if !c.Enabled() || len(keys) == 0 {
		return nil
	}
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("删除渲染缓存失败：%w", err)
	}
	return nil
}

func pageKey(locale string, path string) string {
	return keyPrefix + ":page:" + locale + ":" + path
}

func navKey(locale string) string {
	return keyPrefix + ":nav:" + locale
}

func siteKey(locale string) string {
	return keyPrefix + ":site:" + locale
}
