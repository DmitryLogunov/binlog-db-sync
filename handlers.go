package main

import (
	"binlog-db-sync/lib/binlog"	
	"binlog-db-sync/models"

	"binlog-db-sync/lib/yml"
)

func initHandlers(handlersSettings yml.ThreeLevelMap) *binlog.Handlers {
	schemaTablesModels := map[string]map[string]interface{}{
		"wsdb": map[string]interface{}{
			"customer": &models.WsdbCustomer{},
		},
	}

	handlers := make(binlog.Handlers, 1)
	index := 0
	for schemaName, tablesModels := range schemaTablesModels {
		for tableName, model := range tablesModels {
			handlers[index] = binlog.Handler{
				InvokeType:         handlersSettings[schemaName][tableName]["invokeType"],
				Model:              models.Create(model, models.DBTable{SchemaName: schemaName, TableName: tableName}),
				LambdaFunctionName: handlersSettings[schemaName][tableName]["lambdaFunctionName"],
				Url:                handlersSettings[schemaName][tableName]["url"],
			}
			index += 1
		}
	}

	return &handlers
}
