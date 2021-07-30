package memcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type CahceValue struct {
	value interface{}
	start int64
	ttl   int
}

func NewValue(val interface{}, ttl int) *CahceValue {
	v := &CahceValue{
		value: val,
		start: time.Now().Unix(),
		ttl:   ttl,
	}
	if ttl <= 0 {
		v.ttl = -1
	}
	return v
}

func (c *Cache) doTTL(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(500 * time.Millisecond)
			func() {
				defer func() {
					recover()
				}()
				var now = time.Now().Unix()
				c.Lock()
				defer c.Unlock()
				for k, v := range c.m {
					if v.ttl <= 0 {
						continue
					}
					if v.start+int64(v.ttl) <= now {
						log.Println("delete", k, v.start, v.ttl, now)
						delete(c.m, k)
					}
				}
			}()
		}
	}
}

type Cache struct {
	sync.Mutex
	m map[string]*CahceValue
}

func NewCache(cap int) *Cache {
	c := &Cache{
		m: make(map[string]*CahceValue, cap),
	}
	var ctx = context.Background()
	go c.doTTL(ctx)
	return c
}

func (c *Cache) Ttl(key string) int {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	if ok {
		return v.ttl
	}
	return -2
}

func (c *Cache) ExpireIn(key string, ttl int) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	if ok {
		v.ttl = ttl
	}
}

func (c *Cache) Set(key string, val interface{}, ttl int) {
	c.Lock()
	c.m[key] = NewValue(val, ttl)
	c.Unlock()
}

func (c *Cache) Inc(key string, val int, ttl int) int {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	if !ok {
		c.m[key] = NewValue(val, ttl)
		return val
	}
	if i, ok := v.value.(int); ok {
		c.m[key].value = i + val
		return i + val
	}
	panic("value is not int")
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	return v, ok
}

func (c *Cache) GetString(key string) (string, bool) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	if !ok {
		return "", ok
	}
	if s, ok := v.value.(string); ok {
		return s, ok
	}
	return fmt.Sprint(v), ok
}

func (c *Cache) GetInt(key string) (int, bool) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	if !ok {
		return 0, ok
	}
	if i, ok := v.value.(int); ok {
		return i, ok
	}
	panic("value is not int")
}
