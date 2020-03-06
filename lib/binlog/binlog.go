package binlog

import (
	"os"	
	"fmt"
	"binlog-db-sync/lib/logger"
	"github.com/siddontang/go-mysql/canal"
)

type DbSettings struct {
	Host     string
	Port     int
	User     string
	Password string
}

func getDefaultCanal(dbSettings *DbSettings) (*canal.Canal, error) {
	cfg := canal.NewDefaultConfig()
	
	cfg.Addr = fmt.Sprintf("%s:%d", dbSettings.Host, dbSettings.Port)
	cfg.User = dbSettings.User
	cfg.Password = dbSettings.Password
	cfg.Flavor = "mysql"
	cfg.Dump.ExecutionPath = ""

	return canal.NewCanal(cfg)
}

func StartListen(dbSettings *DbSettings, handlers *Handlers) {
	c, err := getDefaultCanal(dbSettings)
	if err != nil {
		logger.Error("could not create Mysql Canal", nil)
		os.Exit(1)
	}
	
	coords, err := c.GetMasterPos()
	if err != nil {
		logger.Error("could not get start master position of Mysql Binlog", nil)
		os.Exit(1)
	}

	bh := binlogHandler{}
	bh.Handlers = handlers	

	c.SetEventHandler(&bh)
	c.RunFrom(coords)
}
