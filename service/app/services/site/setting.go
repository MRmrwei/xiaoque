package site

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	renderrepository "suki/xqyun/app/repository/render"
	settingrepository "suki/xqyun/app/repository/setting"
)

// globalSettingKeys 是不区分语言的站点设置项，按空 locale 存储。
// 管理后台的站点配置一次编辑中英文文案，值本身是多语言对象，因此这些 key 都按全局保存。
var globalSettingKeys = map[string]struct{}{
	"site_name":      {},
	"company_name":   {},
	"tagline":        {},
	"contact":        {},
	"brand":          {},
	"header":         {},
	"footer":         {},
	"copyright":      {},
	"icp":            {},
	"socials":        {},
	"default_locale": {},
	"allow_index":    {},
	"seo":            {},
}

// SettingService 是站点配置用例服务。
type SettingService struct {
	settings *settingrepository.Repository
	cache    *renderrepository.Cache
}

// NewSettingService 创建站点配置用例服务。
// cache 只在前台渲染缓存启用时生效，保存设置后按语言失效站点配置分片。
func NewSettingService(
	settings *settingrepository.Repository,
	cache *renderrepository.Cache,
) *SettingService {
	return &SettingService{settings: settings, cache: cache}
}

// Get 读取指定语言的站点设置，全局项与语言级设置合并返回。
func (s *SettingService) Get(ctx context.Context, locale string) (map[string]any, error) {
	normalized, err := normalizeLocale(locale, true)
	if err != nil {
		return nil, err
	}

	rows, err := s.settings.List(ctx, normalized)
	if err != nil {
		return nil, err
	}

	result := map[string]any{"locale": normalized}
	for _, row := range rows {
		var value any
		if err := json.Unmarshal(row.Value, &value); err != nil {
			return nil, fmt.Errorf("解析站点设置 %s 失败：%w", row.Key, err)
		}
		result[row.Key] = value
	}
	return result, nil
}

// Save 按 key + locale 保存站点设置；品牌、收录等全局项不区分语言。
func (s *SettingService) Save(
	ctx context.Context,
	locale string,
	values map[string]json.RawMessage,
	operatorID *uint64,
	now int64,
) (map[string]any, error) {
	normalized, err := normalizeLocale(locale, true)
	if err != nil {
		return nil, err
	}

	global := make(map[string]json.RawMessage)
	localized := make(map[string]json.RawMessage)
	for key, value := range values {
		settingKey := strings.TrimSpace(key)
		if settingKey == "" || settingKey == "locale" {
			continue
		}
		if len(settingKey) > 100 {
			return nil, badRequest("设置项名称不能超过 100 个字符")
		}
		if len(value) == 0 {
			continue
		}

		if _, ok := globalSettingKeys[settingKey]; ok {
			global[settingKey] = value
			continue
		}
		localized[settingKey] = value
	}

	if len(global) == 0 && len(localized) == 0 {
		return nil, badRequest("请提供站点设置内容")
	}

	if err := s.settings.UpsertMany(ctx, "", global, operatorID, now); err != nil {
		return nil, err
	}
	if err := s.settings.UpsertMany(ctx, normalized, localized, operatorID, now); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, normalized, len(global) > 0, len(localized) > 0)
	return s.Get(ctx, normalized)
}

// invalidateCache 失效站点配置分片：全局项影响两种语言，语言级项只影响对应语言。
func (s *SettingService) invalidateCache(
	ctx context.Context,
	locale string,
	hasGlobal bool,
	hasLocalized bool,
) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}

	locales := make([]string, 0, 2)
	if hasGlobal {
		locales = append(locales, localeZH, localeEN)
	}
	if hasLocalized {
		locales = append(locales, locale)
	}
	if err := s.cache.InvalidateSite(ctx, locales...); err != nil {
		hlog.CtxWarnf(ctx, "render cache: invalidate site: %v", err)
	}
}
