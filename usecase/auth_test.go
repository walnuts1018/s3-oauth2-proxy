package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository/mock_repository"
)

func TestAuthUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	code := "test_code"
	expectedNonce := "test_nonce"
	expectedSubject := "test_subject"

	t.Run("successful login", func(t *testing.T) {
		mockAuthRepo := mock_repository.NewMockAuthRepository(ctrl)
		token := &oauth2.Token{}
		mockAuthRepo.EXPECT().Exchange(ctx, code).Return(token, nil).Times(1)
		mockAuthRepo.EXPECT().VerifyIDToken(ctx, token, expectedNonce).Return(expectedSubject, nil).Times(1)

		authUsecase := NewAuthUsecase(mockAuthRepo)

		subject, err := authUsecase.Login(ctx, code, expectedNonce)

		assert.NoError(t, err)
		assert.Equal(t, expectedSubject, subject)
	})

	t.Run("exchange token failure", func(t *testing.T) {
		mockAuthRepo := mock_repository.NewMockAuthRepository(ctrl)
		expectedErr := errors.New("failed to exchange token")
		mockAuthRepo.EXPECT().Exchange(ctx, code).Return((*oauth2.Token)(nil), expectedErr).Times(1)

		authUsecase := NewAuthUsecase(mockAuthRepo)

		subject, err := authUsecase.Login(ctx, code, expectedNonce)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Empty(t, subject)
	})

	t.Run("verify id token failure", func(t *testing.T) {
		mockAuthRepo := mock_repository.NewMockAuthRepository(ctrl)
		token := &oauth2.Token{}
		expectedErr := errors.New("failed to verify id token")
		mockAuthRepo.EXPECT().Exchange(ctx, code).Return(token, nil).Times(1)
		mockAuthRepo.EXPECT().VerifyIDToken(ctx, token, expectedNonce).Return("", expectedErr).Times(1)

		authUsecase := NewAuthUsecase(mockAuthRepo)

		subject, err := authUsecase.Login(ctx, code, expectedNonce)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		assert.Empty(t, subject)
	})
}
