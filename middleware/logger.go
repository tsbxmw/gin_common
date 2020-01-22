package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	common "github.com/tsbxmw/gin_common"
	"io/ioutil"
	"os"
	"time"
)

func LoggerInit(e *gin.Engine, file string) {
	common.InitLogger()
	logger, err := loggerMiddleware(file)
	if err != nil {
		panic(err)
	}
	e.Use(logger)
	loggerError, err := loggerErrorMiddleware(file)
	if err != nil {
		panic(err)
	}
	e.Use(loggerError)
	// here using gin.Recovery as Exception Capture
	//loggerRecovery, err := loggerRecoveryMiddleware(file)
	//if err != nil {
	//	panic(err)
	//}
	//e.Use(loggerRecovery)
}

func loggerMiddleware(file string) (logger gin.HandlerFunc, err error) {
	logClient := common.LogrusLogger

	if _, err = os.Stat(file); err != nil {
		logrus.Info(err)
		if _, err = os.Create(file); err != nil {
			panic(err)
		}
	}
	src, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	logClient.Out = src
	logClient.SetLevel(logrus.DebugLevel)
	apiLogPath := file

	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath+"-temp"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel: logWriter,
		logrus.WarnLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)

	stdOut := lfshook.NewHook(os.Stdout, &logrus.TextFormatter{})
	logClient.AddHook(stdOut)

	return func(c *gin.Context) {
		reqBody, _ := c.GetRawData()
		// here put the body back
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		path := c.Request.URL.Path
		params := c.Request.Header
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		//reqBody = c.Request.Body

		logClient.WithFields(logrus.Fields{
			"status_coude": statusCode,
			"latency":      latency,
			"client_ip":    clientIP,
			"method":       method,
			"path":         path,
			"params":       params,
			"req_body":     string(reqBody),
		}).Info()
	}, err
}

func loggerErrorMiddleware(file string) (logger gin.HandlerFunc, err error) {
	logClient := common.LogrusLogger

	if _, err = os.Stat(file); err != nil {
		if _, err = os.Create(file); err != nil {
			panic(err)
		}
	}
	src, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	logClient.Out = src
	logClient.SetLevel(logrus.DebugLevel)
	apiErrorLogPath := file + "-error"

	logErrorWriter, err := rotatelogs.New(
		apiErrorLogPath+".%Y-%m-%d-%H-%M.error.log",
		rotatelogs.WithLinkName(apiErrorLogPath+"-temp"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeErrorMap := lfshook.WriterMap{
		logrus.WarnLevel:  logErrorWriter,
		logrus.ErrorLevel: logErrorWriter,
		logrus.FatalLevel: logErrorWriter,
	}

	errorHook := lfshook.NewHook(writeErrorMap, &logrus.JSONFormatter{})
	logClient.AddHook(errorHook)

	gin.RecoveryWithWriter(logErrorWriter)

	return func(c *gin.Context) {
		if c.Writer.Status() == 200 {
			c.Next()
			return
		}
		reqBody, _ := c.GetRawData()
		// here put the body back
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		path := c.Request.URL.Path
		params := c.Request.Header
		//reqBody := c.Request.Body
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorsInfo := c.Errors
		logClient.WithFields(logrus.Fields{
			"status_coude": statusCode,
			"latency":      latency,
			"client_ip":    clientIP,
			"method":       method,
			"path":         path,
			"params":       params,
			"req_body":     string(reqBody),
			"error":        errorsInfo,
		}).Error()
	}, err
}

func loggerRecoveryMiddleware(file string) (logger gin.HandlerFunc, err error) {
	logClient := logrus.New()

	if _, err = os.Stat(file); err != nil {
		if _, err = os.Create(file); err != nil {
			panic(err)
		}
	}
	src, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	logClient.Out = src
	logClient.SetLevel(logrus.DebugLevel)
	apiRecoveryLogPath := file + "-recovery"

	logRecoveryWriter, err := rotatelogs.New(
		apiRecoveryLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiRecoveryLogPath+"-temp"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	return gin.RecoveryWithWriter(logRecoveryWriter), nil
}
