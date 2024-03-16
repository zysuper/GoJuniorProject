package grpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"time"
)

type GoZeroConsulTestSuite struct {
	suite.Suite
}

func (s *GoZeroTestSuite) TestGoZeroConsulClient() {
	// zrpc/registry/consul/resolver.go 的 `init()` 方法，回调用 resolver.Register(&builder{})
	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, "round_robin")
	zClient := zrpc.MustNewClient(zrpc.RpcClientConf{
		// 手动指定 Target 为 consul 服务器地址+调用的服务名.
		Target: "consul://127.0.0.1:8500/user",
	},
		// 指定负载均衡策略.
		zrpc.WithDialOption(grpc.WithDefaultServiceConfig(svcCfg)))
	client := NewUserServiceClient(zClient.Conn())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetByID(ctx, &GetByIDRequest{
		Id: 123,
	})
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
}

func (s *GoZeroTestSuite) TestGoZeroConsulServer() {
	rpcConf := zrpc.RpcServerConf{
		ListenOn: ":8090",
	}
	consulConf := consul.Conf{
		Host: "127.0.0.1:8500",
		Key:  "user",
	}
	server := zrpc.MustNewServer(rpcConf, func(grpcServer *grpc.Server) {
		RegisterUserServiceServer(grpcServer, &Server{})
	})
	_ = consul.RegisterService(rpcConf.ListenOn, consulConf)

	server.Start()
}
