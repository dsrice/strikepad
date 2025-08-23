package mocks

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
)

// MockMinIOClient は MinIOClient のモック
type MockMinIOClient struct {
	mock.Mock
	BucketName string
}

// MockObject は minio.Object のモック
type MockObject struct {
	mock.Mock
	io.ReadCloser
}

func (m *MockObject) Stat() (minio.ObjectInfo, error) {
	args := m.Called()
	return args.Get(0).(minio.ObjectInfo), args.Error(1)
}

func (m *MockObject) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockObject) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// MockMinioClient は minio.Client のモック
type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*minio.Object), args.Error(1)
}