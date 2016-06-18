// logtrace
package logger

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogTrace struct {
	index int
	buf   bytes.Buffer

	startTime    int64
	subName      string
	subStartTime int64
	subStartLine int
}

func NewLogTrace() *LogTrace {
	return &LogTrace{startTime: time.Now().UnixNano()}
}

func (self *LogTrace) Moudle(module string) {
	self.startTime = time.Now().UnixNano()

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = path.Base(file)
	}
	msg := fmt.Sprintf("[%s:%s:%d],", module, file, line)
	self.buf.WriteString(msg)
}

func (self *LogTrace) Start(m string) {
	self.index += 1
	self.subName = m
	self.subStartTime = time.Now().UnixNano()
	var ok bool
	_, _, self.subStartLine, ok = runtime.Caller(1)
	if !ok {
		self.subStartLine = 0
	}
}

func (self *LogTrace) End(status string) {
	costtime := (time.Now().UnixNano() - self.subStartTime) / 1000000
	msg := fmt.Sprintf("%d[%s:%d:%dms:%s],", self.index, self.subName, self.subStartLine, costtime, status)
	self.buf.WriteString(msg)
}

func (self *LogTrace) String() string {
	costtime := (time.Now().UnixNano() - self.startTime) / 1000000
	self.buf.WriteString(fmt.Sprintf("cost:%dms", costtime))
	return self.buf.String()
}
