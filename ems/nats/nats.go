package nats

import (
	"log"

	"github.com/utkonos-dev/kronk/ems"

	"github.com/nats-io/nats.go"
)

type MS struct {
	conn nats.Conn
}

func NewMS(conn nats.Conn) (ems.MessageSystem, error) {
	return MS{conn: conn}, nil
}

func (s MS) Publish(channel string, data []byte) error {
	return s.conn.Publish(channel, data)
}

func (s MS) Subscribe(channel string, handler ems.Handler) error {
	_, err := s.conn.Subscribe(channel, func(msg *nats.Msg) {
		err := handler(msg.Data)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}

	return nil
}
