package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsKVODInserter struct {
	js jetstream.JetStream
}

func (n NatsKVODInserter) Insert(ctx context.Context, origin Stop, destination Stop) error {
	panic("implement me")
}

func NewODInserter(nc *nats.Conn) NatsKVODInserter {
	js, _ := jetstream.New(nc)
	return NatsKVODInserter{js: js}
}
