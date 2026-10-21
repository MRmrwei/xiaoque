package user

import (
	"context"

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

// LogoutController 处理管理后台退出登录接口。
type LogoutController struct {
	*controller.BaseController
}

// NewLogoutController 创建退出登录控制器，只保存 Application。
func NewLogoutController(app application.Application) *LogoutController {
	return &LogoutController{BaseController: controller.New(app)}
}

// Logout 读取 HttpOnly Cookie 中的刷新令牌并吊销，然后清除 Cookie。
func (l *LogoutController) Logout(ctx context.Context, c *app.RequestContext) {
	jwtConfig := l.App().JWT()
	adminJWT, ok := jwtConfig.Client("admin")
	if !ok {
		hlog.CtxErrorf(ctx, "admin logout: jwt client admin 未配置")
		response.Fail(c, response.CodeInternal, "退出登录失败，请稍后重试")
		return
	}

	refreshToken, _ := cookieutil.RefreshValue(c, jwtConfig.Cookie)

	userRepo := userrepository.New(l.App().DB())
	sessionRepo := sessionrepository.New(l.App().Redis())
	session := authimpl.NewAdminSession(userRepo, sessionRepo, adminAuthConfig(jwtConfig, adminJWT))
	logoutService := adminauth.NewSessionService(session)

	if err := logoutService.Logout(ctx, adminauth.LogoutInput{RefreshToken: refreshToken}); err != nil {
		hlog.CtxErrorf(ctx, "admin logout: %v", err)
		response.Fail(c, response.CodeInternal, "退出登录失败，请稍后重试")
		return
	}

	cookieutil.ClearRefresh(c, jwtConfig.Cookie)
	response.Success(c, nil)
}
