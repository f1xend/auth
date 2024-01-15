package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/f1xend/auth/internal/api/user"
	"github.com/f1xend/auth/internal/converter"
	"github.com/f1xend/auth/internal/service"
	serviceMocks "github.com/f1xend/auth/internal/service/mocks"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.FirstName()
		email = gofakeit.Email()

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Id: id,
			Info: &desc.UpdateUserInfo{
				Name:  &wrappers.StringValue{Value: name},
				Email: &wrappers.StringValue{Value: email},
			},
		}

		res = &empty.Empty{}

		usr = converter.ToUserUpdateFromDesc(req)
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		err             error
		want            *empty.Empty
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err:  nil,
			want: res,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, usr).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, usr).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)

			resHandler, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
