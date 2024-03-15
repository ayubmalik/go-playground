package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cid, err := strconv.Atoi(os.Getenv("TDS_API_CARRIER_ID"))
	if err != nil {
		log.Fatal("could not parse carrier carrierId")
	}
	tc := TdsRestApi{
		url:       os.Getenv("TDS_API_URL"),
		key:       os.Getenv("TDS_API_KEY"),
		carrierId: cid,
	}

	// origins, err := tc.Origins()
	log.Println("DEST")
	ny := StopCity{StopUuid: "83be15f2-118b-45d9-839c-c92e841f10fdstring"}
	dests, err := tc.Destinations(ny)
	log.Println(dests, err)
}
