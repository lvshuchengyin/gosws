// config
package config

import (
	"encoding/xml"
	"fmt"
	"gosws/util"
	"io/ioutil"
)

var (
	gConfig *Config
)

type Config struct {
	SecretKey    string     `xml:"secretkey"`
	ListenAddr   string     `xml:"listenaddr"`
	DB           XMLDb      `xml:"db"`
	Log          XMLLog     `xml:"log"`
	TemplatesDir string     `xml:"templatesdir"`
	Session      XMLSession `xml:"session"`
	Middlewares  []string   `xml:"middlewares>middleware"`
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

type XMLSession struct {
	Lifetime int64              `xml:"lifetime"`
	Sessname string             `xml:"sessname"`
	Config   []XMLSessionConfig `xml:"config"`
}

type XMLSessionConfig struct {
	Name       string `xml:"name"`
	Jsonconfig string `xml:"jsonconfig"`
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

func Session() XMLSession {
	return gConfig.Session
}

func Middlewares() []string {
	return gConfig.Middlewares
}

func SecretKey() string {
	return util.MD5Hex([]byte("gosws_" + gConfig.SecretKey))
}

func ListenAddr() string {
	return gConfig.ListenAddr
}

func MiddlewareNames() []string {
	return gConfig.Middlewares
}

func TemplatesDir() string {
	return gConfig.TemplatesDir
}
