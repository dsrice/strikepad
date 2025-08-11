# エラーコード仕様

このAPIでは統一されたエラーコードシステム（E000形式）を使用して、一貫性のあるエラーレスポンスを提供します。

## エラーレスポンス形式

### 統一エラーレスポンス
全てのエラーは同じ構造体を使用し、必要に応じて詳細情報を含みます。

#### 標準エラーレスポンス
```json
{
  "code": "E102",
  "message": "User already exists",
  "description": "A user with this email address already exists"
}
```

#### バリデーションエラーレスポンス
```json
{
  "code": "E003",
  "message": "Validation failed",
  "description": "One or more fields failed validation",
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
    }
  ]
}
```

## エラーコード一覧

### 一般的なエラーコード (E001-E099)

| コード | HTTPステータス | メッセージ | 説明 |
|--------|---------------|-----------|------|
| `E001` | 500 | Internal server error | サーバー内部エラー |
| `E002` | 400 | Invalid request | リクエスト形式が無効 |
| `E003` | 400 | Validation failed | バリデーション失敗 |
| `E004` | 404 | Resource not found | リソースが見つからない |
| `E005` | 401 | Unauthorized | 認証が必要 |
| `E006` | 403 | Forbidden | アクセス権限なし |
| `E007` | 409 | Conflict | リソースの競合 |

### 認証関連のエラーコード (E100-E199)

| コード | HTTPステータス | メッセージ | 説明 |
|--------|---------------|-----------|------|
| `E100` | 401 | Invalid credentials | メールアドレスまたはパスワードが間違っている |
| `E101` | 404 | User not found | 指定されたメールアドレスのユーザーが見つからない |
| `E102` | 409 | User already exists | 同じメールアドレスのユーザーが既に存在 |
| `E103` | 401 | Token expired | 認証トークンの有効期限切れ |
| `E104` | 401 | Invalid token | 認証トークンが無効 |

### バリデーション関連のエラーコード (E200-E299)

| コード | HTTPステータス | メッセージ | 説明 |
|--------|---------------|-----------|------|
| `E200` | 400 | Email is required | メールアドレスは必須 |
| `E201` | 400 | Invalid email format | メールアドレスの形式が無効 |
| `E202` | 400 | Password is required | パスワードは必須 |
| `E203` | 400 | Password too short | パスワードが短すぎる（8文字未満） |
| `E204` | 400 | Password too long | パスワードが長すぎる（128文字超過） |
| `E205` | 400 | Display name is required | 表示名は必須 |
| `E206` | 400 | Display name too long | 表示名が長すぎる（100文字超過） |

### ビジネスロジック関連のエラーコード (E300-E399)

| コード | HTTPステータス | メッセージ | 説明 |
|--------|---------------|-----------|------|
| `E300` | 403 | Email not verified | メールアドレスが未確認 |
| `E301` | 403 | Account disabled | アカウントが無効化されている |
| `E302` | 403 | Account deleted | アカウントが削除されている |

## API別エラー例

### サインアップAPI (`POST /api/auth/signup`)

#### 成功 (201 Created)
```json
{
  "id": 1,
  "email": "user@example.com",
  "display_name": "John Doe",
  "email_verified": false,
  "created_at": "2025-01-27T10:15:30Z"
}
```

#### バリデーションエラー (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid","password":"123","display_name":""}'
```

```json
{
  "code": "VALIDATION_FAILED",
  "message": "Validation failed",
  "description": "One or more fields failed validation",
  "details": [
    {
      "field": "email",
      "tag": "email",
      "value": "invalid",
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

#### ユーザー重複エラー (409 Conflict)
```json
{
  "code": "E102",
  "message": "User already exists",
  "description": "A user with this email address already exists"
}
```

### ログインAPI (`POST /api/auth/login`)

#### 成功 (200 OK)
```json
{
  "id": 1,
  "email": "user@example.com",
  "display_name": "John Doe",
  "email_verified": false
}
```

#### 認証エラー (401 Unauthorized)
```json
{
  "code": "E100",
  "message": "Invalid credentials",
  "description": "The provided email or password is incorrect"
}
```

#### バリデーションエラー (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"","password":""}'
```

```json
{
  "code": "E003",
  "message": "Validation failed",
  "description": "One or more fields failed validation",
  "details": [
    {
      "field": "email",
      "tag": "required",
      "value": "",
      "message": "email is required"
    },
    {
      "field": "password",
      "tag": "required",
      "value": "",
      "message": "password is required"
    }
  ]
}
```

## フロントエンド実装ガイド

### エラーハンドリング例

```typescript
interface ErrorResponse {
  code: string;
  message: string;
  description?: string;
}

interface ValidationErrorResponse extends ErrorResponse {
  details: ValidationError[];
}

interface ValidationError {
  field: string;
  tag: string;
  value: string;
  message: string;
}

// エラーハンドリング関数
function handleApiError(error: any) {
  if (error.response?.data?.code) {
    const errorCode = error.response.data.code;
    
    switch (errorCode) {
      case 'E102':
        showMessage('このメールアドレスは既に使用されています');
        break;
      case 'E100':
        showMessage('メールアドレスまたはパスワードが間違っています');
        break;
      case 'E003':
        handleValidationErrors(error.response.data.details);
        break;
      default:
        showMessage('エラーが発生しました');
    }
  }
}
```

## 実装詳細

### エラーコード定義
```go
type ErrorCode string

const (
    ErrCodeUserExists    ErrorCode = "E102"
    ErrCodeEmailInvalid  ErrorCode = "E201"
    // ...
)
```

### エラー情報取得
```go
func GetErrorInfo(code ErrorCode) ErrorInfo {
    // エラーコードに対応する詳細情報を返す
}
```

### ハンドラでの使用
```go
errorInfo := errors.GetErrorInfo(errors.ErrCodeUserExists)
return c.JSON(http.StatusConflict, dto.ErrorResponse{
    Code:        string(errorInfo.Code),
    Message:     errorInfo.Message,
    Description: errorInfo.Description,
})
```