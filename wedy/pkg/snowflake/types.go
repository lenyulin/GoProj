package snowflake

type SnowFlakeGenerater interface {
	Generate() (int64, error)
}
