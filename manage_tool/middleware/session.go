package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// セッションの設定
const (
	SessionName          = "session" // デフォルト名を使用
	SessionMaxAge        = 15 * 60   // 15分（秒単位）
	SessionKeyUserID     = "user_id"
	SessionKeyUserName   = "user_name"
	SessionKeyLoginID    = "login_id"
	SessionKeyLastAccess = "last_access"
)

// CreateSessionStore はセッションストアを作成
func CreateSessionStore(secretKey string) sessions.Store {
	store := sessions.NewCookieStore([]byte(secretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionMaxAge,
		HttpOnly: true,
		Secure:   false, // 開発環境ではfalse、本番環境ではtrue
		SameSite: http.SameSiteLaxMode,
	}
	return store
}

// SessionAuth はセッション認証ミドルウェア
func SessionAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ログインページとログイン処理はスキップ
			path := c.Path()
			if path == "/login" || path == "/auth/login" {
				return next(c)
			}

			sess, err := session.Get(SessionName, c)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			// セッションが存在するかチェック
			userID, ok := sess.Values[SessionKeyUserID]
			if !ok || userID == nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			// 最終アクセス時間をチェック（15分タイムアウト）
			lastAccessInterface, ok := sess.Values[SessionKeyLastAccess]
			if !ok {
				return c.Redirect(http.StatusFound, "/login")
			}

			lastAccessUnix, ok := lastAccessInterface.(int64)
			if !ok {
				return c.Redirect(http.StatusFound, "/login")
			}

			lastAccess := time.Unix(lastAccessUnix, 0)

			// 15分以上経過している場合はセッションを削除
			if time.Since(lastAccess) > time.Duration(SessionMaxAge)*time.Second {
				sess.Values = make(map[interface{}]interface{})
				sess.Save(c.Request(), c.Response())
				return c.Redirect(http.StatusFound, "/login")
			}

			// 最終アクセス時間を更新
			sess.Values[SessionKeyLastAccess] = time.Now().Unix()
			sess.Save(c.Request(), c.Response())

			// ユーザー情報をコンテキストに設定
			c.Set("user_id", userID)
			c.Set("user_name", sess.Values[SessionKeyUserName])
			c.Set("login_id", sess.Values[SessionKeyLoginID])

			return next(c)
		}
	}
}