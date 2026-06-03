package config

import (
	"encoding/json"
	"os"

	"example.com/nano_template/cmd/docs"
	"example.com/nano_template/pkg/util"
	"go.uber.org/zap"
)

type (
	SystemRole struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Level       int    `json:"level"`
		State       int    `json:"state"`
	}
	SystemRoleConfig struct {
		Version     string       `json:"version"`
		Description string       `json:"description"`
		Roles       []SystemRole `json:"roles"`
	}
	SwaggerConfig struct {
		Title          string `json:"title"`
		Version        string `json:"v0.1"`
		Description    string `json:"description"`
		TermsOfService string `json:"termsOfService"`
		Host           string `json:"host"`
		BasePath       string `json:"basePath"`
	}
	SystemConfig struct {
		RoleConfigPath    string `yaml:"roleConfigPath"`
		RoleConfig        SystemRoleConfig
		SwaggerConfigPath string `yaml:"swaggerConfigPath"`
		SwaggerConfig     SwaggerConfig
	}
)

func DefaultSystemConfig() SystemConfig {
	return SystemConfig{
		RoleConfigPath: "",
		RoleConfig:     SystemRoleConfig{},
	}
}

func (c *SystemConfig) LoadRole() error {
	if c.RoleConfigPath == "" {
		util.Warn("[system]", zap.String("load-role", "empty role config path"))
		return nil
	}
	fd, err := os.OpenFile(c.RoleConfigPath, os.O_RDONLY, 0655)
	if err != nil {
		util.Warn("[system]", zap.String("load-role", "failed to open role config file"), zap.Error(err))
		return err
	}
	defer fd.Close()
	if err := json.NewDecoder(fd).Decode(&c.RoleConfig); err != nil {
		util.Warn("[system]", zap.String("load-role", "failed to decode role config"), zap.Error(err))
		return err
	}
	util.Info("[system]", zap.String("load-role", "succ"), zap.Int("list-role", len(c.RoleConfig.Roles)))
	return nil
}

func (c *SystemConfig) LoadSwagger() error {
	if c.SwaggerConfigPath == "" {
		util.Warn("[system]", zap.String("load-swagger", "empty swagger config path"))
		return nil
	}
	fd, err := os.OpenFile(c.SwaggerConfigPath, os.O_RDONLY, 0655)
	if err != nil {
		util.Warn("[system]", zap.String("load-swagger", "failed to open swagger config file"), zap.Error(err))
		return err
	}
	defer fd.Close()
	if err := json.NewDecoder(fd).Decode(&c.SwaggerConfig); err != nil {
		util.Warn("[system]", zap.String("load-swagger", "failed to decode swagger config"), zap.Error(err))
		return err
	}
	docs.SwaggerInfo.Title = c.SwaggerConfig.Title
	docs.SwaggerInfo.Version = c.SwaggerConfig.Version
	docs.SwaggerInfo.Description = c.SwaggerConfig.Description
	docs.SwaggerInfo.Host = c.SwaggerConfig.Host
	docs.SwaggerInfo.BasePath = c.SwaggerConfig.BasePath
	util.Info("[system]", zap.String("load-swagger", "succ"))
	return nil
}
