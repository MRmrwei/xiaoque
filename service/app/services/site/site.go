package site

import (
	"context"

	"suki/xqyun/app/dto"
)

// RenderPageInput 是前台页面渲染查询条件。
type RenderPageInput struct {
	// Path 是页面路径，如 /about、/en/about-us；为空时按首页 / 处理。
	Path string
	// Locale 为 zh-CN / en，为空时默认 zh-CN。
	Locale string
}

// PageRenderer 是官网前台页面渲染数据的用例契约。
// 实现只读取线上内容：页面取 published_content，导航与站点配置取已保存的设置；
// 数据访问放在 repository 层，用例层负责聚合、语言回退与 alternates 计算。
type PageRenderer interface {
	// RenderPage 按路径聚合已发布页面、顶部 / 底部导航、站点配置与数据源区块集合。
	// 页面不存在或未发布时返回 ErrPageNotFound，前台据此输出 404 + noindex。
	RenderPage(ctx context.Context, input RenderPageInput) (dto.RenderPageData, error)
}
