package bus

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

var (
  inst *Bus
)

type Bus struct {
  connection *amqp091.Connection
  channel *amqp091.Channel
  queue *amqp091.Queue
  config BusConfig
}

func newBus(config BusConfig) Bus {
  bus := Bus{}
  bus.config = config
  return bus
}

func GetOrCreateBus() *Bus {
  if (inst == nil) {
    dsn := viper.GetString("amqp.dsn")

    b := newBus(BusConfig{
      Connection: BusConnectionConfig{
        Dsn: dsn,
        MaxRetries: 5,
        Timeout: 10 * time.Second,
      },
    })

    inst = &b
  }

  return inst
}

func (b *Bus) Connect() error {
  conn, err := amqp091.Dial(b.config.Connection.Dsn)
  if err != nil {
    return err
  }

  b.connection = conn

  return nil
}

func (b *Bus) Connection() *amqp091.Connection {
  return b.connection
}

func (b *Bus) Channel() (*amqp091.Channel, error) {
  if b.channel == nil {
    c, err := b.connection.Channel()
    if err != nil {
      return nil, err
    }

    b.channel = c
  }

  return b.channel, nil
}

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

