// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/solikewind/happyeat/app/internal/config"
	"github.com/solikewind/happyeat/app/internal/handler"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/common/consts/code/errorx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/happyeatservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(
		func(h http.Header) {
			h.Set("Access-Control-Allow-Origin", "*")
			h.Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			h.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		},
		func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "CORS 拒绝"}`))
		}, "*"))
	defer server.Stop()

	// 服务从 app 目录启动时，Swagger JSON 位于项目根目录。
	openapiPath := filepath.Clean(filepath.Join("..", "happyeat.json"))

	ctx, err := svc.NewServiceContext(c) // 创建相关的服务上下文（db、client）
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.DB.Close()

	server.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/health",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok","service":"happyeat-api"}`))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/openapi/happyeat.json",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if _, err := os.Stat(openapiPath); err != nil {
					http.Error(w, "openapi file not found", http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				http.ServeFile(w, r, openapiPath)
			},
		},
	})

	httpx.SetErrorHandler(func(err error) (int, any) {
		switch e := err.(type) {
		case *errorx.CodeError:
			// 此时 e.Code 是数字，e.Msg 是通过 stringer 自动生成的中文提示
			return http.StatusOK, e
		default:
			return http.StatusInternalServerError, errorx.NewDefaultError(err.Error())
		}
	})
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
