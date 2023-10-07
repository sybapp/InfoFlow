package main

import (
	"flag"
	"fmt"

	"github.com/sybapp/infoflow/applications/article/api/internal/config"
	"github.com/sybapp/infoflow/applications/article/api/internal/handler"
	"github.com/sybapp/infoflow/applications/article/api/internal/svc"
	"github.com/sybapp/infoflow/pkg/xcode"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/article-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
