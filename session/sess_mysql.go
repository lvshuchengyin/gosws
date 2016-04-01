// mysql
package session

import (
	"database/sql"
	"encoding/json"
	"gosws/db"
	"gosws/logger"
	"net/http"
	"time"
)

func init() {
	Register("mysql", &SessionFactoryMysql{})
}

func initMysqlSessionTable() error {
	sql := `CREATE TABLE sws_session (
		session_key char(64) NOT NULL,
		session_data blob NOT NULL,
		session_expiry int(11) unsigned NOT NULL,
		PRIMARY KEY (session_key)
		) CHARSET=utf8;`

	_, err := db.DBExec(sql)
	if err == nil {
		logger.Info("CREATE TABLE sws_session")
	}
	return nil
}

type SessionMysql struct {
	sid       string
	values    map[string]interface{}
	secretKey string
}

func (self *SessionMysql) ID() string {
	return self.sid
}

func (self *SessionMysql) Set(key string, value interface{}) error {
	self.values[key] = value
	return nil
}

func (self *SessionMysql) Get(key string) interface{} {
	if v, ok := self.values[key]; ok {
		return v
	} else {
		return nil
	}
}

func (self *SessionMysql) Delete(key string) error {
	delete(self.values, key)
	return nil
}

func (self *SessionMysql) Clean() error {
	self.values = make(map[string]interface{})
	return nil
}

func (self *SessionMysql) Save(w http.ResponseWriter) error {
	b, err := SessionEncode(self.secretKey, self.values)
	if err != nil {
		logger.Error("SessionEncode fail, err:%v", err)
		return err
	}

	_, err = db.DBExec("UPDATE sws_session SET `session_data`=?, `session_expiry`=? where session_key=?",
		b, time.Now().Unix(), self.sid)
	return nil
}

//--------------------------------Factory--------------------------------------
type ConfigMysql struct {
	SessionGCInterval int64 `json:"sessionGCInterval"`
}

type SessionFactoryMysql struct {
	maxlifetime int64
	secretKey   string
	config      *ConfigMysql
}

// new
func NewSessionFactoryMysql() *SessionFactoryMysql {
	m := SessionFactoryMysql{}
	return &m
}

func (self *SessionFactoryMysql) Init(lifetime int64, secretKey, jsonConf string) error {
	var config *ConfigMysql
	err := json.Unmarshal([]byte(jsonConf), &config)
	if err != nil {
		return err
	}

	initMysqlSessionTable()

	self.maxlifetime = lifetime
	self.secretKey = secretKey
	self.config = config

	go self.GC()
	return nil
}

// get mysql session by sid
func (self *SessionFactoryMysql) Get(sid string) (sess Session) {
	var kv map[string]interface{}
	row := db.DBQueryRow("SELECT session_data FROM sws_session WHERE session_key=?", sid)
	var sessiondata []byte
	err := row.Scan(&sessiondata)
	if err == sql.ErrNoRows {
		db.DBExec("INSERT INTO sws_session(`session_key`,`session_data`,`session_expiry`) VALUES(?,?,?)", sid, "", time.Now().Unix())
	}

	if len(sessiondata) == 0 {
		kv = make(map[string]interface{})
	} else {
		kv, err = SessionDecode(self.secretKey, string(sessiondata))
		if err != nil {
			logger.Error("SessionDecode fail, err:%v", err)
			kv = make(map[string]interface{})
		}
	}

	sess = &SessionMysql{sid: sid, values: kv, secretKey: self.secretKey}
	return sess
}

func (self *SessionFactoryMysql) Exist(sid string) bool {
	row := db.DBQueryRow("SELECT session_data FROM sws_session WHERE session_key=?", sid)
	var sessiondata []byte
	err := row.Scan(&sessiondata)
	return !(err == sql.ErrNoRows)
}

func (self *SessionFactoryMysql) Regenerate(oldsid, sid string) (Session, error) {
	row := db.DBQueryRow("SELECT session_data FROM sws_session WHERE session_key=?", oldsid)
	var sessiondata []byte
	err := row.Scan(&sessiondata)
	if err == sql.ErrNoRows {
		_, err = db.DBExec("INSERT INTO sws_session(`session_key`,`session_data`,`session_expiry`) VALUES(?,?,?)", oldsid, "", time.Now().Unix())
		if err != nil {
			logger.Error("Regenerate insert fail, err:%v", err)
			return nil, err
		}
	}

	_, err = db.DBExec("UPDATE sws_session SET `session_key`=? WHERE session_key=?", sid, oldsid)
	if err != nil {
		logger.Error("Regenerate update fail, err:%v", err)
		return nil, err
	}

	var kv map[string]interface{}
	if len(sessiondata) == 0 {
		kv = make(map[string]interface{})
	} else {
		kv, err = SessionDecode(self.secretKey, string(sessiondata))
		if err != nil {
			logger.Error("SessionDecode fail, err:%v", err)
			kv = make(map[string]interface{})
		}
	}

	rs := &SessionMysql{sid: sid, values: kv, secretKey: self.secretKey}

	return rs, nil
}

func (self *SessionFactoryMysql) Destroy(sid string) error {
	_, err := db.DBExec("DELETE FROM sws_session WHERE session_key=?", sid)
	if err != nil {
		logger.Error("Session Destroy fail, sid:%s, err:%v", sid, err)
	}

	return err
}

func (self *SessionFactoryMysql) GC() (err error) {
	for {
		_, err = db.DBExec("DELETE FROM sws_session WHERE session_expiry < ?", time.Now().Unix()-self.maxlifetime)
		if err != nil {
			logger.Error("Session GC fail, err:%v", err)
		}

		time.Sleep(time.Second * time.Duration(self.config.SessionGCInterval))
	}

	return
}

func (self *SessionFactoryMysql) Count() int64 {
	var total int64
	err := db.DBQueryRow("SELECT COUNT(*) AS num FROM sws_session").Scan(&total)
	if err != nil {
		logger.Error("Session Count fail, err:%v", err)
		return 0
	}

	return total
}
