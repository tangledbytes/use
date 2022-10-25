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
var WorkerID = 4095
var LogLevel = "info"
var DBSyncType = "none"
var DBReadOnly = false

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
	flag.StringVar(
		&LogLevel,
		"log-level",
		utils.GetEnvOrDefault(convertToEnvName("USE", "log-level"), LogLevel),
		"log level",
	)
	flag.StringVar(
		&DBSyncType,
		"db-sync-type",
		utils.GetEnvOrDefault(convertToEnvName("USE", "db-sync-type"), DBSyncType),
		"db sync type",
	)
	flag.BoolVar(
		&DBReadOnly,
		"db-read-only",
		utils.StringToBool(
			utils.GetEnvOrDefault(convertToEnvName("USE", "db-read-only"), utils.BoolToString(DBReadOnly)),
		),
		"db read only",
	)

	flag.Parse()
}

func convertToEnvName(prefix, name string) string {
	name = strings.ReplaceAll(strings.ToUpper(name), "-", "_")

	return prefix + "_" + name
}

func String() string {
	return fmt.Sprintf(
		`
  transport: %s,
  address: %s, 
  storage: %s, 
  version: %s,
  storage-path: %s,
  worker-id: %d,
  log-level: %s
  db-sync-type: %s
  db-read-only: %t`,
		Transport,
		Address,
		Storage,
		Version,
		StoragePath,
		WorkerID,
		LogLevel,
		DBSyncType,
		DBReadOnly,
	)
}
