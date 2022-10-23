package config

import (
	"flag"
	"fmt"
	"os"
)

var Version = "0.1.0"
var Transport = "http"
var Address = ":8080"
var Storage = "stupid"
var StoragePath = ""

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
}

func setupFlags() {
	flag.StringVar(&Transport, "transport", Transport, "transport to use")
	flag.StringVar(&Address, "address", Address, "address to listen on")
	flag.StringVar(&Storage, "storage", Storage, "storage to use")
	flag.StringVar(&Version, "version", Version, "version of the application")
	flag.StringVar(&StoragePath, "storage-path", StoragePath, "path to the storage")
}

func String() string {
	return fmt.Sprintf(
		"transport: %s, address: %s, storage: %s, version: %s, storage-path: %s",
		Transport,
		Address,
		Storage,
		Version,
		StoragePath,
	)
}
