package usecase

import (
	"context"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/model"
	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository"
)

type ProxyUsecase interface {
	GetObject(ctx context.Context, key string) (*model.S3Object, error)
}

type proxyUsecase struct {
	s3Repo repository.S3Repository
}

func NewProxyUsecase(s3Repo repository.S3Repository) ProxyUsecase {
	return &proxyUsecase{s3Repo: s3Repo}
}

func (u *proxyUsecase) GetObject(ctx context.Context, key string) (*model.S3Object, error) {
	return u.s3Repo.GetObject(ctx, key)
}
