package ioc

import (
	grpc2 "GoProj/wedy/interactive/grpc"
	"GoProj/wedy/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Name     string `yaml:"name"`
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int32  `yaml:"port"`
	}
	s := grpc.NewServer()
	intrSvc.Register(s)
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	return &grpcx.Server{
		Server:   s,
		Addr:     cfg.EtcdAddr,
		Port:     cfg.Port,
		EtcdAddr: cfg.EtcdAddr,
	}
}
