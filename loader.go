package main

import (
	"context"
	"log"
	"os"
	"path"
	"plugin"
	"strings"
	"time"
)

var pluginsDir string

func loadPlugins(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Second)
		}
		var dir = pluginsDir
		dirs, err := os.ReadDir(dir)
		if err != nil {
			log.Println("read dir error", err)
			continue
		}
		for _, file := range dirs {
			filename := path.Join(dir, file.Name())
			info, err := file.Info()
			if err != nil {
				log.Println("file info error", err)
				continue
			}
			if info.Size() == 0 {
				continue
			}
			ts := info.ModTime().UnixNano()

			if strings.HasSuffix(filename, ".so") {
				p, err := plugin.Open(filename)
				if err != nil {
					log.Println("plugin open error", err)
					continue
				}

				nameFuncSymbol, err := p.Lookup("Name")
				if err != nil {
					log.Println(err)
					continue
				}

				name, ok := nameFuncSymbol.(*string)
				if !ok {
					log.Printf("%s not contains (func Name() string)\n", filename)
					continue
				}

				if p, ok := ps[*name]; ok && p.TimeStamp >= ts {
					continue
				}

				sb, err := p.Lookup("NewPlugin")
				if err != nil {
					log.Println(err)
					continue
				}
				fn, ok := sb.(newPlugin)
				if !ok {
					log.Println("plugin is not NewPlugin")
					continue
				}
				if _, ok := ps[*name]; !ok {
					log.Printf("load plugin %s \n", *name)
					ps[*name] = Plugin{
						TimeStamp: ts,
						plugin:    fn,
					}
				} else {
					delete(ps, *name)
					log.Printf("reload plugin %s \n", *name)
					ps[*name] = Plugin{
						TimeStamp: ts,
						plugin:    fn,
					}
				}
			}
		}
	}
}
