// logtrace
package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

type LogTrace struct {
	index int
	buf   bytes.Buffer

	startTime      int64
	lastRecordTime int64
}

func NewLogTrace() *LogTrace {
	return &LogTrace{
		startTime:      time.Now().UnixNano(),
		lastRecordTime: time.Now().UnixNano(),
	}
}

func (self *LogTrace) SetTrace(name string, status interface{}) {
	self.index += 1

	now := time.Now().UnixNano()
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		line = 0
	}

	costtime := (now - self.lastRecordTime) / 1000000
	msg := fmt.Sprintf("%d[%s:%d:%dms:%v],", self.index, name, line, costtime, status)
	self.buf.WriteString(msg)

	self.lastRecordTime = now
}

func (self *LogTrace) String() string {
	costtime := (time.Now().UnixNano() - self.startTime) / 1000000
	self.buf.WriteString(fmt.Sprintf("cost:%dms", costtime))
	return self.buf.String()
}
