package aliyun

import (
	"time"

	osspb "example.com/nano_template/proto/oss"
	aliyunoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type AliyunOssModel interface {
	ToPresignResponse(*aliyunoss.PresignResult) *osspb.PresignGeneralResponse
	ToListObjectsResponse(*aliyunoss.ListObjectsResult) *osspb.ListObjectsResponse
}

type aliyunOssModel struct{}

func NewAliyunOssModel() AliyunOssModel {
	return &aliyunOssModel{}
}

func (*aliyunOssModel) ToPresignResponse(result *aliyunoss.PresignResult) *osspb.PresignGeneralResponse {
	if result == nil {
		return &osspb.PresignGeneralResponse{}
	}

	signedHeaders := make([]*osspb.SignedHeader, 0, len(result.SignedHeaders))
	for key, value := range result.SignedHeaders {
		signedHeaders = append(signedHeaders, &osspb.SignedHeader{Key: key, Value: value})
	}
	return &osspb.PresignGeneralResponse{
		SignedUrl:     result.URL,
		SignedMethod:  result.Method,
		Expiration:    result.Expiration.String(),
		SignedHeaders: signedHeaders,
	}
}

func (*aliyunOssModel) ToListObjectsResponse(result *aliyunoss.ListObjectsResult) *osspb.ListObjectsResponse {
	if result == nil {
		return &osspb.ListObjectsResponse{}
	}

	contents := make([]*osspb.ObjectProperties, 0, len(result.Contents))
	for _, object := range result.Contents {
		contents = append(contents, toObjectProperties(object))
	}
	commonPrefixes := make([]*osspb.CommonPrefix, 0, len(result.CommonPrefixes))
	for _, commonPrefix := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, &osspb.CommonPrefix{
			Prefix: aliyunoss.ToString(commonPrefix.Prefix),
		})
	}

	return &osspb.ListObjectsResponse{
		Name:           aliyunoss.ToString(result.Name),
		Prefix:         aliyunoss.ToString(result.Prefix),
		Marker:         aliyunoss.ToString(result.Marker),
		MaxKeys:        result.MaxKeys,
		Delimiter:      aliyunoss.ToString(result.Delimiter),
		IsTruncated:    result.IsTruncated,
		NextMarker:     aliyunoss.ToString(result.NextMarker),
		EncodingType:   aliyunoss.ToString(result.EncodingType),
		Contents:       contents,
		CommonPrefixes: commonPrefixes,
	}
}

func toObjectProperties(object aliyunoss.ObjectProperties) *osspb.ObjectProperties {
	return &osspb.ObjectProperties{
		Key:            aliyunoss.ToString(object.Key),
		Type:           aliyunoss.ToString(object.Type),
		Size:           object.Size,
		Etag:           aliyunoss.ToString(object.ETag),
		LastModified:   optionalTimeString(object.LastModified),
		StorageClass:   aliyunoss.ToString(object.StorageClass),
		Owner:          toObjectOwner(object.Owner),
		RestoreInfo:    aliyunoss.ToString(object.RestoreInfo),
		TransitionTime: optionalTimeString(object.TransitionTime),
	}
}

func toObjectOwner(owner *aliyunoss.Owner) *osspb.ObjectOwner {
	if owner == nil {
		return nil
	}
	return &osspb.ObjectOwner{
		Id:          aliyunoss.ToString(owner.ID),
		DisplayName: aliyunoss.ToString(owner.DisplayName),
	}
}

func optionalTimeString(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.String()
}
