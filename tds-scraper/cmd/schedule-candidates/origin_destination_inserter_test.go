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
	credentials := "/home/ayub/Downloads/NGS-Default-CLI.creds"

	opts := []nats.Option{
		nats.Name("origin_destination_inserter_test"),
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
			Bucket:       "TdsOriginDestinations",
			MaxBytes:     16 * 1024 * 1024,
			MaxValueSize: 512 * 1024,
		})

		if err != nil {
			t.Fatalf("could not create kv %s\n", err)
		}

		origin := Stop{StopUuid: uuid.NewString(), City: "New York"}
		_ = kv.Delete(origin.StopUuid)

		destinations := make([]Stop, 0)
		for range 200 {
			entry, err := kv.Get(origin.StopUuid)
			if err == nil {
				err = json.Unmarshal(entry.Value(), &destinations)
				if err != nil {
					t.Fatalf("could not unmarshal kv %s\n", err)
				}
			}

			id := uuid.New().String()
			destinations = append(destinations, Stop{StopUuid: id, City: "New York"})
			data, _ := json.Marshal(destinations)
			_, err = kv.Put(origin.StopUuid, data)
			if err != nil {
				t.Fatalf("could not update '%s', kv %s\n", origin.StopUuid, err)
			}
		}

		if err != nil {
			t.Fatalf("could not put value: %s\n", err)
		}

		for _, destination := range destinations {
			t.Logf("destination: %s\n", destination.StopUuid)
		}
	})

}
