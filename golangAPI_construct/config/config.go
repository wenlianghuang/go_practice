package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	corsOnce    sync.Once
	corsOrigins []string
)

// LoadCORSOrigins loads white-list from (priority):
// 1. env CORS_ORIGINS (comma separated)
// 2. file config/cors_origins.txt (one origin per line, '#' comments)
func LoadCORSOrigins() []string {
	corsOnce.Do(func() {
		// 1. ENV
		if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
			for _, p := range strings.Split(raw, ",") {
				if v := strings.TrimSpace(p); v != "" {
					corsOrigins = append(corsOrigins, v)
				}
			}
			return
		}
		// 2. File
		path := filepath.Join("config", "cors_origins.txt")
		f, err := os.Open(path)
		if err != nil {
			return // silently fallback to empty (no cross-origin allowed)
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			corsOrigins = append(corsOrigins, line)
		}
	})
	return corsOrigins
}
