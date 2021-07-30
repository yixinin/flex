package main

import (
	"context"
	"flag"
	"flex/http"
	"flex/plugins"
)

var (
	pluginsDir string

	HttpAddr      string
	AdminHttpAddr string
	TcpAddr       string
	UdpAddr       string
)

func main() {
	flag.StringVar(&HttpAddr, "http-addr", ":8080", "listen addr")
	flag.StringVar(&AdminHttpAddr, "admin-http-addr", ":8082", "admin listen addr")
	flag.StringVar(&pluginsDir, "plugin", "./plugins", "plugins dir")
	flag.Parse()
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	go plugins.LoadPlugins(ctx, pluginsDir)
	http.InitHttp(HttpAddr)
}
