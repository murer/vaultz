package lock

import (
	"encoding/json"

	"github.com/murer/vaultz/util"
)

func stringify(o map[string]string) string {
	ret, err := json.Marshal(o)
	util.Check(err)
	return string(ret)
}
