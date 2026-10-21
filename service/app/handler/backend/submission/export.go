package submission

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/response"
)

// Export 导出符合筛选条件的留言 CSV。
// GET /api/v1/admin/submissions/export
func (s *Controller) Export(ctx context.Context, c *app.RequestContext) {
	input, ok := listInput(c)
	if !ok {
		return
	}

	items, err := newService(s.App()).Export(ctx, input)
	if err != nil {
		fail(ctx, c, err)
		return
	}

	var buffer bytes.Buffer
	// UTF-8 BOM，保证 Excel 正确识别中文。
	buffer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(&buffer)
	records := [][]string{
		{"提交时间", "表单", "语言", "姓名", "邮箱", "电话", "公司", "留言", "状态", "IP", "浏览器"},
	}
	for _, item := range items {
		records = append(records, []string{
			formatTimestamp(item.CreatedAt),
			csvSafe(item.FormKey),
			localeLabel(item.Locale),
			csvSafe(item.Name),
			csvSafe(item.Email),
			csvSafe(item.Phone),
			csvSafe(item.Company),
			csvSafe(item.Message),
			statusLabel(item.Status),
			csvSafe(item.IP),
			csvSafe(item.UserAgent),
		})
	}
	if err := writer.WriteAll(records); err != nil {
		hlog.CtxErrorf(ctx, "admin submission export: %v", err)
		response.Fail(c, response.CodeInternal, "导出留言失败，请稍后重试")
		return
	}
	if err := writer.Error(); err != nil {
		hlog.CtxErrorf(ctx, "admin submission export: %v", err)
		response.Fail(c, response.CodeInternal, "导出留言失败，请稍后重试")
		return
	}

	filename := fmt.Sprintf("submissions-%s.csv", time.Now().Format("20060102"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	if _, err := c.Write(buffer.Bytes()); err != nil {
		hlog.CtxErrorf(ctx, "admin submission export: write response: %v", err)
	}
}

func formatTimestamp(seconds int64) string {
	if seconds <= 0 {
		return ""
	}
	return time.Unix(seconds, 0).Format("2006-01-02 15:04:05")
}

func localeLabel(locale string) string {
	switch locale {
	case "zh-CN":
		return "中文"
	case "en":
		return "English"
	default:
		return locale
	}
}

func statusLabel(status string) string {
	switch status {
	case "new":
		return "待处理"
	case "processing":
		return "处理中"
	case "done":
		return "已完成"
	case "spam":
		return "垃圾"
	default:
		return status
	}
}

// csvSafe 避免 Excel 把以 = / @ 开头的值当作公式执行。
// 不处理 + / -，避免电话（+86...）和负数被加上前缀。
func csvSafe(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return value
	}
	switch trimmed[0] {
	case '=', '@', '\t', '\r':
		return "'" + value
	default:
		return value
	}
}
