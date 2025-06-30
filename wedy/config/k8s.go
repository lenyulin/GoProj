//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "remote_user:Pwd970203..@tcp(42.194.164.163:3306)/wedy",
	},
	Redis: RedisConfig{
		Addr: "wedy-record-redis:6379",
	},
}
