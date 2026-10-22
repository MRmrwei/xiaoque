package site

import (
	"encoding/json"

	"suki/xqyun/app/dto"
)

// 本文件导出前台渲染实现（app/services/site/impl）需要复用的规范化与解析逻辑，
// 避免在实现包里重复维护一套路径 / 语言规则。

// NormalizeLocale 校验并规范化语言标识；allowEmpty 为 true 时空值返回默认中文。
func NormalizeLocale(value string, allowEmpty bool) (string, error) {
	return normalizeLocale(value, allowEmpty)
}

// NormalizePath 规范化页面路径：补前导斜杠、去掉末尾斜杠。
func NormalizePath(raw string) string {
	return normalizePath(raw)
}

// ParsePageContent 解析页面 JSON 内容（seo + blocks），空值时返回带默认值的空结构。
func ParsePageContent(raw json.RawMessage) (dto.SEO, []dto.PageBlock, error) {
	content, err := pageContentOf(raw)
	if err != nil {
		return dto.SEO{}, nil, err
	}
	return content.SEO, content.Blocks, nil
}
