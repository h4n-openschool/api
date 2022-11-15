// Package bus implements a message bus connection interface which other
// components of the API service can use to send and receive events.
//
// The bus package should only be used by the API service, as it has been
// designed to work best in ephemeral requests that are stateless.
package bus

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

// Store a private instance of the Bus struct which will be retrieved/set when
// the software requests an instance from [Bus.GetOrCreateBus].
var (
  inst *Bus
)

// Bus represents a connection, channel, and queue definition for a message bus.
type Bus struct {
  // connection is a pointer to an AMQP connection
  connection *amqp091.Connection

  // channel is a pointer to an AMQP channel
  channel *amqp091.Channel

  // queue is a pointer to a software-defined AMQP queue
  queue *amqp091.Queue

  // config is the configuration used to bootstrap the bus
  config BusConfig
}

// newBus creates a new Bus instance given a configuration and returns it.
func newBus(config BusConfig) Bus {
  bus := Bus{}
  bus.config = config
  return bus
}

// GetOrCreateBus retrieves or instantiates a new [Bus], returning a pointer to
// it. It should be used whenever an instance of a [Bus] is required.
func GetOrCreateBus() *Bus {
  if (inst == nil) {
    // Get the DSN from viper so we can decouple from arguments but still
    // respect those as well as the configuration file.
    dsn := viper.GetString("amqp.dsn")

    // Create a new instance of the bus using [newBus].
    b := newBus(BusConfig{
      Connection: BusConnectionConfig{
        Dsn: dsn,
        MaxRetries: 5,
        Timeout: 10 * time.Second,
      },
    })

    // Set the application-wide [Bus] instance value to the newly-created bus.
    inst = &b
  }

  // Return the application-wide [Bus] instance.
  return inst
}

// Connect attempts to connect the [Bus] to the AMQP server.
func (b *Bus) Connect() error {
  // Create (Dial) a new AMQP connection.
  conn, err := amqp091.Dial(b.config.Connection.Dsn)
  if err != nil {
    return err
  }

  // Register the new connection.
  b.connection = conn

  return nil
}

// Connection returns an existing connection if it exists, which can be used to
// create channels, queues, and send/receive messages.
func (b *Bus) Connection() *amqp091.Connection {
  return b.connection
}

// Channel gets or creates an AMQP channel, storing it on the Bus instance for
// re-use if a new channel has been created.
func (b *Bus) Channel() (*amqp091.Channel, error) {
  // If no channel is set up, create one
  if b.channel == nil {
    c, err := b.connection.Channel()
    if err != nil {
      return nil, err
    }

    // Store the new channel on the Bus instance.
    b.channel = c
  }

  // Return the current Bus channel.
  return b.channel, nil
}

// ClassesQueue defines and returns the Classes queue.
func (b *Bus) ClassesQueue() (amqp091.Queue, error) {
  if b.queue == nil {
    c, err := b.Channel()
    if err != nil {
      return amqp091.Queue{}, err
    }

    q, err := c.QueueDeclare(
      "classes",
      false,
      false,
      false,
      false,
      nil,
    )

    if err != nil {
      return amqp091.Queue{}, err
    }

    b.queue = &q
  }

  return *b.queue, nil
}

// ChannelQueue calls [Bus.Channel] as well as [Bus.ClassesQueue], returning both
// or an error if the Bus was unable to establish a new channel or define the
// queue.
func (b *Bus) ChannelQueue() (ch *amqp091.Channel, qu amqp091.Queue, err error) {
  ch, err = b.Channel()
  if err != nil {
    return
  }

  qu, err = b.ClassesQueue()
  if err != nil {
    return
  }

  return
}

