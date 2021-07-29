package main

import (
	"context"
	"flag"
)

var (
	HttpAddr string
	TcpAddr  string
	UdpAddr  string
)

func main() {
	flag.StringVar(&HttpAddr, "http-addr", ":8080", "listen addr")
	flag.StringVar(&pluginsDir, "plugin", "./plugins", "plugins dir")
	flag.Parse()
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	go loadPlugins(ctx)
	InitHttp()
	InitTcp(TcpAddr)
}
