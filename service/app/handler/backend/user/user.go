package user

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	"suki/xqyun/app/models"
	sessionrepository "suki/xqyun/app/repository/session"
	userrepository "suki/xqyun/app/repository/user"
	userservice "suki/xqyun/app/services/user"
	userimpl "suki/xqyun/app/services/user/impl"
	cookieutil "suki/xqyun/app/utils/cookie"
	"suki/xqyun/middleware"
)

// UserController 处理当前管理员的资料与密码。
// 路由前缀：/api/v1/admin/user
type UserController struct {
	*controller.BaseController
}

// NewUserController 创建管理员资料 / 密码控制器。
func NewUserController(app application.Application) *UserController {
	return &UserController{BaseController: controller.New(app)}
}

// UpdateProfile 修改当前管理员资料（显示名称、邮箱、手机号、头像）。
// PUT /api/v1/admin/user/profile
func (u *UserController) UpdateProfile(ctx context.Context, c *app.RequestContext) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized, "请先登录")
		return
	}

	var req dto.AdminProfileUpdateRequest
	if err := c.Bind(&req); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	account, err := userimpl.NewProfileUpdater(
		userrepository.New(u.App().DB()),
	).UpdateProfile(ctx, userservice.UpdateProfileInput{
		UserID: userID,
		Name:   req.Name,
		Email:  req.Email,
		Mobile: req.Mobile,
		Avatar: req.Avatar,
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, toAdminProfile(account))
}

// ChangePassword 修改当前管理员密码。
// PUT /api/v1/admin/user/password
func (u *UserController) ChangePassword(ctx context.Context, c *app.RequestContext) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized, "请先登录")
		return
	}

	var req dto.AdminPasswordChangeRequest
	if err := c.Bind(&req); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	jwtConfig := u.App().JWT()
	tokenType := "admin"
	if client, ok := jwtConfig.Client("admin"); ok && client.TokenType != "" {
		tokenType = client.TokenType
	}
	refreshToken, _ := cookieutil.RefreshValue(c, jwtConfig.Cookie)

	err := userimpl.NewPasswordChanger(
		userrepository.New(u.App().DB()),
		sessionrepository.New(u.App().Redis()),
		tokenType,
	).ChangePassword(ctx, userservice.ChangePasswordInput{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
		RefreshToken:    refreshToken,
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}

	// 改密后刷新令牌已全部吊销，清掉当前浏览器的 Cookie，要求重新登录。
	cookieutil.ClearRefresh(c, jwtConfig.Cookie)
	response.Success(c, nil)
}

// currentUserID 从 Access Token claims 读取当前管理员 ID。
func currentUserID(c *app.RequestContext) (uint64, bool) {
	claims, ok := middleware.AdminClaimsFromContext(c)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseUint(strings.TrimSpace(claims.Subject), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

// toAdminProfile 把账号模型转换为前台响应结构。
func toAdminProfile(account *models.User) dto.AdminProfile {
	profile := dto.AdminProfile{
		ID:       strconv.FormatUint(account.ID, 10),
		Username: account.Username,
		Name:     account.Name,
	}
	if account.Email != nil {
		profile.Email = *account.Email
	}
	if account.Mobile != nil {
		profile.Mobile = *account.Mobile
	}
	if account.Avatar != nil {
		profile.Avatar = *account.Avatar
	}
	return profile
}

// fail 把业务错误映射为统一响应；未知错误记录日志并返回 500。
func fail(ctx context.Context, c *app.RequestContext, err error) {
	switch {
	case errors.Is(err, userservice.ErrUserNotFound):
		response.Fail(c, response.CodeNotFound, err.Error())
	case errors.Is(err, userservice.ErrCurrentPasswordIncorrect),
		errors.Is(err, userservice.ErrPasswordTooShort),
		errors.Is(err, userservice.ErrPasswordTooLong),
		errors.Is(err, userservice.ErrNameRequired),
		errors.Is(err, userservice.ErrNameTooLong),
		errors.Is(err, userservice.ErrEmailInvalid),
		errors.Is(err, userservice.ErrEmailTooLong),
		errors.Is(err, userservice.ErrMobileTooLong),
		errors.Is(err, userservice.ErrAvatarTooLarge):
		response.Fail(c, response.CodeBadRequest, err.Error())
	default:
		hlog.CtxErrorf(ctx, "admin user: %v", err)
		response.Fail(c, response.CodeInternal, "操作失败，请稍后重试")
	}
}
