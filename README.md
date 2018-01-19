# amqp-topic-monitor
A command line utility to watch specified amqp topic

# Usage

```text
Usage: amonitor [options] [routing_keys]...

options:

  -h string
        the amqp host (default "localhost")
  -p uint
        the amqp port (default 5672)
  -l string
        the login name (default "guest")
  -P string
        the password (default "guest")
  -t string
        the name of the topic exchange (default "log")
  -v string
        the vhost (default "/")
  -k    output the routing key before every message
```
