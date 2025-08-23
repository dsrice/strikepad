package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient はMinIOクライアントのラッパー
type MinIOClient struct {
	Client     *minio.Client
	BucketName string
}

// NewMinIOClient は新しいMinIOクライアントを作成
func NewMinIOClient() (*MinIOClient, error) {
	endpoint := getMinIOEnv("MINIO_ENDPOINT", "localhost:9000")
	accessKeyID := getMinIOEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretAccessKey := getMinIOEnv("MINIO_SECRET_KEY", "minioadmin")
	bucketName := getMinIOEnv("MINIO_BUCKET", "dev-strikepad")
	useSSL := false

	// MinIOクライアントを初期化
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("MinIOクライアントの初期化に失敗しました: %w", err)
	}

	client := &MinIOClient{
		Client:     minioClient,
		BucketName: bucketName,
	}

	// バケットを作成（存在しない場合）
	err = client.EnsureBucket()
	if err != nil {
		return nil, fmt.Errorf("バケットの作成に失敗しました: %w", err)
	}

	return client, nil
}

// EnsureBucket はバケットが存在することを確認し、存在しない場合は作成する
func (m *MinIOClient) EnsureBucket() error {
	ctx := context.Background()
	exists, err := m.Client.BucketExists(ctx, m.BucketName)
	if err != nil {
		return fmt.Errorf("バケット存在確認に失敗しました: %w", err)
	}

	if !exists {
		err = m.Client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("バケット作成に失敗しました: %w", err)
		}
		log.Printf("バケット %s を作成しました", m.BucketName)
	}

	return nil
}

// getMinIOEnv は環境変数を取得し、設定されていない場合はデフォルト値を返す
func getMinIOEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}