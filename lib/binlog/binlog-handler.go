package binlog

import (
	"fmt"
	"runtime/debug"

	awsClient "binlog-db-sync/lib/aws-client"
	"binlog-db-sync/lib/parser"	
	"binlog-db-sync/lib/types"
	"binlog-db-sync/lib/logger"
	"binlog-db-sync/models"
	"binlog-db-sync/lib/helpers"
	"binlog-db-sync/lib/request"

	"github.com/siddontang/go-mysql/canal"
)

type Handler struct {
	InvokeType         string
	Model              *models.BaseModel
	LambdaFunctionName string
	Url                string
}

type Handlers []Handler

type binlogHandler struct {
	canal.DummyEventHandler
	parser.BinlogParser
	Handlers *Handlers
}

func (h *binlogHandler) OnRow(event *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r, " ", string(debug.Stack()))
		}
	}()

	// base value for canal.DeleteAction or canal.InsertAction
	var startEventRowsIndex = 0
	var eventRowsIndexIncrement = 1

	if event.Action == canal.UpdateAction {
		startEventRowsIndex = 1
		eventRowsIndexIncrement = 2
	}

	for eventRowsIndex := startEventRowsIndex; eventRowsIndex < len(event.Rows); eventRowsIndex += eventRowsIndexIncrement {
		eventSchemeTable := event.Table.Schema + "." + event.Table.Name

		for _, handler := range *(h.Handlers) {
			handlerSchemeTable := handler.Model.SchemaName() + "." + handler.Model.TableName()

			if eventSchemeTable != handlerSchemeTable {
				continue
			}
			
			eventID := helpers.RandomString(20)
			logger.Info("new event", map[string]string{"eventID": eventID, "action": event.Action, "table": eventSchemeTable})
						 
			entity := handler.Model
			h.GetBinLogData(entity.Properties, event, eventRowsIndex)
			entityData := entity.BuildJSON()
    
			switch event.Action {
			case canal.UpdateAction:
				oldEntity := handler.Model
				h.GetBinLogData(oldEntity.Properties, event, eventRowsIndex-1)
				logger.Info("event new data", map[string]string{"eventID": eventID, "data": entityData})
				logger.Info("event old data", map[string]string{"eventID": eventID, "data": oldEntity.BuildJSON()})
			case canal.InsertAction:
				logger.Info("event new data", map[string]string{"eventID": eventID, "data": entityData})
			case canal.DeleteAction:
				logger.Info("event new data", map[string]string{"eventID": eventID, "data": entityData})
			default:
				logger.Warning("unknown action", map[string]string{"eventID": eventID})
				continue
			}					

			invokeParams := types.InvokeParams{Action: event.Action, Schema: event.Table.Schema, Table: event.Table.Name, Data: entityData}
			go callHandler(&handler, &invokeParams, eventID)

			break
		}
	}
	return nil
}

func callHandler(h *Handler, params *types.InvokeParams, eventID string) {
	var handlerIdentifier = ""
	logger.Info("running invoke", map[string]interface{}{"eventID": eventID, "type": h.InvokeType, "handler": handlerIdentifier})		
	switch h.InvokeType {
	case "http":
		handlerIdentifier = h.Url
		response, _ := request.Post(h.Url, params)
		logger.Info("invoke response", map[string]interface{}{"eventID": eventID, "response": response})	
	case "awsInvoke":
		handlerIdentifier = h.LambdaFunctionName
		response := awsClient.Invoke(h.LambdaFunctionName, params)	
		logger.Info("invoke response", map[string]interface{}{"eventID": eventID, "response": response})			
	default:
		return
	}
	
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}
