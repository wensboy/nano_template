package ai

import (
	"errors"
	"io"
	"net/http"

	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/web/middleware"
	chatpb "example.com/nano_template/proto/chat"
	"github.com/gin-gonic/gin"
)

type AiChatHandler interface {
	CreateSession(c *gin.Context)
	UpdateSession(c *gin.Context)
	GetSession(c *gin.Context)
	ListSessions(c *gin.Context)
	DeleteSession(c *gin.Context)
	Generate(c *gin.Context)
	GenerateModal(c *gin.Context)
	Record(c *gin.Context)
	UpdateRecord(c *gin.Context)
	ListRecords(c *gin.Context)
	DeleteRecord(c *gin.Context)
	UploadAppendixs(c *gin.Context)
	UpdateAppendix(c *gin.Context)
	ListAppendixs(c *gin.Context)
	DeleteAppendix(c *gin.Context)
}

type aiChatHandler struct {
	client chatpb.AiChatServiceClient
}

type AiChatHandlerOption func(*aiChatHandler)

func NewAiChatHandler(opts ...AiChatHandlerOption) AiChatHandler {
	handler := &aiChatHandler{
		client: rpc.GetAiChatServiceClient(),
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}

func (h *aiChatHandler) CreateSession(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req chatpb.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}
	req.UserId = int32(userID)

	resp, err := h.client.CreateSession(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat session created successfully", resp)
}

func (h *aiChatHandler) UpdateSession(c *gin.Context) {
	var req chatpb.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}
	if sessionID := c.Param("session_id"); sessionID != "" {
		req.SessionId = sessionID
	}

	_, err := h.client.UpdateSession(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.SuccNoMore(c, "AI chat session updated successfully")
}

func (h *aiChatHandler) GetSession(c *gin.Context) {
	resp, err := h.client.GetSession(c.Request.Context(), &chatpb.GetSessionRequest{
		SessionId: c.Param("session_id"),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat session retrieved successfully", resp)
}

func (h *aiChatHandler) ListSessions(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, pageSize := middleware.ParsePageParams(c)

	resp, err := h.client.ListSessions(c.Request.Context(), &chatpb.ListSessionsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		UserId:   int32(userID),
		KeyWords: c.Query("key_words"),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat sessions retrieved successfully", resp)
}

func (h *aiChatHandler) DeleteSession(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	_, err := h.client.DeleteSession(c.Request.Context(), &chatpb.DeleteSessionRequest{
		UserId:    int32(userID),
		SessionId: c.Param("session_id"),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.SuccNoMore(c, "AI chat session deleted successfully")
}

func (h *aiChatHandler) Generate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req chatpb.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}
	req.UserId = int32(userID)

	stream, err := h.client.Generate(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	streamResponses(c, stream.Recv)
}

func (h *aiChatHandler) GenerateModal(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req chatpb.GenerateModalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}
	req.UserId = int32(userID)

	stream, err := h.client.GenerateModal(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	streamResponses(c, stream.Recv)
}

func (h *aiChatHandler) Record(c *gin.Context) {
	var req chatpb.RecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}

	resp, err := h.client.Record(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat record created successfully", resp)
}

func (h *aiChatHandler) UpdateRecord(c *gin.Context) {
	recordID, err := middleware.ParseUintParam(c, "record_id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid record ID")
		return
	}

	var req chatpb.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}
	req.RecordId = int32(recordID)

	_, err = h.client.UpdateRecord(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.SuccNoMore(c, "AI chat record updated successfully")
}

func (h *aiChatHandler) ListRecords(c *gin.Context) {
	page, pageSize := middleware.ParsePageParams(c)
	resp, err := h.client.ListRecords(c.Request.Context(), &chatpb.ListRecordsRequest{
		SessionId: c.Param("session_id"),
		Page:      int32(page),
		PageSize:  int32(pageSize),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat records retrieved successfully", resp)
}

func (h *aiChatHandler) DeleteRecord(c *gin.Context) {
	recordID, err := middleware.ParseUintParam(c, "record_id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid record ID")
		return
	}

	_, err = h.client.DeleteRecord(c.Request.Context(), &chatpb.DeleteRecordRequest{RecordId: int32(recordID)})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.SuccNoMore(c, "AI chat record deleted successfully")
}

func (h *aiChatHandler) UploadAppendixs(c *gin.Context) {
	var req chatpb.UploadAppendixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}

	resp, err := h.client.UploadAppendixs(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat appendix uploaded successfully", resp)
}

func (h *aiChatHandler) UpdateAppendix(c *gin.Context) {
	appendixID, err := middleware.ParseUintParam(c, "appendix_id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid appendix ID")
		return
	}

	var req chatpb.UpdateAppendixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}
	req.AppendixId = int32(appendixID)

	_, err = h.client.UpdateAppendix(c.Request.Context(), &req)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.SuccNoMore(c, "AI chat appendix updated successfully")
}

func (h *aiChatHandler) ListAppendixs(c *gin.Context) {
	messageID, err := middleware.ParseUintParam(c, "message_id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid message ID")
		return
	}

	resp, err := h.client.ListAppendixs(c.Request.Context(), &chatpb.ListAppendixsRequest{
		MessageId: int32(messageID),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "AI chat appendixs retrieved successfully", resp)
}

func (h *aiChatHandler) DeleteAppendix(c *gin.Context) {
	appendixID, err := middleware.ParseUintParam(c, "appendix_id")
	if err != nil {
		middleware.Erro(c, http.StatusBadRequest, "Invalid appendix ID")
		return
	}

	_, err = h.client.DeleteAppendix(c.Request.Context(), &chatpb.DeleteAppendixRequest{AppendixId: int32(appendixID)})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.SuccNoMore(c, "AI chat appendix deleted successfully")
}

func currentUserID(c *gin.Context) (uint, bool) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		middleware.Erro(c, http.StatusUnauthorized, "Invalid user ID")
		return 0, false
	}
	return userID, true
}

func streamResponses[T any](c *gin.Context, recv func() (T, error)) {
	c.Stream(func(w io.Writer) bool {
		resp, err := recv()
		if errors.Is(err, io.EOF) {
			return false
		}
		if err != nil {
			c.SSEvent("error", err.Error())
			return false
		}

		c.SSEvent("message", resp)
		return true
	})
}
