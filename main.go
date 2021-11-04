package main

import (
	"backend-image-server/app"
	"backend-image-server/docs"
	"backend-image-server/pkg/config"
	"fmt"
	"net/http"

	_ "backend-image-server/docs"

	"github.com/sirupsen/logrus"
)

// @title Imloader Server API
// @version 0.1
// @description API for images
// @Schemes http https
// @in header
// @name Authorization
func main() {
	config.Init()
	r := app.Setup()
	cfg := config.Get()

	docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s/", cfg.GlobalPrefix)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	var listenErr error
	if cfg.UseSSL {
		logrus.Infof("Starting listening server with SSL on %s", addr)
		listenErr = http.ListenAndServeTLS(addr, "cert.pem", "privkey.pem", r)
	} else {
		logrus.Infof("Starting listening server without SSL on %s", addr)
		listenErr = http.ListenAndServe(addr, r)
	}
	logrus.Fatalf("Listen err :%s ", listenErr)
}
