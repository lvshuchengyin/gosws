// cookie_session
package cookie_session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lvshuchengyin/gosws/util"
)

const (
	EXPIRE_KEY = "_exp"
)

type CookieSession struct {
	secretKey string
	req       *http.Request
	res       http.ResponseWriter
	expire    int64
	values    map[string]interface{}
}

func NewCookieSession(secretKey string, req *http.Request, res http.ResponseWriter, expire int64) *CookieSession {
	return &CookieSession{
		secretKey: secretKey,
		req:       req,
		res:       res,
		expire:    expire,
	}
}

func (self *CookieSession) Parse() (err error) {
	cookie, err := self.req.Cookie("Sws")
	if err != nil {
		return nil
	}

	defer func() {
		if err != nil {
			self.values = map[string]interface{}{}
			self.Flush()
		}
	}()

	vals := strings.Split(cookie.Value, "|")
	if len(vals) != 2 {
		return fmt.Errorf("invalid sws cookie value: %s", cookie.Value)
	}

	content := vals[0]
	sign := vals[1]

	if sign != self.createSignature(content) {
		return fmt.Errorf("sws cookie signature not match")
	}

	originData, err := self.decrypt(content)
	if err != nil {
		return fmt.Errorf("sws cookie decrypt err: %s", err.Error())
	}

	err = json.Unmarshal([]byte(originData), &self.values)
	if err != nil {
		return fmt.Errorf("sws cookie not a json, err: %s", err.Error())
	}

	// check expire
	expire, ok := self.values[EXPIRE_KEY]
	if !ok {
		return fmt.Errorf("sws cookie not found expire")
	}

	fexpire, ok := expire.(float64)
	if !ok {
		return fmt.Errorf("sws cookie expire not float")
	}

	if time.Now().Unix() > int64(fexpire) {
		return fmt.Errorf("sws cookie already expire")
	}

	return nil
}

func (self *CookieSession) GetInt(key string) (v int64) {
	data, ok := self.values[key]
	if !ok {
		return
	}

	fv, ok := data.(float64)
	if !ok {
		return
	}

	return int64(fv)
}

func (self *CookieSession) GetString(key string) (v string) {
	data, ok := self.values[key]
	if !ok {
		return
	}

	v, _ = data.(string)
	return
}

func (self *CookieSession) Set(key string, val interface{}) {
	self.set(key, val)

	self.Flush()
	return
}

func (self *CookieSession) Del(key string) {
	if self.values == nil {
		return
	}

	delete(self.values, key)
	self.Flush()
	return
}

func (self *CookieSession) Flush() (err error) {
	self.set(EXPIRE_KEY, time.Now().Unix()+self.expire)

	expire := time.Now().Add(time.Second * time.Duration(self.expire))
	sessVal := ""

	if len(self.values) > 0 {

		var bs []byte
		bs, err = json.Marshal(self.values)
		if err != nil {
			return
		}
		content := string(bs)
		content, err = self.encrypt(content)
		if err != nil {
			return
		}

		sign := self.createSignature(content)

		sessVal = fmt.Sprintf("%s|%s", content, sign)
	} else {
		expire = time.Now()
	}

	cookie := &http.Cookie{
		Name:     "Sws",
		Value:    sessVal,
		Path:     "/",
		Expires:  expire,
		HttpOnly: true,
	}

	http.SetCookie(self.res, cookie)
	return
}

func (self *CookieSession) set(key string, val interface{}) {
	if self.values == nil {
		self.values = make(map[string]interface{}, 1)
	}
	self.values[key] = val
	return
}

func (self *CookieSession) createSignature(content string) string {
	sMsg := fmt.Sprintf("%s_%s", content, self.req.UserAgent())

	mac := hmac.New(sha256.New, []byte(self.secretKey))
	mac.Write([]byte(sMsg))
	expectedMac := mac.Sum(nil)
	return strings.TrimRight(util.UrlBase64Encode(expectedMac), "=")
}

func (self *CookieSession) encrypt(content string) (string, error) {
	cryptData, err := util.AesEncrypt([]byte(content), []byte(self.secretKey))
	if err != nil {
		return "", err
	}

	encoded := util.UrlBase64Encode(cryptData)
	return strings.TrimRight(string(encoded), "="), nil
}

func (self *CookieSession) decrypt(content string) (string, error) {
	data64, err := util.UrlBase64Decode([]byte(content))
	if err != nil {
		return "", err
	}

	originData, err := util.AesDecrypt(data64, []byte(self.secretKey))
	if err != nil {
		return "", err
	}

	return string(originData), nil
}
