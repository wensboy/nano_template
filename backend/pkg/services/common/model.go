package common

import "example.com/nano_template/pkg/config"

type (
	InspectResponse struct {
		Version     string `json:"version"`
		Author      string `json:"author"`
		Description string `json:"description"`
	}
	TemplateResponse struct {
		Id          string                     `json:"id"`
		FrontMatter config.TemplateFrontMatter `json:"frontmatter"`
		Content     string                     `json:"content"`
	}
)
