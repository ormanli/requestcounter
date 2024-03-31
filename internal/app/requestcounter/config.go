package requestcounter

import (
	"os"
	"time"
)

// Config defines configuration of application. Values are parsed from environment variables.
type Config struct {
	// ServerPort is the port number on which the server will listen for incoming requests. Defaults to 11111 if not specified.
	ServerPort int `split_words:"true" default:"11111"`
	// ServerHost is the hostname or IP address of the server. If not specified, the hostname of the machine running the application is used.
	ServerHost string `split_words:"true"`
	// ServerReadHeaderTimeout is the amount of time to wait for a client to send a request header before giving up and closing the connection. Defaults to 1 second if not specified.
	ServerReadHeaderTimeout time.Duration `split_words:"true" default:"1s"`
	// InitDebug is the flag indicating whether debug mode is enabled. If set to true, additional debugging information will be logged during application startup.
	InitDebug bool `split_words:"true"`
	// RedisHost is the hostname or IP address of the Redis server. Defaults to "localhost" if not specified.
	RedisHost string `split_words:"true" default:"localhost"`
	// RedisPort is the port number on which the Redis server is listening for incoming connections. Defaults to 6379 if not specified.
	RedisPort int `split_words:"true" default:"6379"`
	// RedisRetryDuration is the amount of time to wait before retrying a failed Redis operation. Defaults to 1 second if not specified.
	RedisRetryDuration time.Duration `split_words:"true" default:"1s"`
	// RedisMaxRetry is the maximum number of times to retry a failed Redis operation before giving up and returning an error. Defaults to 10 if not specified.
	RedisMaxRetry int `split_words:"true" default:"10"`
	// RedisCounterKey is the key used to store the counter value in Redis. This field is required.
	RedisCounterKey string `split_words:"true" required:"true"`
}

// GetServerHost returns the hostname or IP address of the server. If the ServerHost field is not specified, it will return the hostname of the machine running the application.
func (c Config) GetServerHost() string {
	if c.ServerHost == "" {
		hostname, _ := os.Hostname()
		return hostname
	}

	return c.ServerHost
}
