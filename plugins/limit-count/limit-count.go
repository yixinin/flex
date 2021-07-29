package main

import (
	"flex/iface"
	"strconv"

	"github.com/valyala/fasthttp"
)

var Name string = "limit-count"

func NewPlugin(name string) iface.Plugin {
	return NewLimitCount(name)
}

type LimitCount struct {
	name string
	m    map[string]int
}

func NewLimitCount(name string) iface.Plugin {
	return &LimitCount{
		name: name,
		m:    make(map[string]int, 1024),
	}
}

func (l *LimitCount) Name() string {
	return l.name
}
func (l *LimitCount) Handle(c *fasthttp.RequestCtx) error {
	var host = string(c.Host())
	l.m[host]++
	c.Request.Header.Set("limit-count", strconv.Itoa(l.m[host]))
	if l.m[host] > 5 {
		return iface.Error(305, "limit gt 5")
	}
	return nil
}
