// httpserver
package gosws

import (
	"runtime"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lvshuchengyin/gosws/config"
	"github.com/lvshuchengyin/gosws/db"
	"github.com/lvshuchengyin/gosws/logger"
	"github.com/lvshuchengyin/gosws/middleware"
	"github.com/lvshuchengyin/gosws/session"
)

var (
	sessionFactory session.SessionFactory
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

	// session
	session.Init()

	// middleware init
	middleware.Init()

	return
}

////-------------------------account------------------------------
//func AccountGetUser(arg *context.Context) error {
//	ui := arg.Session.Get("sws_uid")
//	if ui == nil {
//		logger.Swsd("session get uid fail")
//		return nil
//	}

//	fuid, ok := ui.(float64)
//	if !ok {
//		return errors.New("sws_uid not float64")
//	}
//	uid := int64(fuid)

//	user := &SwsUser{}

//	var struc interface{}
//	struc = user
//	s := reflect.ValueOf(struc).Elem()
//	leng := s.NumField()
//	onerow := make([]interface{}, leng)
//	for i := 0; i < leng; i++ {
//		onerow[i] = s.Field(i).Addr().Interface()
//	}

//	row := models.DBQueryRow("SELECT * FROM sws_user WHERE id=?", uid)
//	err := row.Scan(onerow...)
//	if err != nil {
//		logger.Error("AccountGetUser: %d, err: %s", uid, err.Error())
//		return err
//	}
//	arg.User = user
//	return nil
//}

//func AccountRegister(arg *context.Context, username, password string) error {
//	if len(username) <= 0 || len(password) <= 0 {
//		return errors.New("username or password can't be empty")
//	}

//	row := models.DBQueryRow("SELECT id FROM sws_user WHERE username=?", username)
//	var id int64
//	err := row.Scan(&id)
//	if err == nil {
//		return errors.New("此用户名已经被注册")
//	}

//	email := ""
//	isDelete := 0
//	permission := 0
//	createTime := time.Now().Unix()
//	modifyTime := time.Now().Unix()
//	lastLogin := int64(0)
//	password = AccountPasswordHash(password)
//	result, err := models.DBExec("INSERT INTO sws_user (username, password, email, is_delete, permission, create_time, modify_time, last_login) VALUES(?,?,?,?,?,?,?,?)", username, password, email, isDelete, permission, createTime, modifyTime, lastLogin)
//	if err != nil {
//		return err
//	}

//	uid, err := result.LastInsertId()
//	if err != nil {
//		return err
//	}

//	aff, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}

//	if aff <= int64(0) {
//		return errors.New("register fail")
//	}

//	logger.Swsd("AccountRegister username: %s, uid: %d", username, uid)

//	return arg.Session.Set("sws_uid", uid)
//}

//func AccountLogin(arg *context.Context, username, password string) error {
//	password = AccountPasswordHash(password)
//	var uid int64
//	row := models.DBQueryRow("SELECT id FROM sws_user WHERE username=? AND password=?", username, password)
//	err := row.Scan(&uid)
//	if err != nil {
//		logger.Info("AccountLogin fail, u: %s, err: %s", username, err.Error())
//		return err
//	}

//	arg.Session.Set("sws_uid", int64(uid))
//	logger.Swsd("AccountLogin %d", uid)

//	models.DBExec("UPDATE sws_user SET last_login=? where id=?", time.Now().Unix(), uid)
//	return nil
//}

//func AccountLogout(arg *context.Context) error {
//	return arg.Session.Delete("sws_uid")
//}

//func AccountPasswordHash(data string) string {
//	mac := hmac.New(sha256.New, []byte(config.SecretKey))
//	mac.Write([]byte(data))
//	mac.Write([]byte(config.SecretKey))
//	macbs := mac.Sum(nil)
//	return fmt.Sprintf("%x", macbs)
//}

//----------------------session------------------------
type SwsUser struct {
	Id         int64
	Username   string
	Password   string
	Email      string
	IsDelete   int
	Permission int
	CreateTime int64
	ModifyTime int64
	LastLogin  int64
}

//func initUserTable() error {
//	sql := `CREATE TABLE sws_user (
//		id bigint NOT NULL AUTO_INCREMENT,
//		username char(32) NOT NULL UNIQUE,
//		password char(64) NOT NULL,
//		email char(32) NOT NULL,
//		is_delete int NOT NULL,
//		permission int NOT NULL,
//		create_time bigint NOT NULL,
//		modify_time bigint NOT NULL,
//		last_login bigint NOT NULL,
//		PRIMARY KEY (id)
//	) CHARSET=utf8;`
//	_, err := models.DBExec(sql)
//	if err == nil {
//		logger.Info("CREATE TABLE sws_user")
//	}
//	return err
//}
