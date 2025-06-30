package domian

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id         int64
	Cancel     func()
	Executor   string
	Name       string
	Expression string
}

func (j Job) NextTime() time.Time {
	c := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.Expression)
	return s.Next(time.Now())
}
