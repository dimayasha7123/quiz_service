package repository

import (
	"fmt"
	"strconv"
)

const (
	keyPrefix = "user_"

	fieldName = "name"
	fieldQSID = "qsid"
)

func getKeyByTGID(id int64) string {
	return fmt.Sprintf("%s%d", keyPrefix, id)
}

func getTGIDByKey(key string) (int64, error) {
	noPrefKey := key[len(keyPrefix):]
	id, err := strconv.Atoi(noPrefKey)
	if err != nil {
		return 0, fmt.Errorf("can't convert string tgid to int")
	}
	return int64(id), nil
}
