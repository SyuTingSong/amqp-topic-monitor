package main

import (
	"log"
	"os"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
)

type Conf struct {
	host           string
	port           uint
	login          string
	password       string
	vhost          string
	heartbeat      uint
	withRoutingKey bool
	exchangeName   string
	routingKeys    []string
}

func (cfg *Conf) url() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
		cfg.login,
		cfg.password,
		cfg.host,
		cfg.port,
		cfg.vhost,
	)
}

func configs() (cfg *Conf) {
	cfg = &Conf{}
	const (
		hUsage = "the amqp host"
		pUsage = "the amqp port"
		lUsage = "the login name"
		PUsage = "the password"
		vUsage = "the vhost"
		tUsage = "the name of the topic exchange"
		kUsage = "output the routing key before every message"
	)
	flag.StringVar(&cfg.host, "h", "localhost", hUsage)
	flag.UintVar(&cfg.port, "p", 5672, pUsage)
	flag.StringVar(&cfg.login, "l", "guest", lUsage)
	flag.StringVar(&cfg.password, "P", "guest", PUsage)
	flag.StringVar(&cfg.vhost, "v", "/", vUsage)
	flag.StringVar(&cfg.exchangeName, "t", "log", tUsage)
	flag.BoolVar(&cfg.withRoutingKey, "k", false, kUsage)

	flag.Parse()

	if "" == cfg.exchangeName {
		usage()
	}
	if flag.NArg() == 0 {
		cfg.routingKeys = []string{"#"}
	} else {
		cfg.routingKeys = flag.Args()
	}

	return
}

func main() {
	cfg := configs()

	con, err := amqp.Dial(cfg.url())
	dieOnErr(err)

	ch, err := con.Channel()
	dieOnErr(err)

	err = ch.ExchangeDeclare(
		cfg.exchangeName,
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	dieOnErr(err)

	q, err := ch.QueueDeclare(
		"monitor."+cfg.exchangeName+"."+randomChars(10),
		false,
		false,
		true,
		false,
		nil,
	)
	dieOnErr(err)

	for _, key := range cfg.routingKeys {
		ch.QueueBind(q.Name, key, cfg.exchangeName, false, nil)
	}

	d, err := ch.Consume(
		q.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	dieOnErr(err)

	for msg := range d {
		if cfg.withRoutingKey {
			fmt.Println(msg.RoutingKey)
		}
		fmt.Print(string(msg.Body))
		if msg.Body[len(msg.Body)-1] != '\n' {
			fmt.Print("\n")
		}
	}
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"Usage: %s -t <exchange_name> [options] [<routing_key>...]\n",
		os.Args[0],
	)
	flag.PrintDefaults()
	os.Exit(1)
}

func randomChars(length uint) string {
	var i uint
	var base = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, length)
	baseSize := len(base)
	for i = 0; i < length; i++ {
		b[i] = base[rand.Intn(baseSize)]
	}
	return string(b)
}
