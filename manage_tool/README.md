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
- Air（開発用ホットリロード）

### セットアップ

1. **Air をインストール**
   ```bash
   go install github.com/air-verse/air@latest
   ```

2. **依存関係をインストール**
   ```bash
   cd manage_tool
   go mod download
   ```

3. **開発サーバー起動**
   ```bash
   air
   ```

### ディレクトリ構成

```
manage_tool/
├── main.go              # メインアプリケーション
├── go.mod               # Go modules設定
├── .air.toml            # Air設定ファイル
├── tmp/                 # Air一時ビルドファイル
├── templates/           # HTMLテンプレート
│   ├── template.go      # テンプレートレンダラー
│   ├── login.html       # ログイン画面
│   └── dashboard.html   # ダッシュボード
├── handlers/            # リクエストハンドラー
│   └── auth.go          # 認証ハンドラー
└── README.md            # このファイル
```