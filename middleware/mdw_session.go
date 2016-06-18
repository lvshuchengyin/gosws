// session
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/logger"
	"github.com/lvshuchengyin/gosws/session"
	"github.com/lvshuchengyin/gosws/util"
)

const (
	NAME_SESS = "session"
)

func init() {
	Register(NAME_SESS, &MiddlewareSession{})
}

//-------MiddlewareSession----------
type MiddlewareSession struct {
	sessionFactory session.SessionFactory
	secretKey      string
}

func (self *MiddlewareSession) Init(sf session.SessionFactory, secretKey string) {
	self.sessionFactory = sf
	self.secretKey = secretKey
}

func (self *MiddlewareSession) Name() string {
	return NAME_SESS
}

func (self *MiddlewareSession) ProcessRequest(arg *context.Context) error {
	var sessionKey string
	var sessionId string
	exist := false
	if cookie, err := arg.Req.Cookie("SessionKey"); err == nil {
		// cookie exist
		sessionKey = cookie.Value
		sessionId = self.getSessionId(sessionKey, arg)
		exist = self.sessionFactory.Exist(sessionId)
	}

	if !exist {
		// new
		sessionKey = self.createSessionKey(arg)
		sessionId = self.getSessionId(sessionKey, arg)
		// set
		cookie := &http.Cookie{
			Name:     "SessionKey",
			Value:    sessionKey,
			Path:     "/",
			Expires:  time.Now().Add(time.Second * time.Duration(86400*365*10)),
			HttpOnly: true,
		}

		logger.Swsd("new cookie: %s", sessionKey)
		http.SetCookie(arg.Res, cookie)
	}

	arg.Session = self.sessionFactory.Get(sessionId)
	// user
	//AccountGetUser(arg)
	//Swsd("session: %+v", arg.Session)
	return nil
}
func (self *MiddlewareSession) ProcessResponse(arg *context.Context) error {
	_ = arg.Session.Save(arg.Res)
	return nil
}

func (self *MiddlewareSession) createSessionKey(arg *context.Context) string {
	randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1 << 62)
	msg := fmt.Sprintf("%d;%d;%s;%s", time.Now().UnixNano(), randNum, arg.Req.RemoteAddr)
	mac := hmac.New(sha256.New, []byte(self.secretKey))
	mac.Write([]byte(msg))
	macbys := mac.Sum(nil)

	return fmt.Sprintf("%x", macbys)
}

func (self *MiddlewareSession) getSessionId(sessionKey string, arg *context.Context) string {
	sMsg := fmt.Sprintf("%s_%s", sessionKey, arg.Req.UserAgent())

	mac := hmac.New(sha256.New, []byte(self.secretKey))
	mac.Write([]byte(sMsg))
	expectedMac := mac.Sum(nil)
	return util.UrlBase64Encode(expectedMac)
}
