package role

import (
	"math"
	"net/http"

	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/web/middleware"
	rolepb "example.com/nano_template/proto/role"
	"github.com/gin-gonic/gin"
)

// RoleHandler defines the interface for role CRUD operations.
type RoleHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type roleHandler struct {
	client rolepb.RoleServiceClient
	model  RoleModel
}

type RoleHandlerOption func(*roleHandler)

func NewRoleHandler(opts ...RoleHandlerOption) RoleHandler {
	h := &roleHandler{
		client: rpc.GetRoleServiceClient(),
		model:  NewRoleModel(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Create godoc
// @Summary create role
// @Schemes
// @Description create a new role
// @Tags role
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Create role request"
// @Success 200 {object} middleware.Response{data=CreateResponse}
// @Router /role [post]
func (h *roleHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	resp, err := h.client.Create(c.Request.Context(), &rolepb.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Level:       int32(req.Level),
		State:       uint32(req.State),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.Succ(c, "Role created successfully", h.model.ToCreateResponse(resp))
}

// Get godoc
// @Summary get role by id
// @Schemes
// @Description retrieve a single role by its id
// @Tags role
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} middleware.Response{data=GetResponse}
// @Router /role/{id} [get]
func (h *roleHandler) Get(c *gin.Context) {
	id, err := middleware.ParseUintParam(c, "id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	resp, err := h.client.Get(c.Request.Context(), &rolepb.GetRequest{Id: uint64(id)})
	if err != nil {
		middleware.Fail(c, "Role not found")
		return
	}
	middleware.Succ(c, "Role retrieved successfully", h.model.ToGetResponse(resp))
}

// List godoc
// @Summary list roles
// @Schemes
// @Description retrieve roles with pagination
// @Tags role
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Param keywords query string false "Role keyword"
// @Param level query int false "Role level"
// @Param state query int false "Role state"
// @Success 200 {object} middleware.Response{data=middleware.Pagination[GetResponse]}
// @Router /role [get]
func (h *roleHandler) List(c *gin.Context) {
	page, pageSize := middleware.ParsePageParams(c)
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	resp, err := h.client.List(c.Request.Context(), &rolepb.ListRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keywords: middleware.OptionalZero(req.KeyWords),
		Level:    middleware.OptionalDefault(req.Level, math.MaxInt32),
		State:    middleware.OptionalZero(req.State),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	r := h.model.ToListResponse(resp)
	r.Page = page
	r.PageSize = pageSize
	middleware.Succ(c, "Roles retrieved successfully", r)
}

// Update godoc
// @Summary update role
// @Schemes
// @Description update a role by id
// @Tags role
// @Produce json
// @Param id path int true "Role ID"
// @Param request body UpdateRequest true "Update role request"
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /role/{id} [put]
func (h *roleHandler) Update(c *gin.Context) {
	id, err := middleware.ParseUintParam(c, "id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	_, err = h.client.Update(c.Request.Context(), &rolepb.UpdateRequest{
		Id:          uint64(id),
		Name:        req.Name,
		Description: req.Description,
		Level:       int32(req.Level),
		State:       uint32(req.State),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "Role updated successfully")
}

// Delete godoc
// @Summary delete role
// @Schemes
// @Description delete a role by id
// @Tags role
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /role/{id} [delete]
func (h *roleHandler) Delete(c *gin.Context) {
	id, err := middleware.ParseUintParam(c, "id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	_, err = h.client.Delete(c.Request.Context(), &rolepb.DeleteRequest{Id: uint64(id)})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "Role deleted successfully")
}
