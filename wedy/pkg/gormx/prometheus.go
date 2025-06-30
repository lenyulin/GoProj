package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	vector *prometheus.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus callbacks"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	err := db.Callback().Create().Before("*").Register("prometheus_create_gorm_before", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_gorm_after", c.After("CREATE"))
	if err != nil {
		return err
	}
	err = db.Callback().Query().Before("*").Register("prometheus_query_gorm_before", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Query().After("*").Register("prometheus_query_gorm_after", c.After("QUERY"))
	return err
}

func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	vector := prometheus.NewSummaryVec(opts,
		[]string{"type", "table"})
	prometheus.MustRegister(vector)
	return &Callbacks{
		vector: vector,
	}
}
func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}
func (c *Callbacks) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		start, ok := val.(time.Time)
		if ok {
			duration := time.Since(start).Milliseconds()
			c.vector.WithLabelValues(typ, db.Statement.Table).Observe(float64(duration))
		}
	}
}
