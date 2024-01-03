package luchen

import "time"

const (
	defaultRetries           = 5
	defaultRequestTimeout    = time.Second * 3
	defaultConnectionTimeout = time.Second * 1
	defaultPoolSize          = 50
	defaultMaxPoolSize       = 100
	defaultPoolTTL           = time.Minute
)
