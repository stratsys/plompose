package labels

import (
	"regexp"
	"strings"
)

var (
	portRe     = regexp.MustCompile("(traefik|loadbalancer)\\.(([a-z]*)\\.)?port")
	lbRe       = regexp.MustCompile("loadbalancer\\.(([a-z]*)\\.)?hostregexp")
	trRe       = regexp.MustCompile("traefik\\.(([a-z]*)\\.)?frontend.rule")
	labelRe    = regexp.MustCompile("{[a-z]+:[^}]+}")
	hostregexp = "hostregexp:"
)

func GetPorts(labels map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range labels {
		matches := portRe.FindStringSubmatch(k)
		if len(matches) == 4 {
			result[matches[3]] = v
		}
	}

	return result
}

func GetDomains(labels map[string]string) map[string][]string {
	result := make(map[string][]string)
	for k, v := range labels {
		if matches := lbRe.FindStringSubmatch(k); len(matches) == 3 {
			match := matches[2]
			result[match] = append(result[match], v)
		} else if matches := trRe.FindStringSubmatch(k); len(matches) == 3 && strings.HasPrefix(strings.ToLower(v), hostregexp) {
			s := v[len(hostregexp):]
			match := matches[2]
			for _, val := range strings.Split(s, ",") {
				val = removeRegexpLabels(val)
				if val != ".*" {
					result[match] = append(result[match], val)
				}
			}
		}
	}

	return result
}

func removeRegexpLabels(input string) string {
	return labelRe.ReplaceAllStringFunc(input, func(s string) string {
		ss := strings.TrimRight(strings.SplitN(s, ":", 2)[1], "}")
		return ss
	})
}
