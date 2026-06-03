package user

import userpb "example.com/nano_template/proto/user"

type UserModel interface {
	ToPbUserProfile(*UserProfile) *userpb.GetUserProfileResponse
}

type userModel struct{}

func NewUserModel() UserModel {
	return &userModel{}
}

func (*userModel) ToPbUserProfile(profile *UserProfile) *userpb.GetUserProfileResponse {
	return &userpb.GetUserProfileResponse{
		Id:        int32(profile.ID),
		UserId:    int32(profile.UserID),
		Avatar:    profile.Avatar,
		Nickname:  profile.Nickname,
		Email:     profile.Email,
		Phone:     profile.Phone,
		Signature: profile.Signature,
		CreatedAt: profile.CreatedAt.String(),
		UpdatedAt: profile.UpdatedAt.String(),
	}
}
