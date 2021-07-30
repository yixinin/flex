package http

import (
	"encoding/json"
	"flex/iface"
	"flex/plugins"
	"log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		return &fasthttp.Client{}
	},
}

func init() {
	plugins.UpdatePlugins = UpdatePlugins
}

var routes *RouteTable
var once sync.Once

func GetRoutes() *RouteTable {
	if routes == nil {
		once.Do(func() {
			routes = NewRouteTable()
		})
		return routes
	}

	return routes
}

func UpdatePlugins(name string) {
	GetRoutes().Foreach(func(r Route) bool {
		_, ok := r.Plugins[name]
		if !ok {
			return false
		}
		delete(r.Plugins, name)
		pg, ok := plugins.GetPool().Get(name)
		if !ok {
			log.Println("plugin pool", name, "not exsist")
			return false
		}

		p := pg.NewPlugin(r.Name)
		c, ok := r.Configs[name]
		if !ok {
			log.Println("config", name, "not exsist")
			return false
		}
		buf, err := json.Marshal(c)
		if err != nil {
			log.Println(err)
			return false
		}
		err = p.SetConfig(buf)
		if err != nil {
			log.Println(err)
			return false
		}
		r.Plugins[name] = p
		return false
	})
}

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
	if !GetRoutes().Add(r) {
		log.Println("route add fail, duplicate route name")
	}
}

func handler(c *fasthttp.RequestCtx) {
	var host = string(c.Host())
	var path = string(c.Path())

	var matched = false
	var r Route

	GetRoutes().Foreach(func(route Route) bool {
		if route.Match(host, path) {
			r = route
			matched = true
			return true
		}
		return false
	})

	if !matched {
		log.Println(host, path, "not matched")
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

	client := clientPool.Get().(*fasthttp.Client)
	err := client.DoTimeout(&c.Request, &c.Response, time.Second*time.Duration(r.Timeout))
	if err != nil {
		log.Println(err)
	}
	clientPool.Put(client)
}
