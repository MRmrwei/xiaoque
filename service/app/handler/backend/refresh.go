package backend

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
	sessionrepository "suki/xqyun/app/repository/session"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/adminauth"
	authimpl "suki/xqyun/app/services/auth/impl"
	cookieutil "suki/xqyun/app/utils/cookie"
)

// RefreshController 处理管理后台刷新 Access Token 接口。
type RefreshController struct {
	*controller.BaseController
}

// NewRefreshController 创建刷新令牌控制器，只保存 Application。
func NewRefreshController(app application.Application) *RefreshController {
	return &RefreshController{BaseController: controller.New(app)}
}

// Refresh 从 HttpOnly Cookie 读取刷新令牌，轮换后下发新的 Access Token。
func (r *RefreshController) Refresh(ctx context.Context, c *app.RequestContext) {
	jwtConfig := r.App().JWT()
	adminJWT, ok := jwtConfig.Client("admin")
	if !ok {
		hlog.CtxErrorf(ctx, "admin refresh: jwt client admin 未配置")
		response.Fail(c, response.CodeInternal, "刷新登录状态失败，请稍后重试")
		return
	}

	refreshToken, ok := cookieutil.RefreshValue(c, jwtConfig.Cookie)
	if !ok {
		// 没有刷新令牌等同于未登录，顺手清掉可能残留的空 Cookie。
		cookieutil.ClearRefresh(c, jwtConfig.Cookie)
		response.Fail(c, adminauth.ErrSessionExpired.Code, adminauth.ErrSessionExpired.Message)
		return
	}

	userRepo := userrepository.New(r.App().DB())
	sessionRepo := sessionrepository.New(r.App().Redis())
	session := authimpl.NewAdminSession(userRepo, sessionRepo, adminAuthConfig(jwtConfig, adminJWT))
	refreshService := adminauth.NewSessionService(session)

	result, err := refreshService.Refresh(ctx, adminauth.RefreshInput{
		RefreshToken: refreshToken,
		UserAgent:    string(c.UserAgent()),
		IP:           c.ClientIP(),
	})
	if err != nil {
		// 只有登录状态确实失效时才清 Cookie；数据库等临时故障保留 Cookie，避免用户被迫重新登录。
		if errors.Is(err, adminauth.ErrSessionExpired) || errors.Is(err, adminauth.ErrAccountDisabled) {
			cookieutil.ClearRefresh(c, jwtConfig.Cookie)
		}
		r.fail(ctx, c, err)
		return
	}

	cookieutil.SetRefresh(
		c,
		result.RefreshToken,
		int(adminJWT.RefreshTTL.Seconds()),
		jwtConfig.Cookie,
	)
	response.Success(c, result)
}

// fail 把刷新用例的业务错误转换为统一响应；未知错误统一返回 500 并记录日志。
func (r *RefreshController) fail(ctx context.Context, c *app.RequestContext, err error) {
	var businessError *adminauth.Error
	if errors.As(err, &businessError) {
		response.Fail(c, businessError.Code, businessError.Message)
		return
	}

	hlog.CtxErrorf(ctx, "admin refresh: %v", err)
	response.Fail(c, response.CodeInternal, "刷新登录状态失败，请稍后重试")
}
