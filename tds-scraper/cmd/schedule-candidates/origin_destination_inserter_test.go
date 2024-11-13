package main_test

import (
	"github.com/nats-io/nats.go"
	"testing"
)

func TestKVSpike(t *testing.T) {
	credentials := "/home/ayub/Downloads/NGS-Default-user_pub_ayub.creds"

	opts := []nats.Option{
		nats.UserCredentials(credentials),
	}

	nc, err := nats.Connect("tls://connect.ngs.global", opts...)
	if err != nil {
		t.Fatalf("could not connect to nats %s\n", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("could not connect to nats %s\n", err)
	}

	t.Cleanup(func() {
		_ = nc.Drain()
	})

	t.Run("put", func(t *testing.T) {
		kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "originDestinations",
		})
		if err != nil {
			t.Fatalf("could not create kv %s\n", err)
		}

		_, err = kv.Put("x", []byte("hello"))
		if err != nil {
			t.Fatalf("could not put value: %s\n", err)
		}
	})

}
