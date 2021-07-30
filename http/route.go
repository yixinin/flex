package http

import (
	"flex/iface"
	"math/rand"
	"strings"
)

type Server struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}

type Route struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Tags     []string `json:"tags"`
	Priority int      `json:"priority"`
	Timeout  int      `json:"timeout"`

	Hosts    []string  `json:"hosts"`
	Paths    []string  `json:"paths"`
	Rewrites [2]string `json:"reweite"`

	Servers      []Server `json:"servers"`
	serverIndexs []int    `json:"-"`

	Configs map[string]map[string]interface{} `json:"configs"`
	Plugins map[string]iface.Plugin           `json:"-"`
}

func (r *Route) Match(host, path string) bool {
	var hostMatched = false
	var pathMatched = false
	for _, v := range r.Hosts {
		if host == v {
			hostMatched = true
			break
		}
	}
	for _, v := range r.Paths {
		if path == v {
			pathMatched = true
			break
		}
	}
	return hostMatched && pathMatched
}

func (r *Route) Route(host string) string {
	return r.RandRoute().Addr
}

func (r *Route) RandRoute() Server {
	r.updateBalancer()
	// 随机取一个server
	var i = rand.Intn(len(r.serverIndexs))
	idx := r.serverIndexs[i]
	if idx >= len(r.Servers) {
		if len(r.Servers) != 0 {
			return Server{}
		}
		return r.Servers[rand.Intn(len(r.Servers))]
	}
	return r.Servers[idx]
}

func (r *Route) updateBalancer() {
	if len(r.serverIndexs) == 0 && len(r.Servers) != 0 {
		r.serverIndexs = make([]int, 0, len(r.Servers)*2)
		for idx, v := range r.Servers {
			for i := 0; i < v.Weight; i++ {
				r.serverIndexs = append(r.serverIndexs, idx)
			}
		}
	}
}

func (r *Route) Rewrite(path string) string {
	if strings.HasPrefix(path, r.Rewrites[0]) {
		return strings.Replace(path, r.Rewrites[0], r.Rewrites[1], 1)
	}
	return path
}

func (r *Route) SetPlugin(p iface.Plugin) {
	r.Plugins[p.Name()] = p
}

type RouteSlice []Route

func (a RouteSlice) Len() int           { return len(a) }
func (a RouteSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RouteSlice) Less(i, j int) bool { return a[i].Priority < a[j].Priority }
