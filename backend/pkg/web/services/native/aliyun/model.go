package aliyun

import osspb "example.com/nano_template/proto/oss"

type PresignUploadRequest struct {
	ObjectKey *string `json:"object_key" binding:"required"`
	Mime      string  `json:"mime"`
	Size      int     `json:"size"`
	Sender    string  `json:"sender"` // 发送方, 作为标识字段, 无特殊含义
}

type PresignDownloadRequest struct {
	BucketPrefix string `json:"bucket_prefix"`
	ObjectKey    string `json:"object_key"`
	Getter       string `json:"getter"` // 获取方, 作为标识字段, 无特殊含义
}

type PresignListRequest struct {
	Getter string `json:"getter"` // 获取方, 作为标识字段, 无特殊含义
}

type ObjectOwner struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type ObjectProperties struct {
	Key            string       `json:"key"`
	Type           string       `json:"type"`
	Size           int64        `json:"size"`
	ETag           string       `json:"etag"`
	LastModified   string       `json:"last_modified"`
	StorageClass   string       `json:"storage_class"`
	Owner          *ObjectOwner `json:"owner,omitempty"`
	RestoreInfo    string       `json:"restore_info,omitempty"`
	TransitionTime string       `json:"transition_time,omitempty"`
}

type CommonPrefix struct {
	Prefix string `json:"prefix"`
}

type ListObjectsResponse struct {
	Name           string             `json:"name"`
	Prefix         string             `json:"prefix"`
	Marker         string             `json:"marker"`
	MaxKeys        int32              `json:"max_keys"`
	Delimiter      string             `json:"delimiter"`
	IsTruncated    bool               `json:"is_truncated"`
	NextMarker     string             `json:"next_marker"`
	EncodingType   string             `json:"encoding_type"`
	Contents       []ObjectProperties `json:"contents"`
	CommonPrefixes []CommonPrefix     `json:"common_prefixes"`
}

type SignedHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PresignResponse struct {
	SignedUrl     string         `json:"signed_url"`
	Method        string         `json:"method"`
	Expiration    string         `json:"expiration"`
	SignedHeaders []SignedHeader `json:"signed_headers,omitempty"`
}

type AliyunModel interface {
	ToPresignResponse(*osspb.PresignGeneralResponse) PresignResponse
	ToListObjectsResponse(*osspb.ListObjectsResponse) ListObjectsResponse
}

type aliyunModel struct{}

func NewAliyunModel() AliyunModel {
	return &aliyunModel{}
}

func (*aliyunModel) ToPresignResponse(resp *osspb.PresignGeneralResponse) PresignResponse {
	signedHeaders := make([]SignedHeader, 0)
	if resp != nil {
		signedHeaders = make([]SignedHeader, 0, len(resp.SignedHeaders))
		for _, header := range resp.SignedHeaders {
			signedHeaders = append(signedHeaders, SignedHeader{Key: header.Key, Value: header.Value})
		}
	}
	if resp == nil {
		return PresignResponse{SignedHeaders: signedHeaders}
	}
	return PresignResponse{
		SignedUrl:     resp.SignedUrl,
		Method:        resp.SignedMethod,
		Expiration:    resp.Expiration,
		SignedHeaders: signedHeaders,
	}
}

func (*aliyunModel) ToListObjectsResponse(resp *osspb.ListObjectsResponse) ListObjectsResponse {
	if resp == nil {
		return ListObjectsResponse{}
	}

	contents := make([]ObjectProperties, 0, len(resp.Contents))
	for _, object := range resp.Contents {
		contents = append(contents, toObjectProperties(object))
	}
	commonPrefixes := make([]CommonPrefix, 0, len(resp.CommonPrefixes))
	for _, commonPrefix := range resp.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: commonPrefix.Prefix})
	}

	return ListObjectsResponse{
		Name:           resp.Name,
		Prefix:         resp.Prefix,
		Marker:         resp.Marker,
		MaxKeys:        resp.MaxKeys,
		Delimiter:      resp.Delimiter,
		IsTruncated:    resp.IsTruncated,
		NextMarker:     resp.NextMarker,
		EncodingType:   resp.EncodingType,
		Contents:       contents,
		CommonPrefixes: commonPrefixes,
	}
}

func toObjectProperties(object *osspb.ObjectProperties) ObjectProperties {
	return ObjectProperties{
		Key:            object.Key,
		Type:           object.Type,
		Size:           object.Size,
		ETag:           object.Etag,
		LastModified:   object.LastModified,
		StorageClass:   object.StorageClass,
		Owner:          toObjectOwner(object.Owner),
		RestoreInfo:    object.RestoreInfo,
		TransitionTime: object.TransitionTime,
	}
}

func toObjectOwner(owner *osspb.ObjectOwner) *ObjectOwner {
	if owner == nil {
		return nil
	}
	return &ObjectOwner{
		ID:          owner.Id,
		DisplayName: owner.DisplayName,
	}
}
