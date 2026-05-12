package rpc

import (
	"net"
	"time"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc/services/common"
	"google.golang.org/grpc"
)

type rpcServer struct {
	quit       chan struct{}
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewRpcServer(cfg *config.RpcConfig) *rpcServer {
	lis, err := net.Listen("tcp", cfg.Host+":"+cfg.Port)
	if err != nil {
		panic(err)
	}
	return &rpcServer{
		quit:     make(chan struct{}),
		listener: lis,
		grpcServer: grpc.NewServer(
			grpc.ChainUnaryInterceptor(loggerInterceptor),
		),
	}
}

func (s *rpcServer) start() {
	s.mountServer()

	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()

	<-s.quit
	s.grpcServer.Stop()
}

func (s *rpcServer) StartBg() {
	go func() {
		s.start()
	}()
}

func (s *rpcServer) Start() {
	s.start()
}

func (s *rpcServer) Stop(delay time.Duration) {
	time.Sleep(delay)
	s.quit <- struct{}{}
}

func (s *rpcServer) mountServer() {
	common.MountCommonServer(s.grpcServer)
}
