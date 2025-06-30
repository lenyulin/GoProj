package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Builder struct {
	NameSpace  string
	Subsystem  string
	Name       string
	InstanceId string
}

func Build() *Builder {
	return &Builder{}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.NameSpace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_time",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		ctx.Next()
	}
}
func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "path", "status_code"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.NameSpace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_resp_time",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{0.5: 0.01, 0.75: 0.01, 0.90: 0.01, 0.99: 0.001, 0.999: 0.0001},
	}, labels)
	prometheus.MustRegister(vector)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			partten := ctx.FullPath()
			status_code := ctx.Writer.Status()
			vector.WithLabelValues(method, partten, strconv.Itoa(status_code)).Observe(float64(duration))
		}()
		ctx.Next()
	}
}
