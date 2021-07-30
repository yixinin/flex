package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	fasthttp.ListenAndServe(":8081", handler)
}

func handler(c *fasthttp.RequestCtx) {
	var s = fmt.Sprintf("host:%s,path:%s", c.Host(), c.Path())
	c.WriteString("hello sample server:" + s)
}
