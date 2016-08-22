// httpserver
package gosws

import (
	"runtime"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lvshuchengyin/gosws/config"
	"github.com/lvshuchengyin/gosws/db"
	"github.com/lvshuchengyin/gosws/logger"
)

func Init(confPath string) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// conf
	config.InitConf(confPath)

	//log
	confLog := config.Log()
	logger.Init(confLog.Level, confLog.Path, confLog.Isstd, confLog.Maxfilenum, confLog.Maxfilesize)

	// db
	confDB := config.DB()
	db.InitDBPool(confDB.Source, confDB.Maxidle, confDB.Maxopen)

	return
}
