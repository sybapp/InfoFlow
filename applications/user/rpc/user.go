package main

import (
	"flag"
	"fmt"

	"github.com/sybapp/infoflow/applications/user/rpc/internal/config"
	"github.com/sybapp/infoflow/applications/user/rpc/internal/server"
	"github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/user/rpc/pb"
	"github.com/sybapp/infoflow/pkg/interceptors"

	"github.com/zeromicro/go-zero/core/conf"
	cs "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == cs.DevMode || c.Mode == cs.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
