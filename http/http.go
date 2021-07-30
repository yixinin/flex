package http

import (
	"encoding/json"
	"flex/iface"
	"flex/plugins"
	"log"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
)

var Routes = make(RouteSlice, 0, 10)

func InitHttp(addr string) {
	fasthttp.ListenAndServe(addr, handler)
}

func AddRoute(r Route) {
	for k := range r.Configs {
		pg, ok := plugins.GetPool().Get(k)
		if !ok {
			continue
		}
		p := pg.NewPlugin(r.Name)

		conf, err := json.Marshal(r.Configs[k])
		if err != nil {
			log.Println(err)
			continue
		}
		err = p.SetConfig(conf)
		if err != nil {
			log.Println(err)
			continue
		}
		r.Plugins[k] = p
	}
	Routes = append(Routes, r)
	sort.Sort(Routes)
}

func handler(c *fasthttp.RequestCtx) {
	var host = string(c.Host())
	var path = string(c.Path())

	addr := c.RemoteAddr()
	log.Println(addr.String(), host, path)

	var matched = false
	var r Route
	for _, v := range Routes {
		if v.Match(host, path) {
			r = v
			matched = true
			break
		}
	}
	if !matched {
		log.Println("not matched")
		c.SetStatusCode(500)
		c.WriteString("server not reachable")
		return
	}

	for _, p := range r.Plugins {
		if err := p.Handle(c); err != nil {
			switch err.(type) {
			case *iface.PluginError:
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

	newHost := r.Route(host)
	// 解析负载
	if newHost != "" {
		c.Request.SetHost(newHost)
	} else {
		c.SetStatusCode(500)
		c.Write([]byte("gateway error"))
		return
	}

	// 重写path
	newPath := r.Rewrite(path)
	c.URI().SetPath(newPath)

	err := fasthttp.DoTimeout(&c.Request, &c.Response, time.Second*time.Duration(r.Timeout))
	if err != nil {
		log.Println(err)
	}
}
