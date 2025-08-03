# バリデーション仕様

このプロジェクトでは `go-playground/validator/v10` を使用して構造化されたリクエストバリデーションを実装しています。

## バリデーション機能

### 1. 構造体タグによるバリデーション

DTOにバリデーションタグを設定し、自動的にバリデーションを実行します。

```go
type SignupRequest struct {
    Email       string `json:"email" validate:"required,email,max=255"`
    Password    string `json:"password" validate:"required,min=8,max=128"`
    DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}
```

### 2. サポートしているバリデーションタグ

- `required`: 必須フィールド
- `email`: 有効なメールアドレス形式
- `min=N`: 最小文字数（文字列）/最小値（数値）
- `max=N`: 最大文字数（文字列）/最大値（数値）
- `len=N`: 正確な文字数
- `gt=N`: より大きい
- `gte=N`: 以上
- `lt=N`: より小さい
- `lte=N`: 以下
- `oneof=value1 value2`: 指定された値のいずれか
- `alpha`: アルファベットのみ
- `alphanum`: 英数字のみ
- `numeric`: 数値のみ
- `url`: 有効なURL
- `uri`: 有効なURI

### 3. バリデーションエラーレスポンス

#### 詳細エラー形式
```json
{
  "error": "validation_error",
  "message": "Validation failed",
  "details": [
    {
      "field": "email",
      "tag": "email",
      "value": "invalid-email",
      "message": "email must be a valid email address"
    },
    {
      "field": "password",
      "tag": "min",
      "value": "123",
      "message": "password must be at least 8 characters long"
    },
    {
      "field": "display_name",
      "tag": "required",
      "value": "",
      "message": "display_name is required"
    }
  ]
}
```

### 4. サインアップリクエストのバリデーション

#### 成功例
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe"
}
```

#### 失敗例とエラーメッセージ

**無効なメールアドレス**
```json
{
  "email": "invalid-email",
  "password": "password123",
  "display_name": "John Doe"
}
```
→ `email must be a valid email address`

**短いパスワード**
```json
{
  "email": "user@example.com",
  "password": "123",
  "display_name": "John Doe"
}
```
→ `password must be at least 8 characters long`

**長いパスワード**
```json
{
  "email": "user@example.com",
  "password": "a".repeat(129),
  "display_name": "John Doe"
}
```
→ `password must be at most 128 characters long`

**空の表示名**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": ""
}
```
→ `display_name is required`

### 5. ログインリクエストのバリデーション

#### 成功例
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### 失敗例
```json
{
  "email": "",
  "password": ""
}
```
→ 複数のエラーが返される

### 6. カスタムバリデーションメッセージ

バリデーターは日本語対応の分かりやすいエラーメッセージを生成します：

- `required`: `{field} is required`
- `email`: `{field} must be a valid email address`
- `min`: `{field} must be at least {param} characters long`
- `max`: `{field} must be at most {param} characters long`

### 7. テスト用cURLコマンド

#### バリデーションエラーのテスト
```bash
# 無効な入力でサインアップ
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d @test_validation_errors.json

# 空の入力でログイン
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d @test_login_invalid.json
```

## 実装詳細

### 1. バリデーターの初期化
```go
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        validator:   validator.New(),
    }
}
```

### 2. バリデーションの実行
```go
if err := h.validator.Validate(&req); err != nil {
    if ve, ok := err.(validator.ValidationErrors); ok {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "error":   "validation_error",
            "message": "Validation failed",
            "details": ve.Errors,
        })
    }
}
```

### 3. JSONフィールド名の使用
バリデーターは構造体フィールド名ではなく、JSONタグで指定された名前を使用してエラーメッセージを生成します。

これにより、フロントエンドから送信されるJSONフィールド名と一致したエラーメッセージが返されます。