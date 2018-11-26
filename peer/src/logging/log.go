package logging

import (
	"config"
	"strings"

	"github.com/cihub/seelog"
)

var logger seelog.LoggerInterface

type Logger struct {
}

func MustGetLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func (l *Logger) Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
func (l *Logger) Criticalf(format string, args ...interface{}) {
	logger.Criticalf(format, args...)
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
func (l *Logger) Info(args ...interface{}) {
	logger.Info(args)
}
func (l *Logger) Warning(args ...interface{}) {
	logger.Warning(args)
}
func (l *Logger) Error(args ...interface{}) {
	logger.Error(args)
}
func (l *Logger) Critical(args ...interface{}) {
	logger.Critical(args)
}
func (l *Logger) Debug(args ...interface{}) {
	logger.Debug(args)
}

func Initialize() {
	type_ := config.GetLogType()
	levels := config.GetLogLevel()
	maxdays := config.GetLogMaxdays()
	maxfiles := config.GetLogMaxfiles()
	maxsize := config.GetLogMaxsize()

	var useConfig string
	if strings.EqualFold(type_, "date") {
		useConfig = `
<seelog minlevel="debug" maxlevel="critical">
    <outputs formatid="main">
        <filter levels="` + levels + `"> 
            <console/>
        </filter>
		<filter levels="debug,info,warn,critical,error"> 
            <rollingfile type="date" filename="` + config.BasePath + `/logs/log.txt" datepattern="2006.01.02" fullname="true" maxrolls="` + maxdays + `"/> //the unit of size is byte
        </filter>	
    </outputs>
    <formats>
        <format id="main" format="%Date %Time [%LEV] [%File:%Line] [%Func]->%Msg%n"/>
    </formats>
</seelog>
		`
	} else if strings.EqualFold(type_, "size") {
		useConfig = `
<seelog minlevel="debug" maxlevel="critical">
    <outputs formatid="main">
        <filter levels="` + levels + `"> 
            <console/>
        </filter>
		<filter levels="debug,critical,error"> 
            <rollingfile type="size" filename="` + config.BasePath + `/logs/log.txt" maxsize="` + maxsize + `"maxrolls="` + maxfiles + `"/> //the unit of size is byte
        </filter>
    </outputs>
    <formats>
        <format id="main" format="%Date %Time [%LEV] [%File:%Line] [%Func]->%Msg%n"/>
    </formats>
</seelog>
		`
	}

	logger, _ = seelog.LoggerFromConfigAsBytes([]byte(useConfig))
	//defer Logger.Flush()

	seelog.ReplaceLogger(logger)
}
