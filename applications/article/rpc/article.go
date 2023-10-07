package main

import (
	"flag"
	"fmt"

	"github.com/sybapp/infoflow/applications/article/rpc/internal/config"
	"github.com/sybapp/infoflow/applications/article/rpc/internal/server"
	"github.com/sybapp/infoflow/applications/article/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/article/rpc/pb"
	"github.com/sybapp/infoflow/pkg/interceptors"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterArticleServer(grpcServer, server.NewArticleServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
