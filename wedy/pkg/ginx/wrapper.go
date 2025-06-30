package ginx

import (
	"GoProj/wedy/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var l logger.LoggerV1 = logger.NewNopLogger()
var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
}
