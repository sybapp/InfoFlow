// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/sybapp/infoflow/applications/applet/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: LoginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/info",
				Handler: UserInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithSignature(serverCtx.Config.Signature),
		rest.WithPrefix("/v1/user"),
	)
}
