package env_mapper

import (
	"os"
)

func ReplaceFileWithEnv(fileEnvs map[string]string) map[string]string {
	ret := make(map[string]string)
	for k, v := range fileEnvs {
		realEnv, ok := os.LookupEnv(k)
		if ok {
			ret[k] = realEnv
		} else {
			ret[k] = v
		}
	}
	return ret
}
