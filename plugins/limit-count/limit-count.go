package main

import (
	"encoding/json"
	"flex/iface"
	"flex/memcache"
	"fmt"

	"github.com/valyala/fasthttp"
)

const (
	KeyRemoteAddr = "remote_addr"
	ServerAddr    = "server_addr"
)

var Name string = "limit-count"

func NewPlugin(name string) iface.Plugin {
	return NewLimitCount(name)
}

type Config struct {
	Key      string `json:"key"`
	Limit    int    `json:"limit"`
	Duration int    `json:"duration"`
	Code     int    `json:"code"`
}
type LimitCount struct {
	key    string
	cache  *memcache.Cache
	config Config
}

func NewLimitCount(route string) iface.Plugin {
	return &LimitCount{
		key:   route,
		cache: memcache.NewCache(1024),
	}
}

func (l *LimitCount) Name() string {
	return Name
}

func (l *LimitCount) SetConfig(config string) error {
	return json.Unmarshal([]byte(config), &l.config)
}

func (l *LimitCount) Handle(c *fasthttp.RequestCtx) error {
	var val = string(c.Request.Header.Peek(l.config.Key))
	var key = fmt.Sprintf("%s/%s/%s", l.key, Name, val)
	count := l.cache.Inc(key, 1, l.config.Duration)
	if count > l.config.Limit {
		c.SetStatusCode(l.config.Code)
		return iface.ErrorAbort
	}
	return nil
}
