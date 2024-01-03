package converter

import (
	"github.com/f1xend/auth/internal/model"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

//func ToUserFromDesc(user model.User) *desc.User {
//	return &desc.User{
//		Id:        user.ID,
//		Info:      ToUserInfoFromService(user.Info),
//		CreatedAt: timestamppb.New(user.CreatedAt),
//		UpdatedAt: updatedAt,
//	}
//}

func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
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

func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	var role bool
	if info.Role == desc.Role_admin {
		role = true
	}

	return &model.UserInfo{
		Name:     info.Name,
		Email:    info.Email,
		Password: info.Password,
		Role:     role,
	}
}

func ToUserUpdateFromDesc(user *desc.UpdateRequest) *model.UpdateUser {
	return &model.UpdateUser{
		ID:   user.Id,
		Info: *ToUserUpdateInfoFromDesc(*user.Info),
	}
}

func ToUserUpdateInfoFromDesc(user desc.UpdateUserInfo) *model.UpdateUserInfo {
	log.Println(user.GetName().GetValue(), user.GetEmail().GetValue())

	return &model.UpdateUserInfo{

		Name:  user.Name.GetValue(),
		Email: user.Email.GetValue(),
	}
}
