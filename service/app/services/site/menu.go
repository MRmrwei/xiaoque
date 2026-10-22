package site

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/dto"
	"suki/xqyun/app/models"
	menurepository "suki/xqyun/app/repository/menu"
	renderrepository "suki/xqyun/app/repository/render"
)

// menuNames 是站点编排支持的菜单容器编码及默认名称。
var menuNames = map[string]string{
	"header":        "顶部导航",
	"footer":        "底部导航",
	"footer_bottom": "页脚底部链接",
}

// 菜单项支持的类型与打开方式。
var (
	linkTypes = map[string]struct{}{
		"page":     {},
		"content":  {},
		"external": {},
		"none":     {},
	}
	linkTargets = map[string]struct{}{
		"_self":  {},
		"_blank": {},
	}
)

// MenuService 是顶部 / 底部导航用例服务。
type MenuService struct {
	menus *menurepository.Repository
	cache *renderrepository.Cache
}

// NewMenuService 创建导航菜单用例服务。
// cache 只在前台渲染缓存启用时生效，保存菜单后失效导航分片。
func NewMenuService(
	menus *menurepository.Repository,
	cache *renderrepository.Cache,
) *MenuService {
	return &MenuService{menus: menus, cache: cache}
}

// List 查询菜单容器列表。
func (s *MenuService) List(ctx context.Context) (dto.AdminMenuListResponse, error) {
	menus, err := s.menus.List(ctx)
	if err != nil {
		return dto.AdminMenuListResponse{}, err
	}

	items := make([]dto.AdminMenuSummary, 0, len(menus))
	for index := range menus {
		items = append(items, toMenuSummary(&menus[index]))
	}
	return dto.AdminMenuListResponse{Items: items}, nil
}

// Detail 查询菜单详情与菜单项。
func (s *MenuService) Detail(ctx context.Context, code string) (dto.AdminMenuDetail, error) {
	normalized, err := normalizeMenuCode(code)
	if err != nil {
		return dto.AdminMenuDetail{}, err
	}

	menu, err := s.menus.FindByCode(ctx, normalized)
	if err != nil {
		return dto.AdminMenuDetail{}, menuRepositoryError(err)
	}
	items, err := s.menus.FindItems(ctx, menu.ID)
	if err != nil {
		return dto.AdminMenuDetail{}, err
	}

	return dto.AdminMenuDetail{
		AdminMenuSummary: toMenuSummary(menu),
		Items:            buildMenuTree(items),
	}, nil
}

// Save 整体保存菜单项，Items 的顺序即展示顺序；菜单容器不存在时自动创建。
func (s *MenuService) Save(
	ctx context.Context,
	code string,
	request dto.AdminMenuSaveRequest,
	now int64,
) (dto.AdminMenuDetail, error) {
	normalized, err := normalizeMenuCode(code)
	if err != nil {
		return dto.AdminMenuDetail{}, err
	}

	items, err := toRepositoryItems(request.Items, 0)
	if err != nil {
		return dto.AdminMenuDetail{}, err
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		name = menuNames[normalized]
	}
	if len(name) > 100 {
		return dto.AdminMenuDetail{}, badRequest("菜单名称不能超过 100 个字符")
	}

	if _, err := s.menus.ReplaceItems(ctx, normalized, name, items, now); err != nil {
		return dto.AdminMenuDetail{}, err
	}
	s.invalidateCache(ctx)
	return s.Detail(ctx, normalized)
}

// navMenuCodes 是前台导航使用的三个菜单容器，顺序固定。
var navMenuCodes = []string{"header", "footer", "footer_bottom"}

// Nav 读取前台导航：顶部菜单、页脚栏目与底部链接。
// 只返回可见项，文案按 locale 取值，缺失时回退中文。
func (s *MenuService) Nav(ctx context.Context, locale string) (dto.SiteNav, error) {
	normalized, err := normalizeLocale(locale, true)
	if err != nil {
		return dto.SiteNav{}, err
	}

	items := make(map[string][]dto.NavItem, len(navMenuCodes))
	for _, code := range navMenuCodes {
		navItems := []dto.NavItem{}

		menu, err := s.menus.FindByCode(ctx, code)
		switch {
		case errors.Is(err, menurepository.ErrNotFound):
			// 菜单容器还没保存过：前台按空导航渲染。
		case err != nil:
			return dto.SiteNav{}, err
		default:
			rows, findErr := s.menus.FindItems(ctx, menu.ID)
			if findErr != nil {
				return dto.SiteNav{}, findErr
			}
			navItems = buildPublicNavItems(rows, normalized)
		}
		items[code] = navItems
	}

	return dto.SiteNav{
		Header:       items["header"],
		Footer:       items["footer"],
		FooterBottom: items["footer_bottom"],
	}, nil
}

// buildPublicNavItems 把菜单项按 parent_id 组装成两级树，并过滤不可见项。
func buildPublicNavItems(rows []models.MenuItem, locale string) []dto.NavItem {
	children := make(map[uint64][]models.MenuItem)
	roots := make([]models.MenuItem, 0, len(rows))
	for _, item := range rows {
		if item.ParentID == nil {
			roots = append(roots, item)
			continue
		}
		children[*item.ParentID] = append(children[*item.ParentID], item)
	}

	result := make([]dto.NavItem, 0, len(roots))
	for _, root := range roots {
		if root.Visible != 1 {
			continue
		}
		node := toPublicNavItem(root, locale)
		for _, child := range children[root.ID] {
			if child.Visible != 1 {
				continue
			}
			node.Children = append(node.Children, toPublicNavItem(child, locale))
		}
		result = append(result, node)
	}
	return result
}

