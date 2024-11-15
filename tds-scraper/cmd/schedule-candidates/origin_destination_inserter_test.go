package main_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"testing"
)

type Stop struct {
	StopUuid string `json:"stopUuid"`
	City     string
}

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
			Bucket: "OriginDestinationsX",
		})

		if err != nil {
			t.Fatalf("could not create kv %s\n", err)
		}

		_, err = kv.Put("x", []byte("hello"))

		dests := make([]Stop, 0, 100)
		for range 5 {
			id := uuid.New().String()
			stop := Stop{StopUuid: id, City: "New York"}
			key, _ := json.Marshal(stop)
			entry, err := kv.Get(string(key))
			if err != nil {
				t.Fatalf("could not get kv %s\n", err)
			}

			err = json.Unmarshal(entry.Value(), &dests)
			if err != nil {
				t.Fatalf("could not unmarshal kv %s\n", err)
			}

			t.Logf("dests: %v\n", dests)
		}

		if err != nil {
			t.Fatalf("could not put value: %s\n", err)
		}
	})

}
