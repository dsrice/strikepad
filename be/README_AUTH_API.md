# 認証API仕様

## サインアップAPI

### エンドポイント
```
POST /api/auth/signup
```

### リクエスト
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe"
}
```

### 成功レスポンス (201 Created)
```json
{
  "id": 1,
  "email": "user@example.com",
  "display_name": "John Doe",
  "email_verified": false,
  "created_at": "2025-01-27T10:15:30Z"
}
```

### エラーレスポンス

#### バリデーションエラー (400 Bad Request)
```json
{
  "error": "validation_error",
  "message": "Email is required"
}
```

#### 既存ユーザー (409 Conflict)
```json
{
  "error": "user_exists",
  "message": "User with this email already exists"
}
```

## ログインAPI

### エンドポイント
```
POST /api/auth/login
```

### リクエスト
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### 成功レスポンス (200 OK)
```json
{
  "id": 1,
  "email": "user@example.com",
  "display_name": "John Doe",
  "email_verified": false
}
```

### エラーレスポンス

#### 認証失敗 (401 Unauthorized)
```json
{
  "error": "invalid_credentials",
  "message": "Invalid email or password"
}
```

## バリデーション仕様

### メールアドレス
- 必須フィールド
- 正規表現: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
- 自動的に小文字に変換・空白除去

### パスワード
- 必須フィールド
- 最小8文字、最大128文字
- bcryptでハッシュ化して保存

### 表示名
- 必須フィールド
- 最小1文字、最大100文字

## テスト用cURLコマンド

### サインアップ
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d @test_signup.json
```

### ログイン
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d @test_login.json
```

## セキュリティ機能

1. **パスワードハッシュ化**: bcryptを使用
2. **削除済みユーザーチェック**: is_deletedフラグで論理削除対応
3. **メール重複チェック**: 同一メールでの重複登録を防止
4. **入力値正規化**: メールアドレスの小文字変換・空白除去
5. **エラーハンドリング**: 詳細なエラーメッセージでデバッグ支援

## データベース

- プロバイダー種別: `email` (メールアドレス認証)
- 論理削除対応: `is_deleted`フラグ
- メール確認: `email_verified`フラグ（将来の機能拡張用）