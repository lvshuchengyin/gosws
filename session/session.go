// session
package session

import "net/http"

type Session interface {
	ID() string
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string) error
	Save(w http.ResponseWriter) error
	Clean() error
}

type SessionFactory interface {
	Init(lifetime int64, secretKey, jsonConf string) error
	Get(sid string) Session
	Exist(sid string) bool
	Regenerate(oldsid, sid string) (Session, error)
	Destroy(sid string) error
	Count() int64
	GC() error
}
