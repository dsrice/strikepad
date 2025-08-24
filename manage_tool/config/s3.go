package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3ClientInterface はS3クライアントのインターフェース
type S3ClientInterface interface {
	PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadBucket(ctx context.Context, input *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	CreateBucket(ctx context.Context, input *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	PutBucketPolicy(ctx context.Context, input *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
}

// S3Client はAWS S3クライアントのラッパー
type S3Client struct {
	Client        S3ClientInterface
	PreSignClient *s3.PresignClient
	BucketName    string
	Region        string
}

// NewS3Client は新しいS3クライアントを作成
func NewS3Client() (*S3Client, error) {
	bucketName := getS3Env("S3_BUCKET", "dev-strikepad")
	region := getS3Env("AWS_REGION", "us-east-1")
	endpoint := getS3Env("S3_ENDPOINT", "http://localhost:9000") // MinIO互換性のため
	accessKeyID := getS3Env("AWS_ACCESS_KEY_ID", "minioadmin")
	secretAccessKey := getS3Env("AWS_SECRET_ACCESS_KEY", "minioadmin")

	var cfg aws.Config
	var err error

	if endpoint != "" {
		// MinIO互換エンドポイントを使用（開発環境用）
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
			config.WithBaseEndpoint(endpoint),
		)
	} else {
		// 通常のAWS S3を使用（本番環境用）
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %w", err)
	}

	// Path-style addressingを有効にする（MinIO互換性のため）
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.UsePathStyle = true
		}
	})

	// 署名付きURL生成用のPresignClientを作成
	presignClient := s3.NewPresignClient(client)

	s3Client := &S3Client{
		Client:        client,
		PreSignClient: presignClient,
		BucketName:    bucketName,
		Region:        region,
	}

	// バケットを作成（存在しない場合）
	err = s3Client.EnsureBucket()
	if err != nil {
		return nil, fmt.Errorf("バケットの作成に失敗しました: %w", err)
	}

	return s3Client, nil
}

// EnsureBucket はバケットが存在することを確認し、存在しない場合は作成する
func (s *S3Client) EnsureBucket() error {
	ctx := context.Background()

	// バケットの存在確認
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.BucketName),
	})

	if err != nil {
		// バケットが存在しない場合は作成
		log.Printf("バケット %s が存在しないため作成します", s.BucketName)

		createBucketInput := &s3.CreateBucketInput{
			Bucket: aws.String(s.BucketName),
		}

		// us-east-1以外の場合はLocationConstraintが必要
		if s.Region != "us-east-1" {
			createBucketInput.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(s.Region),
			}
		}

		_, err = s.Client.CreateBucket(ctx, createBucketInput)
		if err != nil {
			return fmt.Errorf("バケット作成に失敗しました: %w", err)
		}
		log.Printf("バケット %s を作成しました", s.BucketName)

		// パブリック読み取りポリシーを設定
		err = s.setPublicReadPolicy()
		if err != nil {
			log.Printf("バケットポリシーの設定に失敗しました: %v", err)
		} else {
			log.Printf("バケット %s にパブリック読み取りポリシーを設定しました", s.BucketName)
		}
	}

	return nil
}

// setPublicReadPolicy はバケットにパブリック読み取りポリシーを設定
func (s *S3Client) setPublicReadPolicy() error {
	ctx := context.Background()

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::%s/*"
			}
		]
	}`, s.BucketName)

	_, err := s.Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(s.BucketName),
		Policy: aws.String(policy),
	})

	return err
}

// GeneratePresignedURL は署名付きURLを生成する
func (s *S3Client) GeneratePresignedURL(objectKey string, expires time.Duration) (string, error) {
	ctx := context.Background()

	request, err := s.PreSignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		return "", fmt.Errorf("署名付きURL生成に失敗しました: %w", err)
	}

	return request.URL, nil
}

// GenerateUploadPresignedURL はアップロード用署名付きURLを生成する
func (s *S3Client) GenerateUploadPresignedURL(objectKey string, contentType string, expires time.Duration) (string, error) {
	ctx := context.Background()

	request, err := s.PreSignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		return "", fmt.Errorf("アップロード用署名付きURL生成に失敗しました: %w", err)
	}

	return request.URL, nil
}

// getS3Env は環境変数を取得し、設定されていない場合はデフォルト値を返す
func getS3Env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}