package env_mapper

import (
	"os"
)

func ReplaceFileWithEnv(fileEnvs map[string]string) (map[string]string, int) {
	ret := make(map[string]string)
	updated := 0
	for k, v := range fileEnvs {
		realEnv, ok := os.LookupEnv(k)
		if ok {
			updated++
			ret[k] = realEnv
		} else {
			ret[k] = v
		}
	}
	return ret, updated
}
