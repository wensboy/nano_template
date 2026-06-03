package role

import (
	"example.com/nano_template/pkg/web/middleware"
	rolepb "example.com/nano_template/proto/role"
)

type (
	CreateRequest struct {
		Name        string `json:"name" binding:"required"`
		Level       int32  `json:"level"`
		State       uint32 `json:"state"`
		Description string `json:"description"`
	}
	ListRequest struct {
		KeyWords *string `form:"keywords"`
		Level    *int32  `form:"level"`
		State    *uint32 `form:"state"`
	}
	UpdateRequest struct {
		Name        string `json:"name"`
		Level       int32  `json:"level"`
		State       uint32 `json:"state"`
		Description string `json:"description"`
	}
)

type (
	CreateResponse struct {
		ID uint64 `json:"id"`
	}
	GetResponse struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Level       int32  `json:"level"`
		State       uint32 `json:"state"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
	ListResponse struct {
		Roles []GetResponse `json:"roles"`
		Total int64         `json:"total"`
	}
)

type RoleModel interface {
	ToCreateResponse(*rolepb.CreateResponse) CreateResponse
	ToGetResponse(*rolepb.GetResponse) GetResponse
	ToListResponse(*rolepb.ListResponse) middleware.Pagination[GetResponse]
}

type roleModel struct{}

func NewRoleModel() RoleModel {
	return &roleModel{}
}

func (*roleModel) ToCreateResponse(resp *rolepb.CreateResponse) CreateResponse {
	return CreateResponse{
		ID: resp.Id,
	}
}

func (*roleModel) ToGetResponse(resp *rolepb.GetResponse) GetResponse {
	return GetResponse{
		ID:          resp.Role.Id,
		Name:        resp.Role.Name,
		Description: resp.Role.Description,
		Level:       resp.Role.Level,
		State:       resp.Role.State,
		CreatedAt:   resp.Role.CreatedAt,
		UpdatedAt:   resp.Role.UpdatedAt,
	}
}

func (*roleModel) ToListResponse(resp *rolepb.ListResponse) middleware.Pagination[GetResponse] {
	items := make([]GetResponse, 0, len(resp.Roles))
	for _, role := range resp.Roles {
		items = append(items, GetResponse{
			ID:          role.Id,
			Name:        role.Name,
			Description: role.Description,
			Level:       role.Level,
			State:       role.State,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return middleware.Pagination[GetResponse]{
		Total: int(resp.Total),
		Items: items,
	}
}
