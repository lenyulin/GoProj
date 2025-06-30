package ioc

import (
	"GoProj/wedy/internal/job"
	"GoProj/wedy/internal/service"
	"GoProj/wedy/pkg/logger"
	rclock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService, c *rclock.Client, l logger.LoggerV1) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30, c, l)
}
func InitJobs(l logger.LoggerV1, rjob *job.RankingJob) *cron.Cron {
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "lenyulin",
		Subsystem: "wedy",
		Name:      "cron_job",
		Help:      "cron job info",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1min", builder.Build(rjob))
	if err != nil {
		panic(err)
	}
	return expr
}
