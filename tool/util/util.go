package util

import (

	"github.com/satori/go.uuid"
	"strings"
)



//生成uuid
func GenerateUuid() string {
	uid := uuid.NewV1()
	uids := strings.Split(uid.String(), "-")
	return uids[0] + uids[1] + uids[2] + uids[4] + uids[3]
}
