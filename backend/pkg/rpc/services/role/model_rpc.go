package role

import rolepb "example.com/nano_template/proto/role"

type RoleModel interface {
	ToPbRole(*Role) *rolepb.Role
}

type roleModel struct{}

func NewRoleModel() RoleModel {
	return &roleModel{}
}

func (*roleModel) ToPbRole(role *Role) *rolepb.Role {
	return &rolepb.Role{
		Id:          uint64(role.ID),
		Name:        role.Name,
		Description: role.Description,
		Level:       int32(role.Level),
		State:       uint32(role.State),
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
	}
}
