package role

import (
	"errors"

	"example.com/nano_template/pkg/rpc/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RoleService interface {
	Create(name string, level int, state uint, description *string) (*Role, error)
	Get(id uint) (*Role, error)
	List(page, pageSize *int, keywords *string, level *int, state *uint) ([]Role, int64, error)
	Update(id uint, name, description *string, level, state *uint) error
	Delete(id uint) error
}

type roleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) RoleService {
	return &roleService{db: db}
}

func (s *roleService) Create(name string, level int, state uint, description *string) (*Role, error) {
	var existingRole Role
	if err := s.db.Where("name = ?", name).First(&existingRole).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "role already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.Internal, "failed to check role")
	}

	role := Role{
		Name:  name,
		Level: level,
		State: state,
	}
	if description != nil {
		role.Description = *description
	}
	if err := s.db.Create(&role).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to create role")
	}
	return &role, nil
}

func (s *roleService) Get(id uint) (*Role, error) {
	var role Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, status.Error(codes.Internal, "failed to fetch role")
	}
	return &role, nil
}

func (s *roleService) List(page, pageSize *int, keywords *string, level *int, state *uint) ([]Role, int64, error) {
	var roles []Role
	var total int64

	query := s.db.Model(&Role{})
	if keywords != nil {
		like := "%" + *keywords + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if level != nil {
		query = query.Where("level = ?", *level)
	}
	if state != nil {
		query = query.Where("state = ?", *state)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "failed to count roles")
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

	if err := query.Offset((pageValue - 1) * pageSizeValue).Limit(pageSizeValue).Find(&roles).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "failed to list roles")
	}

	return roles, total, nil
}

func (s *roleService) Update(id uint, name, description *string, level, state *uint) error {
	var role Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Error(codes.NotFound, "role not found")
		}
		return status.Error(codes.Internal, "failed to fetch role")
	}

	updates := map[string]any{}
	middleware.AssignUpdate(updates, "name", name)
	middleware.AssignUpdate(updates, "description", description)
	middleware.AssignUpdate(updates, "level", level)
	middleware.AssignUpdate(updates, "state", state)
	if len(updates) == 0 {
		return nil
	}

	if err := s.db.Model(&role).Updates(updates).Error; err != nil {
		return status.Error(codes.Internal, "failed to update role")
	}
	return nil
}

func (s *roleService) Delete(id uint) error {
	var role Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Error(codes.NotFound, "role not found")
		}
		return status.Error(codes.Internal, "failed to fetch role")
	}
	if err := s.db.Delete(&Role{}, id).Error; err != nil {
		return status.Error(codes.Internal, "failed to delete role")
	}
	return nil
}
