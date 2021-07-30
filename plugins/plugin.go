package plugins

import "sync"

var pool *PluginPool

var once sync.Once

func GetPool() *PluginPool {
	if pool == nil {
		once.Do(func() {
			pool = NewPluginPool()
		})
		return pool
	}
	return pool
}

type PluginPool struct {
	sync.RWMutex
	ps map[string]Plugin
}

func NewPluginPool() *PluginPool {
	return &PluginPool{
		ps: make(map[string]Plugin),
	}
}

func (pool *PluginPool) Set(name string, p Plugin) bool {
	pool.Lock()
	defer pool.Unlock()
	_, ok := pool.ps[name]
	pool.ps[name] = p
	return ok
}

func (pool *PluginPool) Get(name string) (Plugin, bool) {
	pool.RLock()
	defer pool.RUnlock()
	p, ok := pool.ps[name]
	return p, ok
}

func (pool *PluginPool) Del(name string) bool {
	pool.Lock()
	defer pool.Unlock()
	_, ok := pool.ps[name]
	delete(pool.ps, name)
	return ok
}
