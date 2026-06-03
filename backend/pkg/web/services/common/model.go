package common

import (
	"example.com/nano_template/pkg/config"
	commonpb "example.com/nano_template/proto/common"
)

type (
	InspectResponse struct {
		Version     string `json:"version"`
		Author      string `json:"author"`
		Description string `json:"description"`
	}
	GetTemplateResponse struct {
		Id          string                     `json:"id"`
		FrontMatter config.TemplateFrontMatter `json:"frontmatter"`
		Content     string                     `json:"content"`
	}
)

type CommonModel interface {
	ToInspectResponse(*commonpb.InspectResponse) InspectResponse
	ToGetTemplateResponse(string, *commonpb.GetTemplateResponse) GetTemplateResponse
}

type commonModel struct{}

func NewCommonModel() CommonModel {
	return &commonModel{}
}

func (*commonModel) ToInspectResponse(resp *commonpb.InspectResponse) InspectResponse {
	return InspectResponse{
		Version:     resp.Version,
		Author:      resp.Author,
		Description: resp.Description,
	}
}

func (*commonModel) ToGetTemplateResponse(template_id string, resp *commonpb.GetTemplateResponse) GetTemplateResponse {
	fm := config.TemplateFrontMatter{
		Role: resp.FrontMatter.Role,
	}
	return GetTemplateResponse{
		Id:          template_id,
		FrontMatter: fm,
		Content:     resp.Content,
	}
}
