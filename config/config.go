// config
package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/lvshuchengyin/gosws/util"
)

var (
	gConfig Config
)

type Config struct {
	SecretKey       string `xml:"secretkey"`
	ListenAddr      string `xml:"listenaddr"`
	DB              XMLDb  `xml:"db"`
	Log             XMLLog `xml:"log"`
	TemplatesDir    string `xml:"templatesdir"`
	StaticDir       string `xml:"staticdir"`
	SessionLifeTime int64  `xml:"sessionlifetime"`
}

type XMLDb struct {
	Source  string `xml:"source"`
	Maxidle int    `xml:"maxidle"`
	Maxopen int    `xml:"maxopen"`
}

type XMLLog struct {
	Level       string `xml:"level"`
	Path        string `xml:"path"`
	Isstd       int64  `xml:"isstd"`
	Maxfilenum  int64  `xml:"maxfilenum"`
	Maxfilesize int64  `xml:"maxfilesize"`
}

func InitConf(confPath string) {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &gConfig)
	if err != nil {
		fmt.Println("invalid config xml")
		panic(err)
	}

	//fmt.Println(config)
}

func Log() XMLLog {
	return gConfig.Log
}

func DB() XMLDb {
	return gConfig.DB
}

func SecretKey() string {
	return util.MD5Hex([]byte("gosws_" + gConfig.SecretKey))
}

func ListenAddr() string {
	return gConfig.ListenAddr
}

func TemplatesDir() string {
	return gConfig.TemplatesDir
}

func StaticDir() string {
	return gConfig.StaticDir
}

func SessionLifeTime() int64 {
	return gConfig.SessionLifeTime
}
