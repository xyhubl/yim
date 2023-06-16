package http

import (
	"github.com/gin-gonic/gin"
	"github.com/xyhubl/yim/internal/logic/conf"
	"github.com/xyhubl/yim/internal/logic/http/middleware"
	v1 "github.com/xyhubl/yim/internal/logic/http/router/v1"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

var httpSrvHandler *http.Server

func Server(c *conf.Config) {
	gin.SetMode(c.Base.DebugModule)
	r := v1.InitRouter(middleware.Cors(), middleware.TranslationMiddleware(), gin.Recovery(), gin.Logger())

	httpSrvHandler = &http.Server{
		Addr:         c.HttpServer.Addr,
		Handler:      r,
		ReadTimeout:  time.Duration(c.HttpServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.HttpServer.WriteTimeout) * time.Second,
	}

	go func() {
		log.Println("[INFO] HTTP server start.", c.HttpServer.Addr)
		if err := httpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] HttpServerRun:%s err:%v\n", c.HttpServer.Addr, err)
		}
	}()
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrvHandler.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpServerStop err: %v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
