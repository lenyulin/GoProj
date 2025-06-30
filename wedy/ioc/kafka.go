package ioc

import (
	events2 "GoProj/wedy/interactive/events"
	"GoProj/wedy/internal/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitConsumer(c1 *events2.InteractiveWatchEventConsumer) []events.Consumer {
	return []events.Consumer{c1}
}
func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}
func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string
	}
	var cfg Config
	cfg.Addr = []string{"14.103.175.18:9094"}
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}
