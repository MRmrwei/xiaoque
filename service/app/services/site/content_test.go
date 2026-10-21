package site

import (
	"testing"

	"suki/xqyun/app/dto"
)

func TestNormalizePath(t *testing.T) {
	cases := map[string]string{
		"":             "",
		"/":            "/",
		"about":        "/about",
		"/about/":      "/about",
		" /en/about/ ": "/en/about",
		"/en/about///": "/en/about",
	}
	for input, expected := range cases {
		if actual := normalizePath(input); actual != expected {
			t.Fatalf("normalizePath(%q) = %q，期望 %q", input, actual, expected)
		}
	}
}

func TestNormalizeLocale(t *testing.T) {
	cases := []struct {
		value      string
		allowEmpty bool
		expected   string
		wantError  bool
	}{
		{value: "", allowEmpty: true, expected: localeZH},
		{value: "zh", allowEmpty: true, expected: localeZH},
		{value: "zh-CN", allowEmpty: true, expected: localeZH},
		{value: "en", allowEmpty: false, expected: localeEN},
		{value: "", allowEmpty: false, wantError: true},
		{value: "fr", allowEmpty: true, wantError: true},
	}
	for _, item := range cases {
		actual, err := normalizeLocale(item.value, item.allowEmpty)
		if item.wantError {
			if err == nil {
				t.Fatalf("normalizeLocale(%q, %v) 期望报错", item.value, item.allowEmpty)
			}
			continue
		}
		if err != nil {
			t.Fatalf("normalizeLocale(%q, %v) error = %v", item.value, item.allowEmpty, err)
		}
		if actual != item.expected {
			t.Fatalf(
				"normalizeLocale(%q, %v) = %q，期望 %q",
				item.value,
				item.allowEmpty,
				actual,
				item.expected,
			)
		}
	}
}

func TestClearTranslatableBlocks(t *testing.T) {
	source := []dto.PageBlock{
		{
			ID:   "old_block",
			Type: "hero",
			Props: map[string]any{
				"title":       "主标题",
				"subtitle":    "副标题",
				"image":       "builder",
				"primaryHref": "/contact",
				"theme":       "light",
				"limit":       float64(3),
				"slides": []any{
					map[string]any{
						"id":        "slide_1",
						"mediaType": "image",
						"media":     "orbit",
						"title":     "幻灯片标题",
					},
				},
			},
		},
	}

	result := clearTranslatableBlocks(source)
	if len(result) != 1 {
		t.Fatalf("区块数量 = %d，期望 1", len(result))
	}
	block := result[0]
	if block.ID == "old_block" {
		t.Fatal("翻译复制时没有生成新的区块 ID")
	}
	if block.Type != "hero" {
		t.Fatalf("区块类型 = %q，期望 hero", block.Type)
	}

	if block.Props["title"] != "" || block.Props["subtitle"] != "" {
		t.Fatal("文本字段没有被清空")
	}
	if block.Props["image"] != "builder" {
		t.Fatal("图片字段被误清空")
	}
	if block.Props["primaryHref"] != "/contact" {
		t.Fatal("链接字段被误清空")
	}
	if block.Props["theme"] != "light" {
		t.Fatal("主题字段被误清空")
	}
	if block.Props["limit"] != float64(3) {
		t.Fatal("数字字段被误修改")
	}

	slides, ok := block.Props["slides"].([]any)
	if !ok || len(slides) != 1 {
		t.Fatalf("slides 结构不正确：%#v", block.Props["slides"])
	}
	slide, ok := slides[0].(map[string]any)
	if !ok {
		t.Fatalf("slide 结构不正确：%#v", slides[0])
	}
	if slide["title"] != "" {
		t.Fatal("轮播项标题没有被清空")
	}
	if slide["media"] != "orbit" || slide["mediaType"] != "image" || slide["id"] != "slide_1" {
		t.Fatal("轮播项的非文本字段被误清空")
	}

	if source[0].Props["title"] != "主标题" {
		t.Fatal("清空操作修改了源页面数据")
	}
}
