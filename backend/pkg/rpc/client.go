package rpc

import (
	chatpb "example.com/nano_template/proto/chat"
	commonpb "example.com/nano_template/proto/common"
	osspb "example.com/nano_template/proto/oss"
	rolepb "example.com/nano_template/proto/role"
	userpb "example.com/nano_template/proto/user"
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

// GetCommonServiceClient returns a client for the CommonService.
func (c *RpcClient) GetCommonServiceClient() commonpb.CommonServiceClient {
	return commonpb.NewCommonServiceClient(c.grpcConn)
}

func GetCommonServiceClient() commonpb.CommonServiceClient {
	return _rpcClient.GetCommonServiceClient()
}

// GetRoleServiceClient returns a client for the RoleService.
func (c *RpcClient) GetRoleServiceClient() rolepb.RoleServiceClient {
	return rolepb.NewRoleServiceClient(c.grpcConn)
}

func GetRoleServiceClient() rolepb.RoleServiceClient {
	return _rpcClient.GetRoleServiceClient()
}

// GetUserServiceClient returns a client for the UserService.
func (c *RpcClient) GetUserServiceClient() userpb.UserServiceClient {
	return userpb.NewUserServiceClient(c.grpcConn)
}

func GetUserServiceClient() userpb.UserServiceClient {
	return _rpcClient.GetUserServiceClient()
}

// GetUserProfileServiceClient returns a client for the UserProfileService.
func (c *RpcClient) GetUserProfileServiceClient() userpb.UserProfileServiceClient {
	return userpb.NewUserProfileServiceClient(c.grpcConn)
}

func GetUserProfileServiceClient() userpb.UserProfileServiceClient {
	return _rpcClient.GetUserProfileServiceClient()
}

// GetAliyunOssServiceClient returns a client for the AliyunOssServiceClient.
func (c *RpcClient) GetAliyunOssServiceClient() osspb.AliyunOssServiceClient {
	return osspb.NewAliyunOssServiceClient(c.grpcConn)
}

func GetAliyunOssServiceClient() osspb.AliyunOssServiceClient {
	return _rpcClient.GetAliyunOssServiceClient()
}

// GetAiChatServiceClient returns a client for the AiChatService
func (c *RpcClient) GetAiChatServiceClient() chatpb.AiChatServiceClient {
	return chatpb.NewAiChatServiceClient(c.grpcConn)
}

func GetAiChatServiceClient() chatpb.AiChatServiceClient {
	return _rpcClient.GetAiChatServiceClient()
}
