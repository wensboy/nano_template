package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/middleware"
	"example.com/nano_template/pkg/services/common"
	userSys "example.com/nano_template/pkg/services/user_sys"
	"example.com/nano_template/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server wraps the gin.Engine and http.Server for graceful shutdown and configuration.
type Server struct {
	engine *gin.Engine
	server *http.Server
}

// Option defines a function type for configuring the server.
type Option func(*Server)

// WithPort sets the server's listening port.
func WithPort(port string) Option {
	return func(s *Server) {
		s.server.Addr = ":" + port
	}
}

// WithReadTimeout sets the server's read timeout.
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets the server's write timeout.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// NewServer initializes a new Server with default settings and applies options.
func NewServer(cfg *config.ServerConfig, opts ...Option) *Server {
	engine := gin.Default()
	// 先运用 server config 中的配置, 后执行 option function 覆盖配置
	// 后续检查参数时遵循: cli flag > env > function call > config file
	addr := cfg.Host + ":" + cfg.Port
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	s := &Server{
		engine: engine,
		server: server,
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start starts the server and listens for incoming requests.
func (s *Server) Start(cfg *config.Config) {
	util.Info("Preparing server")
	s.mountConfigCheck(cfg)
	s.mountDatabase(&cfg.DatabaseConfig)
	s.mountHttpProxy(&cfg.HttpProxyConfig)
	s.mountGlobalMiddleware()
	s.mountStatic(&cfg.WebConfig)
	s.mountRouter(cfg)

	util.Info("Starting server")

	// Graceful shutdown handling
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Error("Server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	util.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		util.Error("Server forced to shutdown", zap.Error(err))
	}

	util.Info("Server exiting")
}

func (s *Server) mountGlobalMiddleware() {
	s.engine.Use(middleware.GlobalErrorHandler())
}

func (s *Server) mountStatic(cfg *config.WebConfig) {
	if !cfg.ServerStatic {
		util.Info("Skip mount static server...")
		return
	}
	distDir := filepath.Clean(cfg.Dist)
	indexFile := filepath.Join(distDir, cfg.Entry)

	if _, err := os.Stat(indexFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			util.Info("Frontend dist not found, skip static mount")
			return
		}
		util.Error("Failed to inspect frontend dist", zap.Error(err))
		return
	}

	s.engine.Static("/assets", filepath.Join(distDir, "assets"))
	s.engine.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))

	apiBasePath := strings.TrimRight(strings.TrimSpace(cfg.BaseUri), "/")
	s.engine.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if apiBasePath != "" && (requestPath == apiBasePath || strings.HasPrefix(requestPath, apiBasePath+"/")) {
			middleware.Erro(c, http.StatusNotFound, "API endpoint not found")
			return
		}

		staticFile := filepath.Join(distDir, filepath.FromSlash(strings.TrimPrefix(requestPath, "/")))
		if relPath, err := filepath.Rel(distDir, staticFile); err == nil &&
			relPath != ".." &&
			!strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
			if fileInfo, err := os.Stat(staticFile); err == nil && !fileInfo.IsDir() {
				c.File(staticFile)
				return
			}
		}

		if filepath.Ext(requestPath) != "" {
			middleware.Erro(c, http.StatusNotFound, "Static file not found")
			return
		}

		c.File(indexFile)
	})
	util.Info("mount static server...")
}

func (s *Server) mountRouter(cfg *config.Config) {
	v1 := s.engine.Group(cfg.WebConfig.BaseUri)
	// basic endpoints
	{
		common.MountCommonRouter(v1, cfg)
	}
	// bussiness endpoints
	{
		userSys.MountUserSysRouter(v1, cfg)
	}
}

func (s *Server) mountDatabase(dbConfig *config.DatabaseConfig) {
	if dbConfig.Enable {
		config.InitDB(dbConfig)
		if dbConfig.AutoMigrate {
			autoMigrate(config.GDB)
		}
	}
}

func (s *Server) mountHttpProxy(cfg *config.HttpProxyConfig) {
	middleware.InitHttpProxy(
		middleware.WithHttpProxyTimeout(cfg.Timeout),
	)
}

func (s *Server) mountConfigCheck(cfg *config.Config) {
	if cfg.LLMConfig.ActiveProvider >= len(cfg.LLMConfig.Providers) {
		cfg.LLMConfig.ActiveProvider = 0
	}
	addrParts := strings.Split(s.server.Addr, ":")
	if cfg.FlagConfig.Port > 1024 {
		addrParts[1] = strconv.Itoa(cfg.FlagConfig.Port)
	}
	if cfg.FlagConfig.Host != "" {
		addrParts[0] = cfg.FlagConfig.Host
	}
	s.server.Addr = strings.Join(addrParts, ":")
}

func autoMigrate(db *gorm.DB) {
	if db == nil {
		return
	}

	if err := db.AutoMigrate(
		&userSys.User{},
		&userSys.UserProfile{},
	); err != nil {
		util.Error("Auto migrate failed", zap.Error(err))
	}
}
