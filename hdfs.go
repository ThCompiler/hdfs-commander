package main

import (
	"io"
	"log"

	"github.com/colinmarc/hdfs"
)

var client *hdfs.Client

func init() {
	connect()
}

func connect() {
	var err error

	clientOptions := hdfs.ClientOptions{
		Addresses: []string{hdfsURL},
		User:      hdfsUser,
	}

	client, err = hdfs.NewClient(clientOptions)
	if err != nil {
		log.Fatal("HDFS connect:", err)
	}
}

func downloadFile(path string) (reader io.Reader, err error) {
	reader, err = client.Open(path)

	return
}
