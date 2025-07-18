package ioc

import (
	"GoProj/wedy/interactive/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitConsumer(c1 *events.InteractiveWatchEventConsumer) []events.Consumer {
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
	cfg.Addr = []string{"127.0.0.1:9094"}
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
