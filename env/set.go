package env

import "os"

func Set(envs map[string]string) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
}

func Unset(envs map[string]string) {
	for k := range envs {
		os.Unsetenv(k)
	}
}
