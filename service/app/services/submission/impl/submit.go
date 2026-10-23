package impl

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"suki/xqyun/app/dto"
	"suki/xqyun/app/models"
	submission "suki/xqyun/app/services/submission"
)

var (
	formKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,49}$`)
	emailPattern   = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

// Submit 校验并保存一条前台表单留言。
func (s *Service) Submit(
	ctx context.Context,
	input submission.SubmitInput,
) (dto.FrontendSubmissionResult, error) {
	formKey := strings.ToLower(strings.TrimSpace(input.FormKey))
	if !formKeyPattern.MatchString(formKey) {
		return dto.FrontendSubmissionResult{}, badRequest("表单标识不正确")
	}

	locale, err := normalizeLocale(input.Locale)
	if err != nil {
		return dto.FrontendSubmissionResult{}, err
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return dto.FrontendSubmissionResult{}, badRequest("请填写称呼")
	}
	if utf8.RuneCountInString(name) > 100 {
		return dto.FrontendSubmissionResult{}, badRequest("称呼不能超过 100 个字符")
	}

	phone := strings.TrimSpace(input.Phone)
	if phone == "" {
		return dto.FrontendSubmissionResult{}, badRequest("请填写联系电话")
	}
	if utf8.RuneCountInString(phone) > 50 {
		return dto.FrontendSubmissionResult{}, badRequest("联系电话不能超过 50 个字符")
	}

	email := strings.TrimSpace(input.Email)
	if email != "" {
		if utf8.RuneCountInString(email) > 255 || !emailPattern.MatchString(email) {
			return dto.FrontendSubmissionResult{}, badRequest("请输入有效的邮箱地址")
		}
	}

	company := strings.TrimSpace(input.Company)
	if utf8.RuneCountInString(company) > 200 {
		return dto.FrontendSubmissionResult{}, badRequest("公司名称不能超过 200 个字符")
	}

	message := strings.TrimSpace(input.Message)
	if message == "" {
		return dto.FrontendSubmissionResult{}, badRequest("请填写留言内容")
	}
	if utf8.RuneCountInString(message) > 5000 {
		return dto.FrontendSubmissionResult{}, badRequest("留言内容不能超过 5000 个字符")
	}

	now := input.Now
	if now <= 0 {
		now = time.Now().Unix()
	}

	data, err := json.Marshal(map[string]string{
		"name":    name,
		"email":   email,
		"company": company,
		"phone":   phone,
		"message": message,
	})
	if err != nil {
		return dto.FrontendSubmissionResult{}, err
	}

	row := models.FormSubmission{
		FormKey:   formKey,
		Locale:    nullable(locale),
		Name:      nullable(name),
		Email:     nullable(email),
		Phone:     nullable(phone),
		Company:   nullable(company),
		Message:   nullable(message),
		Data:      data,
		Status:    submission.StatusNew,
		IP:        nullable(input.IP),
		UserAgent: nullable(truncate(input.UserAgent, 500)),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.submissions.Create(ctx, &row); err != nil {
		return dto.FrontendSubmissionResult{}, err
	}

	return dto.FrontendSubmissionResult{ID: strconv.FormatUint(row.ID, 10)}, nil
}

// nullable 把空字符串转换为 NULL，便于后台区分「未填写」。
func nullable(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// truncate 按字符数截断过长的请求头，避免超出数据库列长度。
func truncate(value string, limit int) string {
	value = strings.TrimSpace(value)
	if utf8.RuneCountInString(value) <= limit {
		return value
	}
	runes := []rune(value)
	return string(runes[:limit])
}
