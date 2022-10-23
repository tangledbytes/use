package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/utkarsh-pro/use/pkg/utils"
)

var Version = "0.1.0"
var Transport = "http"
var Address = ":8080"
var Storage = "stupid"
var StoragePath = ""
var WorkerID = 0

func Setup() {
	setupFlags()
}

func setupFlags() {
	flag.StringVar(
		&Transport,
		"transport",
		utils.GetEnvOrDefault(convertToEnvName("USE", "transport"), Transport),
		"transport to use",
	)
	flag.StringVar(
		&Address,
		"address",
		utils.GetEnvOrDefault(convertToEnvName("USE", "address"), Address),
		"address to listen on",
	)
	flag.StringVar(
		&Storage,
		"storage",
		utils.GetEnvOrDefault(convertToEnvName("USE", "storage"), Storage),
		"storage to use",
	)
	flag.StringVar(
		&StoragePath,
		"storage-path",
		utils.GetEnvOrDefault(convertToEnvName("USE", "storage-path"), StoragePath),
		"path to the storage",
	)
	flag.IntVar(
		&WorkerID,
		"worker-id",
		utils.StringToInt(utils.GetEnvOrDefault(convertToEnvName("USE", "worker-id"),
			utils.IntToString(WorkerID))),
		"worker id",
	)

	flag.Parse()
}

func convertToEnvName(prefix, name string) string {
	name = strings.ReplaceAll(strings.ToUpper(name), "-", "_")

	return prefix + "_" + name
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
