package utils

import (
	"github.com/nats-io/go-nats"
	"log"
)

var con *nats.Conn

func GetNatsConnection() *nats.Conn {
	if con != nil {
		return con
	}

	c, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	con = c

	return con
}
