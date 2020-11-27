package env

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var subRe = regexp.MustCompile(`([^$]*)\$([a-zA-Z0-9_]+)(.*)`)

func ReadFrom(paths ...string) (result map[string]string) {
	result = make(map[string]string)

	for _, path := range paths {
		entries, err := ioutil.ReadDir(path)
		if err == nil {
			for _, entry := range entries {
				if entry.Mode().IsRegular() && strings.HasSuffix(entry.Name(), ".env") {
					if file, err := os.Open(filepath.Join(path, entry.Name())); err == nil {
						scanner := bufio.NewScanner(file)
						for scanner.Scan() {
							if line := scanner.Text(); !strings.HasPrefix(line, "#") {
								if kv := strings.SplitN(line, "=", 2); len(kv) == 2 {
									result[kv[0]] = kv[1]
								}
							}
						}
						file.Close()
					}
				}
			}
		}
	}

	// perform recursive substition of all environment vaalues
	substituted := true
	for substituted {
		substituted = false
		for k, v := range result {
			if matches := subRe.FindStringSubmatch(v); matches != nil {
				subKey := matches[2]
				if subVal, ok := result[subKey]; ok {
					result[k] = matches[1] + subVal + matches[3]
					substituted = true
				}
			}
		}
	}

	return
}
