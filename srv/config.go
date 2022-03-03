package srv

type Config interface {
	Set(key string, val interface{})
	Get(key string) (interface{}, error)
}
