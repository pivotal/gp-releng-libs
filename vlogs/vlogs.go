package vlogs

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
)

var cloudLogger *logging.Logger
var enableFileNameLineNumFuncName bool

func init() {
	cloudLogger = nil
	enableFileNameLineNumFuncName = false
}

func NewCloudLogger(cloudBucket string, cloudLogName string) error {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		ctx := context.Background()
		client, err := logging.NewClient(ctx, cloudBucket)
		if err != nil {
			return err
		}

		cloudLogger = client.Logger(cloudLogName, logging.EntryCountThreshold(10))
		return nil
	}
	return fmt.Errorf("ENV GOOGLE_APPLICATION_CREDENTIALS is not set")
}

func SetUpLocalLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	return nil
}

// Output FileName:LinuNum and function name when debug level <= Warn
func EnableFileNameLineNumFuncNameForLocalLogs(enable bool) {
	enableFileNameLineNumFuncName = enable
}

func Debug(s string, v ...interface{}) {
	message := fmt.Sprintf(s, v...)
	if cloudLogger != nil {
		cloudLogger.Log(logging.Entry{Payload: message, Severity: logging.Debug})
	}
	if enableFileNameLineNumFuncName {
		pc, file, line, ok := runtime.Caller(1)
		uerFunc := runtime.FuncForPC(pc)
		if !ok {
			logrus.Debug(message)
			return
		}
		logrus.WithFields(logrus.Fields{logrus.FieldKeyFile: file + ":" + strconv.Itoa(line), logrus.FieldKeyFunc: uerFunc.Name()}).Debug(message)
	} else {
		logrus.Debug(message)
	}
}

func Info(s string, v ...interface{}) {
	message := fmt.Sprintf(s, v...)
	if cloudLogger != nil {
		cloudLogger.Log(logging.Entry{Payload: message, Severity: logging.Info})
	}
	if enableFileNameLineNumFuncName {
		pc, file, line, ok := runtime.Caller(1)
		uerFunc := runtime.FuncForPC(pc)
		if !ok {
			logrus.Info(message)
			return
		}
		logrus.WithFields(logrus.Fields{logrus.FieldKeyFile: file + ":" + strconv.Itoa(line), logrus.FieldKeyFunc: uerFunc.Name()}).Info(message)
	} else {
		logrus.Info(message)
	}
}

func Warn(s string, v ...interface{}) {
	message := fmt.Sprintf(s, v...)
	if cloudLogger != nil {
		cloudLogger.Log(logging.Entry{Payload: message, Severity: logging.Warning})
	}
	if enableFileNameLineNumFuncName {
		pc, file, line, ok := runtime.Caller(1)
		uerFunc := runtime.FuncForPC(pc)
		if !ok {
			logrus.Warn(message)
			return
		}
		logrus.WithFields(logrus.Fields{logrus.FieldKeyFile: file + ":" + strconv.Itoa(line), logrus.FieldKeyFunc: uerFunc.Name()}).Warn(message)
	} else {
		logrus.Warn(message)
	}
}

func Error(s string, v ...interface{}) {
	message := fmt.Sprintf(s, v...)
	if cloudLogger != nil {
		cloudLogger.Log(logging.Entry{Payload: message, Severity: logging.Error})
		cloudLogger.Flush()
	}
	pc, file, line, ok := runtime.Caller(1)
	uerFunc := runtime.FuncForPC(pc)
	if !ok {
		logrus.Error(message)
		return
	}
	logrus.WithFields(logrus.Fields{logrus.FieldKeyFile: file + ":" + strconv.Itoa(line), logrus.FieldKeyFunc: uerFunc.Name()}).Error(message)
}

func Fatal(s string, v ...interface{}) {
	message := fmt.Sprintf(s, v...)
	if cloudLogger != nil {
		cloudLogger.Log(logging.Entry{Payload: message, Severity: logging.Critical})
		cloudLogger.Flush()
	}
	pc, file, line, ok := runtime.Caller(1)
	uerFunc := runtime.FuncForPC(pc)
	if !ok {
		logrus.Fatal(message)
		return
	}
	logrus.WithFields(logrus.Fields{logrus.FieldKeyFile: file + ":" + strconv.Itoa(line), logrus.FieldKeyFunc: uerFunc.Name()}).Fatal(message)
}
