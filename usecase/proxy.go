package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/model"
	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository"
)

var (
	ErrObjectNotFound = errors.New("object not found in S3")
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
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}
	obj, err := u.s3Repo.GetObject(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}
	return obj, nil
}
