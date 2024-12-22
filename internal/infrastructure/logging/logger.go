package logging

import (
	"github.com/phuslu/log"
)

var Logger log.Logger

func InitLogger() {
	Logger = log.DefaultLogger
}
