package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NordSecurity/nordvpn-linux/daemon/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DaemonURL = "unix:///run/nordvpn/nordvpnd.sock"

var Empty = &pb.Empty{}

func main() {
	conn, err := grpc.Dial(
		DaemonURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := pb.NewDaemonClient(conn)
	resp, err := client.Status(ctx, Empty)
	if err != nil {
		panic(err)
	}

	fmt.Println(Status(resp))

	payload, err := client.Cities(ctx, &pb.CitiesRequest{
		Country: "United_States",
	})
	if err != nil {
		panic(err)
	}
	cities := payload.GetData()
	for _, city := range cities {
		fmt.Println(city)
	}

	_, err = client.Connect(ctx, &pb.ConnectRequest{
		ServerGroup: "",
		ServerTag:   "new_york",
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	resp, err = client.Status(ctx, Empty)
	if err != nil {
		panic(err)
	}
	fmt.Println(Status(resp))
}
