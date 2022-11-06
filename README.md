# openschool/classes

The Class microservice of OpenSchool.

## Running

```shell
# Run RabbitMQ (AMQP message broker)
[hayden@hbjy-pc ~]$ docker run -d -p 5672:5672 rabbitmq

# Run the server
[hayden@hbjy-pc ~]$ go run . serve
```
