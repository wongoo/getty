package hello

import (
	"github.com/dubbogo/getty"
	"time"
)

const (
	CronPeriod      = 20 * time.Second
	WritePkgTimeout = 1e8
)

var (
	log = getty.GetLogger()
)
