package main

import (
	"log"
	"os"
)

var (
	hdfsURL    = os.Getenv("HDFS_URL")         // E.g.: `data.center.com:31000`
	hdfsUser   = os.Getenv("HDFS_USER")        // E.g.: `root`
	serverPort = os.Getenv("HTTP_SERVER_PORT") // E.g.: `80`
)

func init() {
	if hdfsURL == "" {
		log.Fatal("Missing HDFS_URL environment variable!")
	}

	if hdfsUser == "" {
		log.Fatal("Missing HDFS_USER environment variable!")
	}

	if serverPort == "" {
		serverPort = "8888"
		log.Println("Port has not been set. Using default 8888.")
	}
}
