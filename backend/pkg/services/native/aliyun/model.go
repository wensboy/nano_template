package aliyun

import "time"

type PresignRequest struct {
	ObjectKey *string `json:"object_key" binding:"required"`
	Mime      string  `json:"mime"`
	Size      int     `json:"size"`
	Sender    string  `json:"sender"` // 发送方, 作为标识字段, 无特殊含义
}

type SignedHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PresignResponse struct {
	SignedUrl     string         `json:"signed_url"`
	Method        string         `json:"method"`
	Expiration    time.Time      `json:"expiration"`
	SignedHeaders []SignedHeader `json:"signed_headers,omitempty"`
}
