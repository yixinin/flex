package main

import (
	"flex/iface"

	"github.com/valyala/fasthttp"
)

type Plugin struct {
	TimeStamp int64
	Version   string
	plugin    newPlugin
}

var ps = make(map[string]Plugin)

type newPlugin = func(string) iface.Plugin

type fn = func(c fasthttp.RequestCtx) error
