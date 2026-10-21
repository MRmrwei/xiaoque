package adminauth

import (
	"context"
	"errors"
	"testing"

	"suki/xqyun/app/models"
	"suki/xqyun/app/services/auth"
)

type fakeAuthenticator struct {
	result   *auth.LoginResult
	err      error
	received auth.Credentials
	calls    int
}

func (f *fakeAuthenticator) Login(_ context.Context, credentials auth.Credentials) (*auth.LoginResult, error) {
	f.calls++
	f.received = credentials
	return f.result, f.err
}

func stringPtr(value string) *string {
	return &value
}

func TestLoginSuccess(t *testing.T) {
	fake := &fakeAuthenticator{
		result: &auth.LoginResult{
			Pair: &auth.TokenPair{
				AccessToken:  "jwt-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    900,
			},
			User: &models.User{
				ID:       7,
				Username: "admin",
				Name:     "小阙",
				Email:    stringPtr("admin@example.com"),
				Mobile:   stringPtr("13800000000"),
				Avatar:   stringPtr("data:image/webp;base64,abc"),
			},
		},
	}

	response, err := NewLoginService(fake).Login(context.Background(), LoginInput{
		Username:  " admin ",
		Password:  "123456",
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("Login() 返回错误：%v", err)
	}
	if fake.calls != 1 {
		t.Fatalf("认证服务调用次数 = %d，期望 1", fake.calls)
	}
	if fake.received.Username != "admin" {
		t.Fatalf("用户名未去除空格：%q", fake.received.Username)
	}
	if fake.received.Password != "123456" ||
		fake.received.UserAgent != "test-agent" ||
		fake.received.IP != "127.0.0.1" {
		t.Fatalf("认证参数未正确传递：%#v", fake.received)
	}

	if response.AccessToken != "jwt-token" ||
		response.RefreshToken != "refresh-token" ||
		response.ExpiresIn != 900 {
		t.Fatalf("令牌响应不正确：%#v", response)
	}
	if response.User.ID != "7" ||
		response.User.Username != "admin" ||
		response.User.Name != "小阙" ||
		response.User.Email != "admin@example.com" ||
		response.User.Mobile != "13800000000" ||
		response.User.Avatar != "data:image/webp;base64,abc" {
		t.Fatalf("用户信息不正确：%#v", response.User)
	}
	if response.User.Roles == nil || len(response.User.Roles) != 0 {
		t.Fatalf("roles 应为空数组：%#v", response.User.Roles)
	}
	if response.User.Permissions == nil || len(response.User.Permissions) != 0 {
		t.Fatalf("permissions 应为空数组：%#v", response.User.Permissions)
	}
}

func TestLoginRejectsEmptyInput(t *testing.T) {
	inputs := []LoginInput{
		{},
		{Username: "admin"},
		{Username: "   ", Password: "123456"},
	}

	for _, input := range inputs {
		if _, err := NewLoginService(&fakeAuthenticator{}).Login(context.Background(), input); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("输入 %#v 应返回 ErrInvalidInput，实际 %v", input, err)
		}
	}
}

func TestLoginMapsAuthErrors(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want *Error
	}{
		{name: "invalid credentials", err: auth.ErrInvalidCredentials, want: ErrInvalidCredentials},
		{name: "account disabled", err: auth.ErrAccountDisabled, want: ErrAccountDisabled},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := NewLoginService(&fakeAuthenticator{err: testCase.err}).Login(context.Background(), LoginInput{
				Username: "admin",
				Password: "123456",
			})
			if !errors.Is(err, testCase.want) {
				t.Fatalf("错误映射不正确：期望 %v，实际 %v", testCase.want, err)
			}
		})
	}
}

func TestLoginPassesThroughUnknownError(t *testing.T) {
	unknown := errors.New("database unavailable")
	if _, err := NewLoginService(&fakeAuthenticator{err: unknown}).Login(context.Background(), LoginInput{
		Username: "admin",
		Password: "123456",
	}); !errors.Is(err, unknown) {
		t.Fatalf("未知错误应原样返回，实际 %v", err)
	}
}

func TestLoginRejectsNilResult(t *testing.T) {
	if _, err := NewLoginService(&fakeAuthenticator{}).Login(context.Background(), LoginInput{
		Username: "admin",
		Password: "123456",
	}); err == nil {
		t.Fatal("认证服务返回空结果时应报错")
	}
}

func TestLoginRejectsNilAuthenticator(t *testing.T) {
	if _, err := NewLoginService(nil).Login(context.Background(), LoginInput{
		Username: "admin",
		Password: "123456",
	}); err == nil {
		t.Fatal("未注入认证服务时应报错")
	}
}

