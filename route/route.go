package route

import (
	"flex/iface"
	"fmt"
	"strings"
)

type Server struct {
	Host   string
	Port   int
	Weight int
}

type Route struct {
	Name     string
	Desc     string
	Tags     []string
	Priority string

	Hosts    []string
	Paths    []string
	Rewrites [2]string

	Servers []Server

	Timeout int

	Plugins map[string]iface.Plugin
}

func (r *Route) Match(host, path string) bool {

	return false
}

func (r *Route) Route(host string) string {
	for _, v := range r.Servers {
		if v.Port == 0 || v.Port == 80 {
			return v.Host
		} else {
			return fmt.Sprintf("%s:%v", v.Host, v.Port)
		}
	}
	return ""
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
