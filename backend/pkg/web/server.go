package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/util"
	"example.com/nano_template/pkg/web/middleware"
	"example.com/nano_template/pkg/web/services/common"
	"example.com/nano_template/pkg/web/services/native"
	"example.com/nano_template/pkg/web/services/role"
	"example.com/nano_template/pkg/web/services/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

// WebServer wraps the gin.Engine and http.Server for graceful shutdown and configuration.
type WebServer struct {
	engine *gin.Engine
	server *http.Server
}

// Option defines a function type for configuring the server.
type WebServerOption func(*WebServer)

// NewServer initializes a new Server with default settings and applies options.
func NewServer(opts ...WebServerOption) *WebServer {
	engine := gin.Default()
	server := &http.Server{
		Handler: engine,
	}

	s := &WebServer{
		engine: engine,
		server: server,
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *WebServer) mountModules(cfg *config.Config) {
	util.Info("[moudles]", zap.String("web-mount", "mounting"))
	s.mountDatabase(&cfg.DatabaseConfig)
	s.mountOss(cfg)
	s.mountMemDatabase(&cfg.ValkeyConfig)
	s.mountHttpProxy(&cfg.HttpProxyConfig)
	s.mountGlobalMiddleware()
	s.mountStatic(cfg)
	s.mountRpcClient(&cfg.RpcConfig)
	s.mountRouter(cfg)
	s.mountConfig(cfg)
	util.Info("[moudles]", zap.String("web-mount", "mounted"))
}

// Start starts the server and listens for incoming requests.
func (s *WebServer) Start(cfg *config.Config) {
	s.mountModules(cfg)

	util.Info("[web-server]", zap.String("listen", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		util.Error("[web-server]", zap.Error(err))
	}
}

func (s *WebServer) StartBg(cfg *config.Config) {
	s.mountModules(cfg)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Error("[web-server]", zap.Error(err))
		}
	}()
	util.Info("[web-server]", zap.String("listen", s.server.Addr))
}

func (s *WebServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		util.Error("[web-server]", zap.Error(err))
	}
	util.Info("[web-server]", zap.String("stopped", s.server.Addr))
}

func (s *WebServer) mountGlobalMiddleware() {
	s.engine.Use(middleware.ErrorHandler())
	util.Info("[modules]", zap.String("mount-global-middleware", "mounted"))
}

func (s *WebServer) mountStatic(cfg *config.Config) {
	if !cfg.WebConfig.ServerStatic {
		util.Info("[modules]", zap.String("mount-static", "skipped"), zap.String("reason", "option is disabled"))
		return
	}
	distDir := filepath.Clean(cfg.WebConfig.Dist)
	indexFile := filepath.Join(distDir, cfg.WebConfig.Entry)

	if _, err := os.Stat(indexFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			util.Info("[modules]", zap.String("mount-static", "skipped"), zap.String("reason", "frontend dist not found"))
			return
		}
		util.Error("[modules]", zap.String("mount-static", "failed"), zap.Error(err))
		return
	}

	s.engine.Static("/assets", filepath.Join(distDir, "assets"))
	s.engine.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))

	apiBasePath := strings.TrimRight(strings.TrimSpace(cfg.ServerConfig.BaseUri), "/")
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
	util.Info("[modules]", zap.String("mount-static", "mounted"))
}

func (s *WebServer) mountRouter(cfg *config.Config) {
	v1 := s.engine.Group(cfg.ServerConfig.BaseUri)
	// basic endpoints
	{
		common.MountCommonRouter(v1, cfg)
		native.MountNativeRouter(v1, cfg)
	}
	// bussiness endpoints
	{
		role.MountRoleRouter(v1, cfg)
		user.MountUserRouter(v1, cfg)
	}
	util.Info("[modules]", zap.String("mount-router", "mounted"))
}

func (s *WebServer) mountDatabase(dbConfig *config.DatabaseConfig) {
	if dbConfig.Enable {
		config.InitDB(dbConfig)
		if dbConfig.AutoMigrate {
			s.autoMigrate(config.GetGDB())
		}
		util.Info("[modules]", zap.String("mount-database", "mounted"))
	} else {
		util.Info("[modules]", zap.String("mount-database", "skipped"))
	}
}

func (s *WebServer) mountMemDatabase(vConfig *config.ValkeyConfig) {
	if vConfig.Enable {
		config.InitValkey(vConfig)
		util.Info("[modules]", zap.String("mount-valkey-database", "mounted"))
	} else {
		util.Info("[modules]", zap.String("mount-valkey-database", "skipped"))
	}
}

func (s *WebServer) mountHttpProxy(cfg *config.HttpProxyConfig) {
	middleware.InitHttpProxy(
		middleware.WithHttpProxyTimeout(cfg.Timeout),
	)
	util.Info("[modules]", zap.String("mount-http-proxy", "mounted"))
}

func (s *WebServer) mountConfig(cfg *config.Config) {
	if cfg.LLMConfig.ActiveProvider >= len(cfg.LLMConfig.Providers) {
		cfg.LLMConfig.ActiveProvider = 0
	}
	addrParts := make([]string, 2)
	addrParts[0] = cfg.ServerConfig.Host
	addrParts[1] = cfg.ServerConfig.Port
	if cfg.FlagConfig.Port > 1024 {
		addrParts[1] = strconv.Itoa(cfg.FlagConfig.Port)
	}
	if cfg.FlagConfig.Host != "" {
		addrParts[0] = cfg.FlagConfig.Host
	}
	if s.server.Addr == "" {
		s.server.Addr = strings.Join(addrParts, ":")
	}
	s.server.ReadTimeout = time.Duration(cfg.ServerConfig.ReadTimeout) * time.Second
	s.server.WriteTimeout = time.Duration(cfg.ServerConfig.WriteTimeout) * time.Second
	s.server.IdleTimeout = time.Duration(cfg.ServerConfig.IdleTimeout) * time.Second
	s.server.MaxHeaderBytes = cfg.ServerConfig.MaxHeaderBytes
	util.Info("[modules]", zap.String("mount-config", "mounted"))
}

func (s *WebServer) mountRpcClient(cfg *config.RpcConfig) {
	if !cfg.Enable {
		util.Info("[modules]", zap.String("mount-rpc-client", "skipped"), zap.String("reason", "option is disabled"))
		return
	}
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Warn("grpc client connection init fail", zap.Error(err))
	}
	client := rpc.NewRpcClient(rpc.RpcClientOption{
		Grpc: conn,
	})
	rpc.SetRpcClient(client)
	util.Info("[modules]", zap.String("mount-rpc-client", "mounted"))
}

func (s *WebServer) mountOss(cfg *config.Config) {
	// aliyun oss
	if cfg.AliyunOssConfig.Enable {
		config.InitAliyunOss(&cfg.AliyunOssConfig)
		util.Info("[modules]", zap.String("mount-oss", "mounted"))
	} else {
		util.Info("[modules]", zap.String("mount-aliyun-oss", "skipped"), zap.String("reason", "option is disabled"))
	}
}

func (s *WebServer) autoMigrate(db *gorm.DB) {
	if db == nil {
		return
	}

	if err := db.AutoMigrate(); err != nil {
		util.Error("[auto-migrate]", zap.Error(err))
	}
}
