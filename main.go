package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/cli/cli/compose/types"
	"github.com/kubernetes/kompose/pkg/kobject"
	"github.com/kubernetes/kompose/pkg/loader/compose"
	"github.com/kubernetes/kompose/pkg/transformer/kubernetes"
	"github.com/stratsys/plompose/env"
	"github.com/stratsys/plompose/labels"
)

const path = `C:\source\Playbook-dev\stack\`

func main() {
	opt := kobject.ConvertOptions{CreateD: true} // create deployment by default if no controller has been set
	entries, _ := ioutil.ReadDir(path)
	var files []string
	for _, entry := range entries {
		if entry.Mode().IsRegular() && filepath.Ext(entry.Name()) == ".yml" {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}

	// files = []string{`C:\source\Playbook-dev\stack\echo.yml`, `C:\source\Playbook-dev\stack\http-debugger.yml`}
	Convert(opt, files)
}

// Convert transforms docker compose or dab file to k8s objects
func Convert(opt kobject.ConvertOptions, files []string) {
	// Get a transformer that maps komposeObject to provider's primitives
	transformer := &kubernetes.Kubernetes{Opt: opt}

	mergedComposeObject := kobject.KomposeObject{ServiceConfigs: make(map[string]kobject.ServiceConfig), Secrets: make(map[string]types.SecretConfig)}

	for _, file := range files {
		envs := readEnvs(file)
		envs["CA_CERTIFICATES_PATH"] = "/etc/ssl/certs/ca-certificates.crt"

		env.Set(envs)

		// loader parses input from file into komposeObject.
		l, err := new(compose.Compose), error(nil)

		komposeObject, err := l.LoadFile([]string{file})
		if err != nil {
			log.Fatalf("loading '%s' failed %v", file, err.Error())
		}

		env.Unset(envs)

		for k, v := range komposeObject.ServiceConfigs {
			mergedComposeObject.ServiceConfigs[k] = v
		}

		for k, v := range komposeObject.Secrets {
			mergedComposeObject.Secrets[k] = v
		}

		mergedComposeObject.LoadedFrom = komposeObject.LoadedFrom
	}

	exposeServices(mergedComposeObject)

	// Do the transformation
	objects, err := transformer.Transform(mergedComposeObject, opt)

	if err != nil {
		log.Fatalf(err.Error())
	}

	os.MkdirAll("output", 0644)
	os.Chdir("output")

	// Print output
	err = kubernetes.PrintList(objects, opt)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func readEnvs(file string) map[string]string {
	path, base := filepath.Dir(file), filepath.Base(file)
	ext := filepath.Ext(base)
	parent := filepath.Dir(path)

	return env.ReadFrom(parent, filepath.Join(path, base[:len(base)-len(ext)]))
}

func exposeServices(komposeObject kobject.KomposeObject) {
	for k, v := range komposeObject.ServiceConfigs {
		ports := labels.GetPorts(v.DeployLabels)
		domains := labels.GetDomains(v.DeployLabels)
		if len(ports) == 1 {
			for key, port := range ports {
				if domain, ok := domains[key]; ok && len(domain) == 1 {
					d := domain[0]
					porti, _ := strconv.Atoi(port)
					if strings.HasPrefix(d, "^") {
						d = d[1:]
					}
					if strings.HasSuffix(d, "\\..*") {
						d = d[:len(d)-4] + ".*"
					} else if strings.HasSuffix(d, "..*") {
						d = d[:len(d)-3] + ".*"
					} else if strings.HasSuffix(d, "\\.") {
						d = d[:len(d)-2] + ".*"
					}
					v.ExposeService = d
					v.Port = append(v.Port, kobject.Ports{HostPort: int32(porti), ContainerPort: int32(porti)})
				}
			}
		}
		komposeObject.ServiceConfigs[k] = v
	}
}
