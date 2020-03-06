package main

import (
	"strconv"

	"binlog-db-sync/lib/binlog"

	"binlog-db-sync/lib/errors"
	"binlog-db-sync/lib/helpers"
	"binlog-db-sync/lib/files"
)

func main() {
	rootPath := helpers.GetEnv("ROOT_PATH", ".")
	relativeConfigPath := helpers.GetEnv("CONFIGURATION_PATH", "")
	config, err := files.ReadTwoLevelYALM(rootPath + relativeConfigPath + "/config.yaml")
	errors.CheckAndExitIfError(err, "There is impossible to read configuration file ./config.yaml. The service has been terminated.")

	port, _ := strconv.Atoi(config["dbSettings"]["port"])
	dbSettings := binlog.DbSettings{
		Host:     config["dbSettings"]["host"],
		Port:     port,	
		User:     config["dbSettings"]["user"],
		Password: config["dbSettings"]["password"]}

	handlersSettings, err := files.ReadThreeLevelYALM(rootPath + "/handlers-settings.yaml")
	errors.CheckAndExitIfError(err, "There is impossible to read configuration file ./handlers-settings.yaml. The service has been terminated.")

	handlers := initHandlers(handlersSettings)

	binlog.StartListen(&dbSettings, handlers)
}
