package datasource

import "time"

type Cache interface {
	Set(string, interface{}, time.Duration) error
	Get(string) (interface{}, error)
	Clear()
	Invalidate(string)
}