// toPublicNavItem 把菜单项转换为前台结构，文案缺失时回退中文。
func toPublicNavItem(item models.MenuItem, locale string) dto.NavItem {
	labels := make(map[string]string)
	if len(item.Label) > 0 {
		_ = json.Unmarshal(item.Label, &labels)
	}
	label := strings.TrimSpace(labels[locale])
	if label == "" {
		label = strings.TrimSpace(labels["zh-CN"])
	}

	path := ""
	if item.LinkValue != nil {
		path = *item.LinkValue
	}

	return dto.NavItem{
		ID:       strconv.FormatUint(item.ID, 10),
		Label:    label,
		Path:     path,
		LinkType: item.LinkType,
		Target:   item.Target,
		Visible:  item.Visible == 1,
		Sort:     item.Sort,
	}
}

// invalidateCache 删除所有语言的导航分片；失败只记日志，不影响保存结果。
func (s *MenuService) invalidateCache(ctx context.Context) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}
	if err := s.cache.InvalidateNav(ctx); err != nil {
		hlog.CtxWarnf(ctx, "render cache: invalidate nav: %v", err)
	}
}

func normalizeMenuCode(code string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(code))
	if _, ok := menuNames[normalized]; !ok {
		return "", badRequest("不支持的菜单编码")
	}
	return normalized, nil
}

func toMenuSummary(menu *models.Menu) dto.AdminMenuSummary {
	description := ""
	if menu.Description != nil {
		description = *menu.Description
	}
	return dto.AdminMenuSummary{
		ID:          strconv.FormatUint(menu.ID, 10),
		Code:        menu.Code,
		Name:        menu.Name,
		Description: description,
		UpdatedAt:   menu.UpdatedAt,
	}
}

// buildMenuTree 把菜单项按 parent_id 组装成两级树。
func buildMenuTree(items []models.MenuItem) []dto.AdminMenuItem {
	children := make(map[uint64][]models.MenuItem)
	roots := make([]models.MenuItem, 0, len(items))
	for _, item := range items {
		if item.ParentID == nil {
			roots = append(roots, item)
			continue
		}
		children[*item.ParentID] = append(children[*item.ParentID], item)
	}

	result := make([]dto.AdminMenuItem, 0, len(roots))
	for _, root := range roots {
		node := toMenuItem(root)
		for _, child := range children[root.ID] {
			node.Children = append(node.Children, toMenuItem(child))
		}
		result = append(result, node)
	}
	return result
}

func toMenuItem(item models.MenuItem) dto.AdminMenuItem {
	label := make(map[string]string)
	if len(item.Label) > 0 {
		_ = json.Unmarshal(item.Label, &label)
	}
	path := ""
	if item.LinkValue != nil {
		path = *item.LinkValue
	}
	visible := item.Visible == 1
	return dto.AdminMenuItem{
		ID:       strconv.FormatUint(item.ID, 10),
		Label:    label,
		Path:     path,
		LinkType: item.LinkType,
		Target:   item.Target,
		Visible:  &visible,
		Sort:     item.Sort,
	}
}

// toRepositoryItems 校验并转换请求中的菜单项，子项（depth > 0）不再允许继续嵌套。
func toRepositoryItems(
	items []dto.AdminMenuItem,
	depth int,
) ([]menurepository.Item, error) {
	result := make([]menurepository.Item, 0, len(items))
	for index, item := range items {
		if depth > 0 && len(item.Children) > 0 {
			return nil, badRequest("菜单项最多支持两级")
		}

		label := make(map[string]string, len(item.Label))
		hasLabel := false
		for locale, value := range item.Label {
			trimmed := strings.TrimSpace(value)
			label[locale] = trimmed
			if trimmed != "" {
				hasLabel = true
			}
		}
		if !hasLabel {
			return nil, badRequest("请填写菜单文案")
		}
		labelJSON, err := json.Marshal(label)
		if err != nil {
			return nil, badRequest("菜单文案格式不正确")
		}

		linkType := strings.ToLower(strings.TrimSpace(item.LinkType))
		if linkType == "" {
			linkType = "page"
		}
		if _, ok := linkTypes[linkType]; !ok {
			return nil, badRequest("菜单链接类型不正确")
		}

		target := strings.TrimSpace(item.Target)
		if target == "" {
			target = "_self"
		}
		if _, ok := linkTargets[target]; !ok {
			return nil, badRequest("菜单打开方式不正确")
		}

		path := strings.TrimSpace(item.Path)
		if len(path) > 500 {
			return nil, badRequest("菜单链接不能超过 500 个字符")
		}
		var linkValue *string
		if path != "" {
			value := path
			linkValue = &value
		}

		visible := int8(1)
		if item.Visible != nil && !*item.Visible {
			visible = 0
		}

		sort := item.Sort
		if sort == 0 {
			sort = index
		}

		children, err := toRepositoryItems(item.Children, depth+1)
		if err != nil {
			return nil, err
		}

		result = append(result, menurepository.Item{
			Label:     labelJSON,
			LinkType:  linkType,
			LinkValue: linkValue,
			Target:    target,
			Sort:      sort,
			Visible:   visible,
			Children:  children,
		})
	}
	return result, nil
}

func menuRepositoryError(err error) error {
	if errors.Is(err, menurepository.ErrNotFound) {
		return ErrMenuNotFound
	}
	return err
}
