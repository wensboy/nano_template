package rpc

import (
	"net"
	"os"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc/middleware"
	"example.com/nano_template/pkg/rpc/services/common"
	aliyunoss "example.com/nano_template/pkg/rpc/services/native/oss/aliyun"
	"example.com/nano_template/pkg/rpc/services/role"
	"example.com/nano_template/pkg/rpc/services/user"
	"example.com/nano_template/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type RpcServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

type RpcServerOption func(*RpcServer)

func NewRpcServer(opts ...RpcServerOption) *RpcServer {
	s := &RpcServer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *RpcServer) mountModules(cfg *config.Config) {
	util.Info("[moudles]", zap.String("rpc-mount", "mounting"))
	s.mountDatabase(&cfg.DatabaseConfig)
	s.mountOss(cfg)
	s.mountMemDatabase(&cfg.ValkeyConfig)
	s.mountMiddleware(cfg)
	s.mountServer(cfg)
	s.mountConfig(cfg)
	util.Info("[moudles]", zap.String("rpc-mount", "mounted"))
}

func (s *RpcServer) Start(cfg *config.Config) {
	s.mountModules(cfg)

	if err := s.grpcServer.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
		util.Error(err.Error())
		os.Exit(1)
	}

}

func (s *RpcServer) StartBg(cfg *config.Config) {
	s.mountModules(cfg)

	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
			util.Error(err.Error())
			os.Exit(1)
		}
	}()
	util.Info("[server]", zap.String("start", "rpc"), zap.String("listen", cfg.ServerConfig.RpcServerConfig.Host+":"+cfg.ServerConfig.RpcServerConfig.Port))
}

func (s *RpcServer) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *RpcServer) mountDatabase(dbConfig *config.DatabaseConfig) {
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

func (s *RpcServer) mountOss(cfg *config.Config) {
	// aliyun oss
	if cfg.AliyunOssConfig.Enable {
		config.InitAliyunOss(&cfg.AliyunOssConfig)
		util.Info("[modules]", zap.String("mount-oss", "mounted"))
	} else {
		util.Info("[modules]", zap.String("mount-aliyun-oss", "skipped"), zap.String("reason", "option is disabled"))
	}
}

func (s *RpcServer) mountMemDatabase(vConfig *config.ValkeyConfig) {
	if vConfig.Enable {
		config.InitValkey(vConfig)
		util.Info("[modules]", zap.String("mount-valkey-database", "mounted"))
	} else {
		util.Info("[modules]", zap.String("mount-valkey-database", "skipped"))
	}
}

func (s *RpcServer) mountServer(cfg *config.Config) {
	common.MountCommonServer(s.grpcServer, cfg)
	role.MountRoleServer(s.grpcServer, cfg)
	user.MountUserServer(s.grpcServer, cfg)
	// native services
	aliyunoss.MountAliyunOssServer(s.grpcServer, cfg)
}

func (s *RpcServer) mountMiddleware(cfg *config.Config) {
	rpcServerConfig := cfg.ServerConfig.RpcServerConfig
	unaryInterceptors := []grpc.UnaryServerInterceptor{}
	var ui middleware.UnaryInterceptor
	streamInterceptors := []grpc.StreamServerInterceptor{}
	if rpcServerConfig.MiddlewareConfig.RecoveryConfig.Enable {
		unaryInterceptors = append(unaryInterceptors, ui.RecoveryInterceptor(&rpcServerConfig))
	}
	if rpcServerConfig.MiddlewareConfig.LoggerConfig.Enable {
		unaryInterceptors = append(unaryInterceptors, ui.LoggerInterceptor(&rpcServerConfig))
	}
	if rpcServerConfig.MiddlewareConfig.AuthConfig.Enable {
		unaryInterceptors = append(unaryInterceptors, ui.AuthInterceptor(&rpcServerConfig))
	}
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
}

func (s *RpcServer) mountConfig(cfg *config.Config) {
	address := cfg.ServerConfig.RpcServerConfig.Host + ":" + cfg.ServerConfig.RpcServerConfig.Port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		util.Error("[modules]", zap.String("mount-config", "failed"), zap.Error(err))
		os.Exit(1)
	}
	s.listener = lis
}

func (s *RpcServer) autoMigrate(db *gorm.DB) {
	if db == nil {
		return
	}

	if err := db.AutoMigrate(
		&role.Role{},
		&user.User{},
		&user.UserProfile{},
		&user.UserRole{},
	); err != nil {
		util.Error("[auto-migrate]", zap.Error(err))
	}
}
