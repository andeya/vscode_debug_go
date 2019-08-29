package vscode_debug_go

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/henrylee2cn/goutil"
)

func init() {
	type Config struct {
		Version        string                   `json:"version"`
		Configurations []map[string]interface{} `json:"configurations"`
	}
	pid := os.Getpid()
	fmt.Println("Get processId:", pid)
	for _, src := range build.Default.SrcDirs() {
		p := filepath.Join(src, ".vscode", "launch.json")
		goutil.RewriteFile(p, func(cnt []byte) ([]byte, error) {
			var v Config
			if err := json.Unmarshal(cnt, &v); err != nil || len(v.Configurations) == 0 {
				if err != nil {
					fmt.Printf("Unmarshal error: %s, %v\n", p, err)
				}
				return cnt, nil
			}
			var has bool
			for _, m := range v.Configurations {
				if _, ok := m["processId"]; ok {
					m["processId"] = pid
					fmt.Println("Set processId to:", p)
					has = true
				}
			}
			if !has {
				return cnt, nil
			}
			return json.MarshalIndent(v, "", "    ")
		})
	}
}
