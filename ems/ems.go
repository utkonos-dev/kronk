package ems

type MessageSystem interface {
	Publish(channel string, data []byte) error
	Subscribe(channel string, handler Handler) error
}

type Handler func(data []byte) error