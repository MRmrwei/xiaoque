// Package backend 提供管理后台的 HTTP 控制器（登录、刷新令牌、退出登录等）。
package backend

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	sessionrepository "suki/xqyun/app/repository/session"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/adminauth"
	authimpl "suki/xqyun/app/services/auth/impl"
	cookieutil "suki/xqyun/app/utils/cookie"
)

// LoginController 处理管理后台登录接口。
type LoginController struct {
	*controller.BaseController
}

// NewLoginController 创建登录控制器，只保存 Application。
func NewLoginController(app application.Application) *LoginController {
	return &LoginController{BaseController: controller.New(app)}
}

// Login 解析请求、创建登录用例服务并输出响应。
func (l *LoginController) Login(ctx context.Context, c *app.RequestContext) {
	var req dto.AdminLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	jwtConfig := l.App().JWT()
	adminJWT, ok := jwtConfig.Client("admin")
	if !ok {
		hlog.CtxErrorf(ctx, "admin login: jwt client admin 未配置")
		response.Fail(c, response.CodeInternal, "登录失败，请稍后重试")
		return
	}

	userRepo := userrepository.New(l.App().DB())
	sessionRepo := sessionrepository.New(l.App().Redis())
	authenticator := authimpl.NewAdminLogin(userRepo, sessionRepo, adminAuthConfig(jwtConfig, adminJWT))
	loginService := adminauth.NewLoginService(authenticator)

	result, err := loginService.Login(ctx, adminauth.LoginInput{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: string(c.UserAgent()),
		IP:        c.ClientIP(),
	})
	if err != nil {
		l.fail(ctx, c, err)
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

// adminAuthConfig 把应用级 JWT 配置转换为认证实现所需的参数。
// 登录、刷新令牌、退出登录共用同一套客户端配置。
func adminAuthConfig(
	jwtConfig application.JWTConfig,
	client application.JWTClientConfig,
) authimpl.AdminLoginConfig {
	return authimpl.AdminLoginConfig{
		Secret:     client.Secret,
		Issuer:     jwtConfig.Issuer,
		Audience:   client.Audience,
		TokenType:  client.TokenType,
		AccessTTL:  client.AccessTTL,
		RefreshTTL: client.RefreshTTL,
	}
}

// fail 把登录用例的业务错误转换为统一响应；未知错误统一返回 500 并记录日志。
func (l *LoginController) fail(ctx context.Context, c *app.RequestContext, err error) {
	var businessError *adminauth.Error
	if errors.As(err, &businessError) {
		response.Fail(c, businessError.Code, businessError.Message)
		return
	}

	hlog.CtxErrorf(ctx, "admin login: %v", err)
	response.Fail(c, response.CodeInternal, "登录失败，请稍后重试")
}
