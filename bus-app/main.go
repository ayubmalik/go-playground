package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("TDS_API_URL")
	log.Printf("URL = %s", url)

	tc := TdsWebClient{
		url: os.Getenv("TDS_API_URL"),
		key: os.Getenv("TDS_API_KEY"),
	}

	err = tc.Origins2()

	/**
	 * Bal
	 * sjsjs
	 *
	 */

	log.Println(err)
}
