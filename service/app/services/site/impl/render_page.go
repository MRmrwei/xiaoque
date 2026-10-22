// Package impl 提供 site.PageRenderer 的具体实现。
package impl

import (
	"context"
	"errors"
	"strconv"

	"suki/xqyun/app/dto"
	"suki/xqyun/app/models"
	pagerepository "suki/xqyun/app/repository/page"
	renderrepository "suki/xqyun/app/repository/render"
	settingrepository "suki/xqyun/app/repository/setting"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/site"
)

// pageStatusPublished 是页面已发布的取值，与 xqyun_pages.status 一致。
const pageStatusPublished = "published"

// Renderer 是 site.PageRenderer 的实现，聚合官网前台渲染所需的页面、导航、站点配置与内容集合。
// 只读数据：页面取 published_content，导航与站点配置取后台已保存的设置。
type Renderer struct {
	pages    *pagerepository.Repository
	menus    *site.MenuService
	settings *settingrepository.Repository
	users    *userrepository.Repository
	cache    *renderrepository.Cache
}

var _ site.PageRenderer = (*Renderer)(nil)

// NewRenderer 创建前台页面渲染实现。
// cache 只在前台渲染缓存启用时生效：页面、导航、站点配置分别按分片读取，
// 未命中时查库回填；发布 / 保存等写操作在 service 层主动失效。
func NewRenderer(
	pages *pagerepository.Repository,
	menus *site.MenuService,
	settings *settingrepository.Repository,
	users *userrepository.Repository,
	cache *renderrepository.Cache,
) *Renderer {
	return &Renderer{
		pages:    pages,
		menus:    menus,
		settings: settings,
		users:    users,
		cache:    cache,
	}
}

// RenderPage 按路径聚合已发布页面、导航、站点配置与数据源区块集合。
// 页面、导航、站点配置分别读取缓存，未命中时回源数据库并回填。
func (s *Renderer) RenderPage(
	ctx context.Context,
	input site.RenderPageInput,
) (dto.RenderPageData, error) {
	locale, err := site.NormalizeLocale(input.Locale, true)
	if err != nil {
		return dto.RenderPageData{}, err
	}

	path := site.NormalizePath(input.Path)
	if path == "" {
		path = "/"
	}

	page, ok := s.cache.Page(ctx, locale, path)
	if !ok {
		page, err = s.loadPage(ctx, locale, path)
		if err != nil {
			return dto.RenderPageData{}, err
		}
		s.cache.StorePage(ctx, locale, path, page)
	}

	nav, ok := s.cache.Nav(ctx, locale)
	if !ok {
		nav, err = s.menus.Nav(ctx, locale)
		if err != nil {
			return dto.RenderPageData{}, err
		}
		s.cache.StoreNav(ctx, locale, nav)
	}

	settings, ok := s.cache.Site(ctx, locale)
	if !ok {
		settings, err = s.siteConfig(ctx, locale)
		if err != nil {
			return dto.RenderPageData{}, err
		}
		s.cache.StoreSite(ctx, locale, settings)
	}

	return dto.RenderPageData{
		Page:        page,
		Nav:         nav,
		Site:        settings,
		Collections: emptyCollections(),
	}, nil
}

// loadPage 从数据库读取已发布页面并组装前台结构；不存在或未发布返回 ErrPageNotFound。
func (s *Renderer) loadPage(
	ctx context.Context,
	locale string,
	path string,
) (dto.PublishedPage, error) {
	page, err := s.pages.FindByPath(ctx, locale, path)
	if errors.Is(err, pagerepository.ErrNotFound) {
		return dto.PublishedPage{}, site.ErrPageNotFound
	}
	if err != nil {
		return dto.PublishedPage{}, err
	}
	if page.Status != pageStatusPublished {
		return dto.PublishedPage{}, site.ErrPageNotFound
	}

	seo, blocks, err := site.ParsePageContent(page.PublishedContent)
	if err != nil {
		return dto.PublishedPage{}, err
	}
	alternates, err := s.alternates(ctx, page)
	if err != nil {
		return dto.PublishedPage{}, err
	}
	publishedBy, err := s.publishedBy(ctx, page)
	if err != nil {
		return dto.PublishedPage{}, err
	}

	return dto.PublishedPage{
		ID:          formatID(page.ID),
		GroupID:     page.TranslationGroupID,
		Locale:      page.Locale,
		Path:        page.Path,
		Name:        page.Name,
		SEO:         seo,
		Blocks:      blocks,
		Version:     int(page.Version),
		PublishedAt: page.UpdatedAt,
		PublishedBy: publishedBy,
		Alternates:  alternates,
	}, nil
}

// alternates 计算语言配对地址，只包含已发布语言，供前台输出 hreflang 与语言切换。
func (s *Renderer) alternates(
	ctx context.Context,
	page *models.Page,
) (map[string]string, error) {
	result := map[string]string{page.Locale: page.Path}
	if page.TranslationGroupID == "" {
		return result, nil
	}

	siblings, err := s.pages.FindByGroupIDs(ctx, []string{page.TranslationGroupID})
	if err != nil {
		return nil, err
	}
	for _, sibling := range siblings {
		if sibling.ID == page.ID || sibling.Status != pageStatusPublished {
			continue
		}
		result[sibling.Locale] = sibling.Path
	}
	return result, nil
}

// publishedBy 返回最后发布人名称；没有操作人时返回空字符串。
func (s *Renderer) publishedBy(
	ctx context.Context,
	page *models.Page,
) (string, error) {
	if page.UpdatedBy == nil {
		return "", nil
	}
	names, err := s.users.FindNamesByIDs(ctx, []uint64{*page.UpdatedBy})
	if err != nil {
		return "", err
	}
	return names[*page.UpdatedBy], nil
}

// formatID 把数据库主键转换为前台使用的字符串 ID。
func formatID(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// emptyCollections 返回数据源区块的集合。
// 内容库（新闻 / 产品 / 案例）尚未上线，当前固定为空数组，前台按空列表降级。
func emptyCollections() map[string][]any {
	return map[string][]any{
		"news":     {},
		"products": {},
		"cases":    {},
	}
}
