package job

import (
	"GoProj/wedy/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	l      logger.LoggerV1
	vector *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.LoggerV1, opts prometheus.SummaryOpts) *CronJobBuilder {
	return &CronJobBuilder{l: l, vector: prometheus.NewSummaryVec(opts, []string{"jobs", "success"})}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cron.FuncJob(func() {
		start := time.Now()
		b.l.Debug("Job start running", logger.String("name", name))
		err := job.Run()
		if err != nil {
			b.l.Error("Job Run Error",
				logger.Error(err),
				logger.String("name", name))
		}
		b.l.Debug("Job stopped", logger.String("name", name))
		duration := time.Since(start)
		b.vector.WithLabelValues(name, strconv.FormatBool(err == nil)).Observe(float64(duration.Milliseconds()))
	})
}
