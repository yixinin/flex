package plugins

import (
	"context"
	"log"
	"os"
	"path"
	"plugin"
	"strings"
	"time"
)

var LoadConfigs func(ctx context.Context, dir string)

func LoadPlugins(ctx context.Context, dir string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Second)
		}
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

				name, pluginExsist := nameFuncSymbol.(*string)
				if !pluginExsist {
					log.Printf("%s not contains (func Name() string)\n", filename)
					continue
				}

				pg, pluginExsist := GetPool().Get(*name)

				if pluginExsist && pg.TimeStamp >= ts {
					continue
				}

				np, err := p.Lookup("NewPlugin")
				if err != nil {
					log.Println(err)
					continue
				}
				fn, ok := np.(newPlugin)
				if !ok {
					log.Println("plugin is not NewPlugin")
					continue
				}
				if !pluginExsist {
					log.Printf("load plugin %s \n", *name)
					GetPool().Set(*name, Plugin{
						TimeStamp: ts,
						NewPlugin: fn,
					})
					LoadConfigs(ctx, "./http")
				} else {
					GetPool().Del(*name)
					log.Printf("reload plugin %s \n", *name)
					GetPool().Set(*name, Plugin{
						TimeStamp: ts,
						NewPlugin: fn,
					})
				}
			}
		}
	}
}
