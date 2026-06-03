package user

import (
	"errors"

	"example.com/nano_template/pkg/rpc/middleware"
	"example.com/nano_template/pkg/rpc/services/role"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserService interface {
	// 用户基本服务
	Register(username, password string, roleId uint) (*User, *role.Role, error)
	Login(username, password string, roleId uint) (*User, *role.Role, error)
	UpdatePassword(userId uint, oldPassword, newPassword string) error
	Cancel(userId uint) error
	// 用户首选项服务
	CreateProfile(userId uint, avatar, nickname, email, phone, signature *string) (*UserProfile, error)
	GetProfile(userId, id *uint) (*UserProfile, error)
	ListProfiles(page, pageSize *int, keywords *string) ([]*UserProfile, int64, error)
	UpdateProfile(userId, id *uint, avatar, nickname, email, phone, signature *string) (*UserProfile, error)
	DeleteProfile(userId, id *uint) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) Register(username, password string, roleId uint) (*User, *role.Role, error) {
	var existingUser User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, nil, status.Error(codes.AlreadyExists, "user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, status.Error(codes.Internal, "failed to check user")
	}

	var targetRole role.Role
	if err := s.db.First(&targetRole, roleId).Error; err != nil {
		return nil, nil, status.Error(codes.NotFound, "role not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to hash password")
	}

	user := User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		userRole := UserRole{
			UserID: user.ID,
			RoleID: roleId,
			State:  1,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to create user")
	}

	return &user, &targetRole, nil
}

func (s *userService) Login(username, password string, roleId uint) (*User, *role.Role, error) {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, status.Error(codes.NotFound, "username not found")
		}
		return nil, nil, status.Error(codes.Internal, "failed to fetch user")
	}

	// Verify user-role binding
	var userRole UserRole
	if err := s.db.Where("user_id = ? AND role_id = ?", user.ID, roleId).First(&userRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, status.Error(codes.PermissionDenied, "user does not have the specified role")
		}
		return nil, nil, status.Error(codes.Internal, "failed to fetch user role")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	// Fetch role details
	var targetRole role.Role
	if err := s.db.First(&targetRole, roleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, nil, status.Error(codes.Internal, "failed to fetch role details")
	}

	return &user, &targetRole, nil
}

func (s *userService) UpdatePassword(userId uint, oldPassword, newPassword string) error {
	var user User
	if err := s.db.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Error(codes.NotFound, "user not found")
		}
		return status.Error(codes.Internal, "failed to fetch user")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return status.Error(codes.Unauthenticated, "invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return status.Error(codes.Internal, "failed to hash password")
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return status.Error(codes.Internal, "failed to update password")
	}
	return nil
}

func (s *userService) Cancel(userId uint) error {
	if err := s.db.Delete(&User{}, userId).Error; err != nil {
		return status.Error(codes.Internal, "failed to cancel user")
	}
	return nil
}

/* 用户首选项服务 */

func (s *userService) CreateProfile(userId uint, avatar, nickname, email, phone, signature *string) (*UserProfile, error) {
	profile := UserProfile{UserID: userId}
	middleware.AssignString(avatar, &profile.Avatar)
	middleware.AssignString(nickname, &profile.Nickname)
	middleware.AssignString(email, &profile.Email)
	middleware.AssignString(phone, &profile.Phone)
	middleware.AssignString(signature, &profile.Signature)

	if err := s.db.Create(&profile).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to create user profile")
	}
	return &profile, nil
}

func (s *userService) GetProfile(userId, id *uint) (*UserProfile, error) {
	query, ok := s.profileQuery(userId, id)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "profile filter not found")
	}

	var profile UserProfile
	if err := query.First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user profile not found")
		}
		return nil, status.Error(codes.Internal, "failed to fetch user profile")
	}
	return &profile, nil
}

func (s *userService) ListProfiles(page, pageSize *int, keywords *string) ([]*UserProfile, int64, error) {
	query := s.db.Model(&UserProfile{})
	if keywords != nil {
		like := "%" + *keywords + "%"
		query = query.Where(
			"avatar LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ? OR signature LIKE ?",
			like, like, like, like, like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "failed to count user profiles")
	}

	pageValue := 1
	if page != nil {
		pageValue = *page
	}
	pageSizeValue := 10
	if pageSize != nil {
		pageSizeValue = *pageSize
	}
	if pageValue < 1 {
		pageValue = 1
	}

	var profiles []UserProfile
	if err := query.Offset((pageValue - 1) * pageSizeValue).Limit(pageSizeValue).Find(&profiles).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "failed to list user profiles")
	}

	result := make([]*UserProfile, 0, len(profiles))
	for i := range profiles {
		result = append(result, &profiles[i])
	}
	return result, total, nil
}

func (s *userService) UpdateProfile(userId, id *uint, avatar, nickname, email, phone, signature *string) (*UserProfile, error) {
	query, ok := s.profileQuery(userId, id)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "profile filter not found")
	}

	var profile UserProfile
	if err := query.First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user profile not found")
		}
		return nil, status.Error(codes.Internal, "failed to fetch user profile")
	}

	updates := map[string]any{}
	middleware.AssignUpdate(updates, "avatar", avatar)
	middleware.AssignUpdate(updates, "nickname", nickname)
	middleware.AssignUpdate(updates, "email", email)
	middleware.AssignUpdate(updates, "phone", phone)
	middleware.AssignUpdate(updates, "signature", signature)
	if len(updates) == 0 {
		return &profile, nil
	}

	if err := s.db.Model(&profile).Updates(updates).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to update user profile")
	}
	return &profile, nil
}

func (s *userService) DeleteProfile(userId, id *uint) error {
	query, ok := s.profileQuery(userId, id)
	if !ok {
		return status.Error(codes.InvalidArgument, "profile filter not found")
	}

	if err := query.Delete(&UserProfile{}).Error; err != nil {
		return status.Error(codes.Internal, "failed to delete user profile")
	}
	return nil
}

func (s *userService) profileQuery(userId, id *uint) (*gorm.DB, bool) {
	query := s.db.Model(&UserProfile{})
	hasFilter := false
	if userId != nil {
		query = query.Where("user_id = ?", *userId)
		hasFilter = true
	}
	if id != nil {
		query = query.Where("id = ?", *id)
		hasFilter = true
	}
	return query, hasFilter
}
