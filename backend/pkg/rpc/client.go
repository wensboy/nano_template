package rpc

import (
	common "example.com/nano_template/proto"
	"google.golang.org/grpc"
)

type (
	RpcClient struct {
		grpcConn *grpc.ClientConn
	}
	RpcClientOption struct {
		Grpc *grpc.ClientConn
	}
)

var (
	_rpcClient *RpcClient
)

func SetRpcClient(client *RpcClient) {
	_ = client.Alive()
	_rpcClient = client
}

func GetRpcClient() *RpcClient {
	return _rpcClient
}

func NewRpcClient(opt RpcClientOption) *RpcClient {
	return &RpcClient{grpcConn: opt.Grpc}
}

func (c *RpcClient) GetGrpcConn() *grpc.ClientConn {
	return c.grpcConn
}

func (c *RpcClient) Alive() bool {
	c.grpcConn.Connect()
	return true
}

func (c *RpcClient) Close() {
	if c.grpcConn != nil {
		c.grpcConn.Close()
	}
}

func (c *RpcClient) GetCommonServiceClient() common.CommonServiceClient {
	return common.NewCommonServiceClient(c.grpcConn)
}
