package impl

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"suki/xqyun/app/models"
	sessionrepository "suki/xqyun/app/repository/session"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/auth"
)

type fakeUsers struct {
	users map[uint64]*models.User
}

var _ userrepository.UserRepository = (*fakeUsers)(nil)

func newFakeUsers() *fakeUsers {
	return &fakeUsers{users: map[uint64]*models.User{}}
}

func (f *fakeUsers) FindByUsername(_ context.Context, username string) (*models.User, error) {
	for _, user := range f.users {
		if user.Username == username {
			copy := *user
			return &copy, nil
		}
	}
	return nil, userrepository.ErrNotFound
}

func (f *fakeUsers) FindByID(_ context.Context, id uint64) (*models.User, error) {
	user, ok := f.users[id]
	if !ok {
		return nil, userrepository.ErrNotFound
	}
	copy := *user
	return &copy, nil
}

func (f *fakeUsers) UpdateLoginState(context.Context, uint64, int64, *string) error {
	return nil
}

// fakeSessions 用 map 模拟 Redis：删除旧 key 即代表令牌已失效。
type fakeSessions struct {
	tokens map[string]sessionrepository.Token
}

var _ sessionrepository.Repository = (*fakeSessions)(nil)

func newFakeSessions() *fakeSessions {
	return &fakeSessions{tokens: map[string]sessionrepository.Token{}}
}

func sessionKey(tokenType, hash string) string {
	return tokenType + ":" + hash
}

func (f *fakeSessions) Save(_ context.Context, token sessionrepository.Token) error {
	f.tokens[sessionKey(token.TokenType, token.TokenHash)] = token
	return nil
}

func (f *fakeSessions) FindByHash(
	_ context.Context,
	tokenHash string,
	tokenType string,
) (*sessionrepository.Token, error) {
	token, ok := f.tokens[sessionKey(tokenType, tokenHash)]
	if !ok {
		return nil, sessionrepository.ErrNotFound
	}
	token.TokenHash = tokenHash
	token.TokenType = tokenType
	return &token, nil
}

func (f *fakeSessions) Rotate(
	_ context.Context,
	currentHash string,
	tokenType string,
	next sessionrepository.Token,
) error {
	oldKey := sessionKey(tokenType, currentHash)
	if _, ok := f.tokens[oldKey]; !ok {
		return sessionrepository.ErrRefreshTokenRotated
	}
	delete(f.tokens, oldKey)
	f.tokens[sessionKey(next.TokenType, next.TokenHash)] = next
	return nil
}

func (f *fakeSessions) Revoke(_ context.Context, tokenHash, tokenType string) error {
	delete(f.tokens, sessionKey(tokenType, tokenHash))
	return nil
}

func newTestSession(
	users userrepository.UserRepository,
	sessions sessionrepository.Repository,
	now time.Time,
) *AdminSession {
	session := NewAdminSession(users, sessions, AdminLoginConfig{
		Secret:     []byte("0123456789abcdef0123456789abcdef"),
		Issuer:     "test-issuer",
		Audience:   auth.AudienceAdminAPI,
		TokenType:  auth.TokenTypeAdmin,
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 7 * 24 * time.Hour,
	})
	session.now = func() time.Time { return now }
	return session
}

