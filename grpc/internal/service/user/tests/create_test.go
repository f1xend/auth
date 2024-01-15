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

func TestCreate(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.UserInfo
	}

	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		txMock = mocks.NewTxManagerMock(mc)

		id       = gofakeit.Int64()
		name     = gofakeit.FirstName()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 10)
		role     = gofakeit.Bool()

		repoErr = fmt.Errorf("repo error")

		req = &model.UserInfo{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     role,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "Success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
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

			newID, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}

	//for _, tt := range tests {
	//	tt := tt
	//	t.Run(tt.name, func(t *testing.T) {
	//		userRepositoryMock := tt.userRepositoryMock(mc)
	//		service := user.NewMockService(userRepositoryMock)
	//
	//		newID, err := service.Create(tt.args.ctx, tt.args.req)
	//		require.Equal(t, tt.err, err)
	//		require.Equal(t, tt.want, newID)
	//	})
	//}
}
