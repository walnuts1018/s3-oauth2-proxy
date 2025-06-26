package repository

import (
	"context"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/model"
)

type S3Repository interface {
	GetObject(ctx context.Context, key string) (*model.S3Object, error)
}