func TestRefreshRotatesTokenAndIssuesAccessToken(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	users := newFakeUsers()
	users.users[7] = &models.User{ID: 7, Username: "admin", Name: "小阙", Status: 1}
	sessions := newFakeSessions()
	oldPlain := "old-refresh-token"
	oldHash := auth.HashRefreshToken(oldPlain)
	sessions.tokens[sessionKey(auth.TokenTypeAdmin, oldHash)] = sessionrepository.Token{
		UserID:    7,
		TokenHash: oldHash,
		TokenType: auth.TokenTypeAdmin,
		ExpiresAt: now.Add(time.Hour).Unix(),
	}
	session := newTestSession(users, sessions, now)

	result, err := session.Refresh(context.Background(), auth.RefreshCredentials{
		RefreshToken: oldPlain,
		UserAgent:    "test-agent",
		IP:           "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("Refresh() 返回错误：%v", err)
	}
	if result == nil || result.Pair == nil || result.User == nil {
		t.Fatal("Refresh() 结果不完整")
	}
	if result.User.ID != 7 {
		t.Fatalf("用户 ID = %d，期望 7", result.User.ID)
	}
	if result.Pair.AccessToken == "" || result.Pair.RefreshToken == "" {
		t.Fatal("应返回新的 Access Token 与刷新令牌")
	}
	if result.Pair.RefreshToken == oldPlain {
		t.Fatal("刷新令牌必须轮换")
	}
	if result.Pair.ExpiresIn != 900 {
		t.Fatalf("ExpiresIn = %d，期望 900", result.Pair.ExpiresIn)
	}

	claims, err := auth.VerifyAccessToken(
		[]byte("0123456789abcdef0123456789abcdef"),
		result.Pair.AccessToken,
		auth.TokenVerifyOptions{Issuer: "test-issuer", Audience: auth.AudienceAdminAPI, TokenType: auth.TokenTypeAdmin},
	)
	if err != nil {
		t.Fatalf("新 Access Token 校验失败：%v", err)
	}
	if claims.Subject != strconv.FormatUint(7, 10) || claims.Username != "admin" {
		t.Fatalf("claims 不正确：%#v", claims)
	}

	if _, ok := sessions.tokens[sessionKey(auth.TokenTypeAdmin, oldHash)]; ok {
		t.Fatal("旧刷新令牌应被删除")
	}
	rotated, ok := sessions.tokens[sessionKey(auth.TokenTypeAdmin, auth.HashRefreshToken(result.Pair.RefreshToken))]
	if !ok {
		t.Fatal("新刷新令牌未写入")
	}
	if rotated.UserID != 7 ||
		rotated.TokenType != auth.TokenTypeAdmin ||
		rotated.ExpiresAt != now.Add(7*24*time.Hour).Unix() {
		t.Fatalf("新刷新令牌记录不正确：%#v", rotated)
	}
	if rotated.UserAgent == nil || *rotated.UserAgent != "test-agent" ||
		rotated.IP == nil || *rotated.IP != "127.0.0.1" {
		t.Fatalf("设备信息未记录：%#v", rotated)
	}
}

func TestRefreshRejectsOldTokenAfterRotation(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	users := newFakeUsers()
	users.users[7] = &models.User{ID: 7, Username: "admin", Status: 1}
	sessions := newFakeSessions()
	oldPlain := "old-refresh-token"
	oldHash := auth.HashRefreshToken(oldPlain)
	sessions.tokens[sessionKey(auth.TokenTypeAdmin, oldHash)] = sessionrepository.Token{
		UserID: 7, TokenHash: oldHash, TokenType: auth.TokenTypeAdmin,
		ExpiresAt: now.Add(time.Hour).Unix(),
	}
	session := newTestSession(users, sessions, now)

	if _, err := session.Refresh(context.Background(), auth.RefreshCredentials{
		RefreshToken: oldPlain,
	}); err != nil {
		t.Fatalf("首次刷新应成功：%v", err)
	}
	if _, err := session.Refresh(context.Background(), auth.RefreshCredentials{
		RefreshToken: oldPlain,
	}); !errors.Is(err, auth.ErrInvalidRefreshToken) {
		t.Fatalf("旧令牌重放应返回 ErrInvalidRefreshToken，实际 %v", err)
	}
}

func TestRefreshRejectsInvalidTokens(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)

	cases := []struct {
		name  string
		setup func(*fakeSessions) string
	}{
		{
			name:  "missing token",
			setup: func(*fakeSessions) string { return "" },
		},
		{
			name:  "unknown token",
			setup: func(*fakeSessions) string { return "unknown" },
		},
		{
			name: "expired token",
			setup: func(sessions *fakeSessions) string {
				plain := "expired"
				hash := auth.HashRefreshToken(plain)
				sessions.tokens[sessionKey(auth.TokenTypeAdmin, hash)] = sessionrepository.Token{
					UserID: 7, TokenHash: hash, TokenType: auth.TokenTypeAdmin,
					ExpiresAt: now.Unix(),
				}
				return plain
			},
		},
		{
			name: "frontend token",
			setup: func(sessions *fakeSessions) string {
				plain := "frontend"
				hash := auth.HashRefreshToken(plain)
				sessions.tokens[sessionKey(auth.TokenTypeFrontend, hash)] = sessionrepository.Token{
					UserID: 7, TokenHash: hash, TokenType: auth.TokenTypeFrontend,
					ExpiresAt: now.Add(time.Hour).Unix(),
				}
				return plain
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			users := newFakeUsers()
			users.users[7] = &models.User{ID: 7, Username: "admin", Status: 1}
			sessions := newFakeSessions()
			plain := testCase.setup(sessions)
			session := newTestSession(users, sessions, now)

			if _, err := session.Refresh(context.Background(), auth.RefreshCredentials{
				RefreshToken: plain,
			}); !errors.Is(err, auth.ErrInvalidRefreshToken) {
				t.Fatalf("应返回 ErrInvalidRefreshToken，实际 %v", err)
			}
		})
	}
}

func TestRefreshRejectsDisabledUser(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	users := newFakeUsers()
	users.users[7] = &models.User{ID: 7, Username: "admin", Status: 0}
	sessions := newFakeSessions()
	plain := "refresh-token"
	hash := auth.HashRefreshToken(plain)
	sessions.tokens[sessionKey(auth.TokenTypeAdmin, hash)] = sessionrepository.Token{
		UserID: 7, TokenHash: hash, TokenType: auth.TokenTypeAdmin,
		ExpiresAt: now.Add(time.Hour).Unix(),
	}
	session := newTestSession(users, sessions, now)

	if _, err := session.Refresh(context.Background(), auth.RefreshCredentials{
		RefreshToken: plain,
	}); !errors.Is(err, auth.ErrAccountDisabled) {
		t.Fatalf("应返回 ErrAccountDisabled，实际 %v", err)
	}
}

func TestLogoutRevokesToken(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	users := newFakeUsers()
	sessions := newFakeSessions()
	plain := "refresh-token"
	hash := auth.HashRefreshToken(plain)
	sessions.tokens[sessionKey(auth.TokenTypeAdmin, hash)] = sessionrepository.Token{
		UserID: 7, TokenHash: hash, TokenType: auth.TokenTypeAdmin,
		ExpiresAt: now.Add(time.Hour).Unix(),
	}
	session := newTestSession(users, sessions, now)

	if err := session.Logout(context.Background(), auth.LogoutCredentials{RefreshToken: plain}); err != nil {
		t.Fatalf("Logout() 返回错误：%v", err)
	}
	if _, ok := sessions.tokens[sessionKey(auth.TokenTypeAdmin, hash)]; ok {
		t.Fatal("刷新令牌未删除")
	}
}

func TestLogoutIsIdempotent(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	users := newFakeUsers()
	sessions := newFakeSessions()
	session := newTestSession(users, sessions, now)

	for _, token := range []string{"", "unknown-token"} {
		if err := session.Logout(context.Background(), auth.LogoutCredentials{RefreshToken: token}); err != nil {
			t.Fatalf("Logout(%q) 应幂等成功，实际 %v", token, err)
		}
	}
}
