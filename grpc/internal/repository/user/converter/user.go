package converter

import (
	"github.com/f1xend/auth/internal/repository/user/model"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromRepo(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromRepo(info *model.Info) *desc.UserInfo {
	var roleAdmin desc.Role
	if info.Role == true {
		roleAdmin = desc.Role_admin
	} else {
		roleAdmin = desc.Role_user
	}
	return &desc.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.Password,
		Role:            roleAdmin,
	}
}
