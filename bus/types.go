package bus

import "time"

// BusConnectionConfig defines options for managing the bus connection.
type BusConnectionConfig struct {
  // Dsn is the Datasource Name of the AMQP service, formatted as:
  //
  //  amqp://username:password@host:port/vhost
	Dsn        string

  // MaxRetries is the maximum number of times the server should attempt to
  // connect to the AMQP service.
	MaxRetries int

  // Timeout is how long the AMQP client should attempt to connect for before
  // marking the attempt as failed.
	Timeout    time.Duration
}

// BusConfig represents configuration for the [Bus].
type BusConfig struct {
  // Connection defines connection configuration options.
	Connection BusConnectionConfig
}
