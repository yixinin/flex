package main

import (
	"encoding/json"
	"flex/iface"
	"flex/memcache"
	"fmt"
	"log"

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

func (l *LimitCount) SetConfig(config []byte) error {
	return json.Unmarshal(config, &l.config)
}

func (l *LimitCount) Handle(c *fasthttp.RequestCtx) error {
	var key string
	switch l.config.Key {
	case ServerAddr:
		key = fmt.Sprintf("%s/%s/%s", l.key, Name, c.Host())
	case KeyRemoteAddr:
		key = fmt.Sprintf("%s/%s/%s", l.key, Name, c.RemoteAddr())
	default:
		var val = string(c.Request.Header.Peek(l.config.Key))
		key = fmt.Sprintf("%s/%s/%s", l.key, Name, val)
	}

	count := l.cache.Inc(key, 1, l.config.Duration)
	log.Println(key, count, l.config.Duration)
	if count > l.config.Limit {
		c.SetStatusCode(l.config.Code)
		return iface.Error(l.config.Code, "request limited")
	}
	return nil
}
