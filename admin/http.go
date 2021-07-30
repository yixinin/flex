package admin

import "github.com/valyala/fasthttp"

func InitHttp(addr string) {
	fasthttp.ListenAndServe(addr, handler)
}

func handler(c *fasthttp.RequestCtx) {

}
