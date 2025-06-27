package repository

import (
	"context"
	"errors"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/model"
)

var (
	ErrObjectNotFound = errors.New("object not found in S3")
)

//go:generate go tool mockgen -source=s3.go -destination=mock_repository/mock_s3_repository.go -package=mock_repository
type S3Repository interface {
	GetObject(ctx context.Context, key string) (*model.S3Object, error)
}
