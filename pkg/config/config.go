package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var Version = "0.1.0"
var Transport = "http"
var Address = ":8080"
var Storage = "stupid"
var StoragePath = ""
var WorkerID = 0

func Setup() {
	setupEnvs()
	setupFlags()
}

func setupEnvs() {
	if transport := os.Getenv("USE_TRANSPORT"); transport != "" {
		Transport = transport
	}

	if address := os.Getenv("USE_ADDRESS"); address != "" {
		Address = address
	}

	if storage := os.Getenv("USE_STORAGE"); storage != "" {
		Storage = storage
	}

	if version := os.Getenv("USE_VERSION"); version != "" {
		Version = version
	}

	if storagePath := os.Getenv("USE_STORAGE_PATH"); storagePath != "" {
		StoragePath = storagePath
	}

	if workerIDStr := os.Getenv("USE_WORKER_ID"); workerIDStr != "" {
		if workerID, err := strconv.Atoi(workerIDStr); err == nil {
			WorkerID = workerID
		} else {
			panic(err)
		}
	}
}

func setupFlags() {
	flag.StringVar(&Transport, "transport", Transport, "transport to use")
	flag.StringVar(&Address, "address", Address, "address to listen on")
	flag.StringVar(&Storage, "storage", Storage, "storage to use")
	flag.StringVar(&Version, "version", Version, "version of the application")
	flag.StringVar(&StoragePath, "storage-path", StoragePath, "path to the storage")
	flag.IntVar(&WorkerID, "worker-id", WorkerID, "worker id")
	flag.Parse()
}

func String() string {
	return fmt.Sprintf(
		"transport: %s, address: %s, storage: %s, version: %s, storage-path: %s, worker-id: %d",
		Transport,
		Address,
		Storage,
		Version,
		StoragePath,
		WorkerID,
	)
}
