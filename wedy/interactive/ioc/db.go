package ioc

import (
	"GoProj/wedy/config"
	"GoProj/wedy/pkg/gormx"
	"context"
	"fmt"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic("failed to connect database")
	}
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "wedy",
		RefreshInterval: 5,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "lenyulin",
		Subsystem: "wedy",
		Name:      "gorm_db",
		Help:      "GORM DB DATABASE",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{0.5: 0.01, 0.75: 0.01, 0.90: 0.01, 0.99: 0.001, 0.999: 0.0001},
	})
	err = db.Use(cb)
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithDBSystem("wedy")))
	if err != nil {
		panic(err)
	}
	return db
}
func InitMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://14.103.175.18:27017/").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	return client.Database(config.Config.DB.DSN)
}
