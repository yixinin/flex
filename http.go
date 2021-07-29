package main

import (
	"encoding/json"
	"errors"
	"flex/iface"
	"flex/plugins"
	"flex/route"
	"log"
	"sort"

	"github.com/valyala/fasthttp"
)

var Routes = make(route.RouteSlice, 0, 10)

func InitHttp() {
	fasthttp.ListenAndServe(HttpAddr, handler)
}

func AddRoute(r route.Route) {
	for k := range r.Plugins {
		r.Plugins[k] = plugins.Pool[k]
	}
	Routes = append(Routes, r)
	sort.Sort(Routes)
}

func handler(c *fasthttp.RequestCtx) {
	var host = string(c.Host())
	var path = string(c.Path())
	var matched = false
	var r route.Route
	for _, v := range Routes {
		if v.Match(host, path) {
			r = v
			matched = true
			break
		}
	}
	if !matched {
		return
	}

	for _, p := range r.Plugins {
		if err := p.Handle(c); err != nil {
			if errors.Is(err, iface.PluginError{}) {
				b, err := json.Marshal(err)
				if err != nil {
					c.Write([]byte(err.Error()))
					return
				}
				c.Write(b)
				return
			}
			c.Write([]byte(err.Error()))
			return
		}
	}

	// 解析负载
	if host := r.Route(host); host != "" {
		c.Request.SetHost(host)
	} else {
		c.SetStatusCode(500)
		c.Write([]byte("gateway error"))
		return
	}

	// 重写path
	c.URI().SetPath(r.Rewrite(path))

	err := fasthttp.Do(&c.Request, &c.Response)
	if err != nil {
		log.Println(err)
	}
}
