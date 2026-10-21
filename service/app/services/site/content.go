package site

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"suki/xqyun/app/dto"
	"suki/xqyun/app/utils/ulid"
)

// 页面支持的中英文标识。前台内部的 zh 在接口层统一映射为 zh-CN。
const (
	localeZH = "zh-CN"
	localeEN = "en"
)

// pageContent 是 xqyun_pages.draft_content / published_content 的结构。
type pageContent struct {
	SEO    dto.SEO         `json:"seo"`
	Blocks []dto.PageBlock `json:"blocks"`
}

// pageContentOf 解析页面 JSON 内容，空值时返回带默认值的空结构。
func pageContentOf(raw json.RawMessage) (pageContent, error) {
	content := pageContent{
		SEO:    dto.SEO{Keywords: []string{}},
		Blocks: []dto.PageBlock{},
	}
	if len(raw) == 0 || string(raw) == "null" {
		return content, nil
	}
	if err := json.Unmarshal(raw, &content); err != nil {
		return pageContent{}, fmt.Errorf("页面内容解析失败：%w", err)
	}
	if content.SEO.Keywords == nil {
		content.SEO.Keywords = []string{}
	}
	if content.Blocks == nil {
		content.Blocks = []dto.PageBlock{}
	}
	return content, nil
}

// emptyPageContent 生成新建页面的初始内容。
func emptyPageContent(name string) pageContent {
	return pageContent{
		SEO: dto.SEO{
			Title:    name,
			Keywords: []string{},
		},
		Blocks: []dto.PageBlock{},
	}
}

// normalizeSEO 去除首尾空白，并保证 keywords 输出为数组而不是 null。
func normalizeSEO(seo *dto.SEO) {
	seo.Title = strings.TrimSpace(seo.Title)
	seo.Description = strings.TrimSpace(seo.Description)
	if seo.Keywords == nil {
		seo.Keywords = []string{}
	}
}

// normalizeLocale 校验并规范化语言标识；allowEmpty 为 true 时空值返回默认中文。
func normalizeLocale(value string, allowEmpty bool) (string, error) {
	locale := strings.TrimSpace(value)
	switch locale {
	case "", "zh", "zh-CN":
		if locale == "" && !allowEmpty {
			return "", badRequest("请选择页面语言")
		}
		return localeZH, nil
	case "en":
		return localeEN, nil
	default:
		return "", badRequest("语言参数不正确")
	}
}

// otherLocale 返回另一种语言。
func otherLocale(locale string) string {
	if locale == localeEN {
		return localeZH
	}
	return localeEN
}

// normalizePath 规范化页面路径：补前导斜杠、去掉末尾斜杠。
func normalizePath(raw string) string {
	path := strings.TrimSpace(raw)
	if path == "" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	if path == "" {
		path = "/"
	}
	return path
}

// newBlockID 生成翻译 / 复制页面时的新区块 ID。
func newBlockID() string {
	value, err := ulid.New()
	if err != nil {
		return fmt.Sprintf("blk_%d", time.Now().UnixNano())
	}
	return "blk_" + strings.ToLower(value[:10])
}

// keepTextKeys 是复制翻译结构时需要保留的字段：
// 链接、图片、主题、版式等与语言无关的值不参与清空。
var keepTextKeys = map[string]struct{}{
	"id":                {},
	"href":              {},
	"url":               {},
	"videoUrl":          {},
	"primaryHref":       {},
	"secondaryHref":     {},
	"ctaHref":           {},
	"contactHref":       {},
	"moreHref":          {},
	"modalCtaHref":      {},
	"path":              {},
	"link":              {},
	"linkValue":         {},
	"key":               {},
	"email":             {},
	"phone":             {},
	"slug":              {},
	"value":             {},
	"suffix":            {},
	"metric":            {},
	"metricSuffix":      {},
	"icon":              {},
	"variant":           {},
	"theme":             {},
	"layout":            {},
	"tone":              {},
	"visual":            {},
	"cover":             {},
	"mapStyle":          {},
	"imageSide":         {},
	"width":             {},
	"image":             {},
	"images":            {},
	"media":             {},
	"mediaType":         {},
	"videoPoster":       {},
	"backgroundImage":   {},
	"limit":             {},
	"interval":          {},
	"autoplay":          {},
	"showFilters":       {},
	"marquee":           {},
	"date":              {},
	"year":              {},
	"readMinutes":       {},
	"backgroundOverlay": {},
}

// clearTranslatableBlocks 复制区块结构并生成新区块 ID，清空标题、正文等可翻译文本。
// 图片、视频、链接、图标与版式字段保留；非字符串值（数字、布尔）原样保留。
func clearTranslatableBlocks(blocks []dto.PageBlock) []dto.PageBlock {
	result := make([]dto.PageBlock, 0, len(blocks))
	for _, block := range blocks {
		result = append(result, dto.PageBlock{
			ID:    newBlockID(),
			Type:  block.Type,
			Props: clearTranslatableProps(block.Props),
		})
	}
	return result
}

func clearTranslatableProps(props map[string]any) map[string]any {
	result := make(map[string]any, len(props))
	for key, value := range props {
		if _, keep := keepTextKeys[key]; keep {
			result[key] = cloneJSONValue(value)
			continue
		}
		result[key] = clearTranslatableValue(value)
	}
	return result
}

func clearTranslatableValue(value any) any {
	switch typed := value.(type) {
	case string:
		return ""
	case map[string]any:
		return clearTranslatableProps(typed)
	case []any:
		items := make([]any, len(typed))
		for index := range typed {
			items[index] = clearTranslatableValue(typed[index])
		}
		return items
	default:
		return value
	}
}

func cloneJSONValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		cloned := make(map[string]any, len(typed))
		for key, child := range typed {
			cloned[key] = cloneJSONValue(child)
		}
		return cloned
	case []any:
		items := make([]any, len(typed))
		for index := range typed {
			items[index] = cloneJSONValue(typed[index])
		}
		return items
	default:
		return value
	}
}

// regenerateBlockIDs 复制页面时重新生成区块 ID，避免不同页面出现相同区块 ID。
func regenerateBlockIDs(blocks []dto.PageBlock) []dto.PageBlock {
	result := make([]dto.PageBlock, 0, len(blocks))
	for _, block := range blocks {
		props := make(map[string]any)
		if block.Props != nil {
			if cloned, ok := cloneJSONValue(block.Props).(map[string]any); ok {
				props = cloned
			}
		}
		result = append(result, dto.PageBlock{
			ID:    newBlockID(),
			Type:  block.Type,
			Props: props,
		})
	}
	return result
}
