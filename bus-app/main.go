package main

import (
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("TDS_API_URL")
	log.Printf("URL = %s", url)

	tc := TdsWebClient{
		url:     os.Getenv("TDS_API_URL"),
		key:     os.Getenv("TDS_API_KEY"),
		carrier: "304",
		resty:   resty.New(),
	}

	err = tc.Origins()

	log.Println(err)
}
