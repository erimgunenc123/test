package database_logger

import (
	"context"
	"errors"
	"fmt"
	"genericAPI/api/environment"
	"genericAPI/api/os_specific_constants"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"io"
	"log"
	"os"
	"time"
)

type DBLogger struct {
	logger.Config
	logger.Writer
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func (l *DBLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l DBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l DBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l DBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

var dbLogger DBLogger

func InitDbLogger() *DBLogger {
	wd, _ := os.Getwd()
	f, _ := os.OpenFile(fmt.Sprintf("%s%s%s", wd, os_specific_constants.PATH_SEPERATOR, time.Now().Format("db_2006-01-02.log")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	var ioWriter io.Writer
	if environment.IsTestEnvironment() {
		ioWriter = io.MultiWriter(f, os.Stdout)
	} else {
		ioWriter = f
	}
	dbLogger = DBLogger{
		Config: logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			Colorful:                  false,
		},
		Writer:       log.New(ioWriter, "\r\n", log.LstdFlags),
		infoStr:      "%s\n[info] ",
		warnStr:      "%s\n[warn] ",
		errStr:       "%s\n[error] ",
		traceStr:     "%s\n[%.3fms] [rows:%v] %s",
		traceWarnStr: "%s %s\n[%.3fms] [rows:%v] %s",
		traceErrStr:  "%s %s\n[%.3fms] [rows:%v] %s",
	}
	return &dbLogger
}
