package plugins

import (
	"flex/iface"

	"github.com/valyala/fasthttp"
)

type Plugin struct {
	TimeStamp int64
	Version   string
	NewPlugin newPlugin
}

type newPlugin = func(string) iface.Plugin

type handleFunc = func(c fasthttp.RequestCtx) error
