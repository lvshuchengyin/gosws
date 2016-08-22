// logger
package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SWSD = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
)

var (
	loggerObj *Logger
)

type Logger struct {
	time           time.Time
	file           *os.File
	level          int
	logpath        string
	toStd          bool
	logdays        int
	msgChan        chan string
	maxlogfilenum  int64
	maxlogfilesize int64
	logobj         *log.Logger
	nowfilesize    int64
	lock           sync.Mutex
}

func logLevel(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "swsd":
		return SWSD
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "critical":
		return CRITICAL
	default:
		return DEBUG
	}
}

func Init(slevel, logpath string, isStd int64, maxlogfilenum, maxlogfilesize int64) *Logger {
	level := logLevel(slevel)
	toStd := false
	if isStd > 0 {
		toStd = true
	}

	log := &Logger{
		time:           time.Now(),
		level:          level,
		logpath:        logpath,
		toStd:          toStd,
		maxlogfilenum:  maxlogfilenum,
		maxlogfilesize: maxlogfilesize,
	}

	log.delExpire()
	log.openNewLog()

	loggerObj = log

	return log
}

func (self *Logger) isFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func (self *Logger) makesureFileDirExist(filePath string) {
	dirPath := path.Dir(filePath)
	b, _ := self.isFileExist(dirPath)
	if b {
		return
	}

	os.MkdirAll(dirPath, os.ModeDir)
}

func (self *Logger) getDirFiles(dirpath string) (fps []string, err error) {
	dir, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}

	for _, info := range dir {
		if info.IsDir() {
			continue
		}
		filepath := dirpath + "/" + info.Name()
		fps = append(fps, filepath)
	}
	return
}

func (self *Logger) delExpire() {
	dirPath := path.Dir(self.logpath)
	fps, err := self.getDirFiles(dirPath)
	if err != nil {
		fmt.Println("delExpireLog err:", err)
		return
	}

	logname := path.Base(self.logpath)
	for _, fp := range fps {
		fpbn := path.Base(fp)
		if !strings.HasPrefix(fpbn, logname) {
			fmt.Println(fp, logname)
			continue
		}

		items := strings.Split(fp, ".")
		if len(items) < 2 {
			continue
		}

		num, err := strconv.ParseInt(items[len(items)-1], 10, 64)
		if err != nil {
			continue
		}

		if num > self.maxlogfilenum {
			err := os.Remove(fp)
			if err != nil {
				fmt.Println("remove log err:", err)
			}
		}
	}
}

func (self *Logger) rotate() {
	for i := self.maxlogfilenum + 1; i > 0; i-- {
		o := fmt.Sprintf("%s.%d", self.logpath, i)
		ok, erre := self.isFileExist(o)
		if !ok || erre != nil {
			continue
		}

		n := fmt.Sprintf("%s.%d", self.logpath, i+1)
		err := os.Rename(o, n)
		if err != nil {
			fmt.Println("sws rotate log file err:", err)
		}
	}

	err := os.Rename(self.logpath, fmt.Sprintf("%s.%d", self.logpath, 1))
	if err != nil {
		fmt.Println("sws rotate log file err1:", err)
	}
}

func (self *Logger) check() {
	if self.file != nil && self.nowfilesize < self.maxlogfilesize {
		return
	}
	if self.file != nil {
		self.file.Close()
	}

	self.rotate()
	self.delExpire()
	self.openNewLog()
}

func (self *Logger) openNewLog() {
	self.makesureFileDirExist(self.logpath)
	file, err := os.OpenFile(self.logpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		fmt.Println("OpenLogFile err: ", self.logpath, err.Error())
		return
	}

	self.file = file
	info, err := os.Stat(self.logpath)
	if err != nil {
		self.nowfilesize = int64(0)
	} else {
		self.nowfilesize = info.Size()
	}

	self.logobj = log.New(file, "", 0)
}

func (self *Logger) log(level int, format string, args ...interface{}) {
	if self.level > level {
		return
	}

	msg := getLogMsg(level, format, args...)
	self.writeLog(msg)
}

func (self *Logger) writeLog(msg string) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.nowfilesize += int64(len(msg))
	if self.toStd {
		fmt.Print(msg)
	}

	self.check()
	self.logobj.Print(msg)
}

func getLogMsg(level int, format string, args ...interface{}) string {
	timePrefix := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	_, filename := path.Split(file)
	msg := fmt.Sprintf("%s [%s:%d] %s\n", timePrefix, filename, line, fmt.Sprintf(format, args...))
	return msg
}

func printLog(level int, format string, args ...interface{}) {
	msg := getLogMsg(level, format, args...)
	fmt.Print(msg)
}

func Swsd(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(SWSD, "[S] "+format, args...)
		return
	}

	loggerObj.log(SWSD, "[S] "+format, args...)
}

func Debug(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(DEBUG, "[D] "+format, args...)
		return
	}

	loggerObj.log(DEBUG, "[D] "+format, args...)
}

func Info(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(INFO, "[I] "+format, args...)
		return
	}

	loggerObj.log(INFO, "[I] "+format, args...)
}

func Warn(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(WARN, "[W] "+format, args...)
		return
	}

	loggerObj.log(WARN, "[W] "+format, args...)
}

func Warning(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(WARN, "[W] "+format, args...)
		return
	}

	loggerObj.log(WARN, "[W] "+format, args...)
}

func Error(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(ERROR, "[E] "+format, args...)
		return
	}

	loggerObj.log(ERROR, "[E] "+format, args...)
}

func Critical(format string, args ...interface{}) {
	if loggerObj == nil {
		printLog(CRITICAL, "[C] "+format, args...)
		return
	}

	loggerObj.log(CRITICAL, "[C] "+format, args...)
	panic(fmt.Sprintf(format, args...))
}
