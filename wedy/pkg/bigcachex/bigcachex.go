package bigcachex

type BigCachex interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Serve() error
	Shutdown() error
}