type fakeSession struct {
	refreshResult   *auth.LoginResult
	refreshErr      error
	receivedRefresh auth.RefreshCredentials
	refreshCalls    int

	logoutErr      error
	receivedLogout auth.LogoutCredentials
	logoutCalls    int
}

var _ auth.SessionManager = (*fakeSession)(nil)

func (f *fakeSession) Refresh(_ context.Context, credentials auth.RefreshCredentials) (*auth.LoginResult, error) {
	f.refreshCalls++
	f.receivedRefresh = credentials
	return f.refreshResult, f.refreshErr
}

func (f *fakeSession) Logout(_ context.Context, credentials auth.LogoutCredentials) error {
	f.logoutCalls++
	f.receivedLogout = credentials
	return f.logoutErr
}

func TestRefreshSuccess(t *testing.T) {
	session := &fakeSession{
		refreshResult: &auth.LoginResult{
			Pair: &auth.TokenPair{
				AccessToken:  "new-jwt",
				RefreshToken: "new-refresh",
				ExpiresIn:    900,
			},
			User: &models.User{ID: 7, Username: "admin", Name: "小阙"},
		},
	}

	response, err := NewSessionService(session).Refresh(context.Background(), RefreshInput{
		RefreshToken: " refresh-token ",
		UserAgent:    "test-agent",
		IP:           "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("Refresh() 返回错误：%v", err)
	}
	if session.refreshCalls != 1 {
		t.Fatalf("刷新服务调用次数 = %d，期望 1", session.refreshCalls)
	}
	if session.receivedRefresh.RefreshToken != "refresh-token" ||
		session.receivedRefresh.UserAgent != "test-agent" ||
		session.receivedRefresh.IP != "127.0.0.1" {
		t.Fatalf("刷新参数未正确传递：%#v", session.receivedRefresh)
	}
	if response.AccessToken != "new-jwt" ||
		response.RefreshToken != "new-refresh" ||
		response.ExpiresIn != 900 ||
		response.User.ID != "7" {
		t.Fatalf("刷新响应不正确：%#v", response)
	}
}

func TestRefreshRejectsEmptyToken(t *testing.T) {
	session := &fakeSession{}
	if _, err := NewSessionService(session).Refresh(context.Background(), RefreshInput{}); !errors.Is(err, ErrSessionExpired) {
		t.Fatalf("应返回 ErrSessionExpired，实际 %v", err)
	}
	if session.refreshCalls != 0 {
		t.Fatal("缺少刷新令牌时不应调用刷新服务")
	}
}

func TestRefreshMapsAuthErrors(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want *Error
	}{
		{name: "invalid refresh token", err: auth.ErrInvalidRefreshToken, want: ErrSessionExpired},
		{name: "account disabled", err: auth.ErrAccountDisabled, want: ErrAccountDisabled},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := NewSessionService(&fakeSession{refreshErr: testCase.err}).Refresh(context.Background(), RefreshInput{
				RefreshToken: "refresh-token",
			})
			if !errors.Is(err, testCase.want) {
				t.Fatalf("错误映射不正确：期望 %v，实际 %v", testCase.want, err)
			}
		})
	}
}

func TestRefreshPassesThroughUnknownError(t *testing.T) {
	unknown := errors.New("database unavailable")
	if _, err := NewSessionService(&fakeSession{refreshErr: unknown}).Refresh(context.Background(), RefreshInput{
		RefreshToken: "refresh-token",
	}); !errors.Is(err, unknown) {
		t.Fatalf("未知错误应原样返回，实际 %v", err)
	}
}

func TestLogoutCallsSession(t *testing.T) {
	session := &fakeSession{}
	if err := NewSessionService(session).Logout(context.Background(), LogoutInput{
		RefreshToken: " refresh-token ",
	}); err != nil {
		t.Fatalf("Logout() 返回错误：%v", err)
	}
	if session.logoutCalls != 1 || session.receivedLogout.RefreshToken != "refresh-token" {
		t.Fatalf("退出登录参数未正确传递：%#v", session.receivedLogout)
	}
}

func TestLogoutIgnoresEmptyToken(t *testing.T) {
	session := &fakeSession{}
	if err := NewSessionService(session).Logout(context.Background(), LogoutInput{}); err != nil {
		t.Fatalf("缺少刷新令牌时应幂等成功，实际 %v", err)
	}
	if session.logoutCalls != 0 {
		t.Fatal("缺少刷新令牌时不应调用退出登录服务")
	}
}
