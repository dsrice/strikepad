# ログ設定とローテーション

このプロジェクトでは構造化ログ（slog）を使用し、logsフォルダに1時間単位でログローテーションを行います。

## ログ設定

### 環境変数
- `LOG_LEVEL`: ログレベル（DEBUG, INFO, WARN, ERROR）
- `APP_ENV`: 環境設定（dev, production）

### 出力先
- **開発環境**: ファイル + コンソール両方に出力
- **本番環境**: ファイルのみに出力

### ログファイル設定
- **ファイル名**: `logs/app.log`
- **ローテーション**: 1時間単位で自動実行
- **最大ファイルサイズ**: 100MB
- **保持期間**: 24時間分のバックアップ（24ファイル）
- **保存期間**: 7日間
- **圧縮**: 古いログファイルは自動圧縮（gzip）

## ログ形式

### 開発環境（TextHandler）
```
time=2025-01-27T10:15:30.123+09:00 level=INFO source=main.go:106 msg="Initializing migration runner" environment=dev
time=2025-01-27T10:15:30.456+09:00 level=INFO source=runner.go:59 msg="Successfully applied migrations" count=2
```

### 本番環境（JSONHandler）
```json
{"time":"2025-01-27T10:15:30.123+09:00","level":"INFO","source":{"function":"main.runMigrations","file":"main.go","line":106},"msg":"Initializing migration runner","environment":"dev"}
{"time":"2025-01-27T10:15:30.456+09:00","level":"INFO","source":{"function":"github.com/example/internal/migrations.(*MigrationRunner).RunMigrations","file":"runner.go","line":59},"msg":"Successfully applied migrations","count":2}
```

## ファイル構成

```
be/
├── logs/
│   ├── .gitkeep           # ディレクトリ構造保持用
│   ├── app.log            # 現在のログファイル
│   ├── app.log.2025012709 # 1時間前のログ（圧縮前）
│   └── app.log.2025012708.gz # 2時間前のログ（圧縮済み）
└── ...
```

## 使用方法

### 基本的なログ出力
```go
import "log/slog"

// 情報ログ
slog.Info("User logged in", "user_id", 123, "ip", "192.168.1.1")

// エラーログ
slog.Error("Database connection failed", "error", err, "database", "users")

// デバッグログ
slog.Debug("Processing request", "request_id", "abc123", "path", "/api/users")

// 警告ログ
slog.Warn("Rate limit approaching", "user_id", 456, "requests", 95, "limit", 100)
```

### ログレベル設定
```bash
# デバッグレベル（開発時）
LOG_LEVEL=DEBUG go run main.go

# 本番環境（情報レベル）
APP_ENV=production LOG_LEVEL=INFO go run main.go
```

## 特徴

1. **構造化ログ**: キー・値ペアでログ情報を構造化
2. **自動ローテーション**: 1時間単位で自動的にログファイルを分割
3. **圧縮保存**: 古いログファイルは自動的にgzip圧縮
4. **ソース情報**: ログ出力箇所のファイル名・行番号を自動記録
5. **環境別設定**: 開発・本番環境に応じた適切な出力形式

## 注意事項

- ログファイルは `.gitignore` で除外されています
- `logs/.gitkeep` でディレクトリ構造のみGit管理されています
- アプリケーション起動時に自動的に `logs` ディレクトリが作成されます
- ログローテーションは別goroutineで実行されるため、メインプロセスに影響しません