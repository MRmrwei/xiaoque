package impl

import (
	"context"
	"encoding/json"
	"fmt"
)

// siteConfig 读取脱敏后的站点配置：全局项与语言级设置合并，同 key 时语言级覆盖全局。
// 只返回 xqyun_settings 中已保存的业务配置，不读取任何后台专属数据。
func (s *Renderer) siteConfig(ctx context.Context, locale string) (map[string]any, error) {
	rows, err := s.settings.List(ctx, locale)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(rows)+1)
	result["locale"] = locale
	for _, row := range rows {
		var value any
		if err := json.Unmarshal(row.Value, &value); err != nil {
			return nil, fmt.Errorf("解析站点设置 %s 失败：%w", row.Key, err)
		}
		result[row.Key] = value
	}
	return result, nil
}
