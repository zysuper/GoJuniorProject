package ioc

import (
	grpc2 "gitee.com/geekbang/basic-go/webook/follow/grpc"
	"gitee.com/geekbang/basic-go/webook/pkg/grpcx"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(
	ecli *clientv3.Client,
	followRelation *grpc2.FollowServiceServer,
	l logger.LoggerV1,
) *grpcx.Server {
	type Config struct {
		Port    int   `yaml:"port"`
		EtcdTTL int64 `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	followRelation.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		EtcdClient: ecli,
		Name:       "follow",
		EtcdTTL:    cfg.EtcdTTL,
		L:          l,
	}
}
