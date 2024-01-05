package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/f1xend/auth/internal/client/db/mocks"
	"github.com/f1xend/auth/internal/model"
	"github.com/f1xend/auth/internal/repository"
	repoMocks "github.com/f1xend/auth/internal/repository/mocks"
	"github.com/f1xend/auth/internal/service/user"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdate(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.UpdateUser
	}

	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		txMock = mocks.NewTxManagerMock(mc)

		id    = gofakeit.Int64()
		name  = gofakeit.FirstName()
		email = gofakeit.Email()

		repoErr = fmt.Errorf("repo error")

		info = model.UpdateUserInfo{
			Name:  name,
			Email: email,
		}
		req = &model.UpdateUser{
			ID:   id,
			Info: info,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "Success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			service := user.NewService(userRepoMock, txMock)

			err := service.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
