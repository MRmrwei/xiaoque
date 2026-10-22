package impl

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"suki/xqyun/app/models"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/user"
)

const (
	maxNameRunes   = 100
	maxEmailRunes  = 255
	maxMobileRunes = 32
	maxAvatarBytes = 2 << 20 // 2MB
)

// ProfileUpdater 是 user.ProfileUpdater 的实现。
type ProfileUpdater struct {
	users userrepository.UserRepository
	now   func() time.Time
}

var _ user.ProfileUpdater = (*ProfileUpdater)(nil)

// NewProfileUpdater 创建修改资料用例。
func NewProfileUpdater(users userrepository.UserRepository) *ProfileUpdater {
	return &ProfileUpdater{users: users, now: time.Now}
}

// UpdateProfile 校验并更新显示名称、邮箱、手机号与头像，返回更新后的账号信息。
func (s *ProfileUpdater) UpdateProfile(
	ctx context.Context,
	input user.UpdateProfileInput,
) (*models.User, error) {
	if input.UserID == 0 {
		return nil, user.ErrUserNotFound
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, user.ErrNameRequired
	}
	if utf8.RuneCountInString(name) > maxNameRunes {
		return nil, user.ErrNameTooLong
	}

	email := strings.TrimSpace(input.Email)
	if utf8.RuneCountInString(email) > maxEmailRunes {
		return nil, user.ErrEmailTooLong
	}
	if email != "" {
		parsed, err := mail.ParseAddress(email)
		if err != nil || parsed.Address != email {
			return nil, user.ErrEmailInvalid
		}
	}

	mobile := strings.TrimSpace(input.Mobile)
	if utf8.RuneCountInString(mobile) > maxMobileRunes {
		return nil, user.ErrMobileTooLong
	}

	avatar := strings.TrimSpace(input.Avatar)
	if len(avatar) > maxAvatarBytes {
		return nil, user.ErrAvatarTooLarge
	}

	if _, err := s.findUser(ctx, input.UserID); err != nil {
		return nil, err
	}

	err := s.users.UpdateProfile(
		ctx,
		input.UserID,
		name,
		optionalString(email),
		optionalString(mobile),
		optionalString(avatar),
		s.now().Unix(),
	)
	if errors.Is(err, userrepository.ErrNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.findUser(ctx, input.UserID)
}

func (s *ProfileUpdater) findUser(ctx context.Context, userID uint64) (*models.User, error) {
	account, err := s.users.FindByID(ctx, userID)
	if errors.Is(err, userrepository.ErrNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

// optionalString 把空字符串转换为 nil，表示数据库存 NULL。
func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
