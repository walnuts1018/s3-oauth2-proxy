package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"github.com/walnuts1018/s3-oauth2-proxy/domain/model"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository/mock_repository"
)

func TestProxyUsecase_GetObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	key := "test_key"

	t.Run("successful get object", func(t *testing.T) {
		mockS3Repo := mock_repository.NewMockS3Repository(ctrl)
		expectedObject := &model.S3Object{
			Body:          io.NopCloser(bytes.NewReader([]byte("test_content"))),
			ContentLength: 12,
			ContentType:   "text/plain",
		}
		mockS3Repo.EXPECT().GetObject(ctx, key).Return(expectedObject, nil).Times(1)

		proxyUsecase := NewProxyUsecase(mockS3Repo)

		obj, err := proxyUsecase.GetObject(ctx, key)

		assert.NoError(t, err)
		assert.Equal(t, expectedObject, obj)
	})

	t.Run("get object failure", func(t *testing.T) {
		mockS3Repo := mock_repository.NewMockS3Repository(ctrl)
		expectedErr := errors.New("failed to get object")
		mockS3Repo.EXPECT().GetObject(ctx, key).Return((*model.S3Object)(nil), expectedErr).Times(1)

		proxyUsecase := NewProxyUsecase(mockS3Repo)

		obj, err := proxyUsecase.GetObject(ctx, key)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, obj)
	})
}
