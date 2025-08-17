# Strikepad Management Tool

Strikepadの管理ツールです。Golang + Echoフレームワークで構築されています。

## 機能

- Strikepadの各種管理機能
- フロントエンド（FE）の静的ファイル配信
- 管理用API提供

## 起動方法

### 通常起動

```bash
cd manage_tool
go run main.go
```

### 開発用（ホットリロード）

```bash
cd manage_tool
air
```

サーバーはポート8080で起動します。

**Air使用時の利点:**

- ファイル変更時の自動リロード
- `.go`, `.html`ファイルの変更を監視
- 開発効率の向上

## エンドポイント

- `GET /api/health` - ヘルスチェック
- `/*` - 静的ファイル配信（FEのdistフォルダ）

## 開発

### 必要なもの

- Go 1.24.1以上
- Echo v4フレームワーク
- GORM v1.30.1以上
- PostgreSQL
- Air（開発用ホットリロード）
- golangci-lint（コード品質チェック）

### セットアップ

1. **開発ツールをインストール**
   ```bash
   make setup
   ```
   または個別にインストール：
   ```bash
   go install github.com/air-verse/air@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **依存関係をインストール**
   ```bash
   cd manage_tool
   make deps
   ```

3. **開発サーバー起動**
   ```bash
   make dev
   # または
   air
   ```

## コード品質管理

### Lint実行

```bash
# 基本的なLint実行
make lint

# 自動修正可能な問題を修正
make lint-fix

# 詳細な出力でLint実行
make lint-verbose

# CI用のJSON形式レポート生成
make lint-ci
```

### その他のコマンド

```bash
# コードフォーマット
make fmt

# セキュリティチェック
make security

# 全チェック実行（lint + security + build）
make check

# ビルド
make build

# クリーンアップ
make clean
```

### ディレクトリ構成

```
manage_tool/
├── main.go              # メインアプリケーション
├── go.mod               # Go modules設定
├── .air.toml            # Air設定ファイル
├── .golangci.yml        # golangci-lint設定
├── Makefile             # 開発用コマンド
├── tmp/                 # Air一時ビルドファイル
├── config/              # 設定関連
│   └── database.go      # データベース設定
├── models/              # データモデル
│   └── shared.go        # 共有モデル定義
├── repository/          # データアクセス層
│   └── user.go          # ユーザーリポジトリ
├── templates/           # HTMLテンプレート
│   ├── template.go      # テンプレートレンダラー
│   ├── login.html       # ログイン画面
│   ├── dashboard.html   # ダッシュボード
│   └── users.html       # ユーザー管理画面
├── handlers/            # リクエストハンドラー
│   ├── auth.go          # 認証ハンドラー
│   └── admin.go         # 管理機能ハンドラー
└── README.md            # このファイル
```