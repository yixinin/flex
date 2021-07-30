package http

import (
	"context"
	"encoding/json"
	"flex/iface"
	"flex/plugins"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func init() {
	plugins.LoadConfigs = LoadConfigs
}

func LoadConfigs(ctx context.Context, dir string) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		log.Println("read dir error", err)
		return
	}
	for _, file := range dirs {
		var info, err = file.Info()
		if err != nil {
			log.Println(err)
			return
		}
		var filename = path.Join(dir, info.Name())
		if strings.HasSuffix(filename, ".json") {
			f, err := os.Open(filename)
			if err != nil {
				log.Println(err)
				continue
			}
			buf, err := io.ReadAll(f)
			if err != nil {
				log.Println(err)
				continue
			}
			var route Route
			err = json.Unmarshal(buf, &route)
			if err != nil {
				log.Println(err)
				continue
			}
			route.Plugins = make(map[string]iface.Plugin, len(route.Configs))
			AddRoute(route)
		}
	}
}
