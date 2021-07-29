package main

import (
	"encoding/json"
	"errors"
	"flex/iface"
	"log"

	"github.com/valyala/fasthttp"
)

func InitHttp() {
	fasthttp.ListenAndServe(HttpAddr, handler)
}

var handles = make(map[string]iface.Plugin)

func handler(c *fasthttp.RequestCtx) {
	for _, p := range ps {
		host := string(c.Host())
		pp, ok := handles[host]
		if !ok {
			pp = p.plugin(host)
			handles[host] = pp
		}
		if err := pp.Handle(c); err != nil {
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

	err := fasthttp.Do(&c.Request, &c.Response)
	if err != nil {
		log.Println(err)
	}
}
