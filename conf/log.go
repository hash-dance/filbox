package conf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"gitee.com/szxjyt/filbox-backend/modules/common"
)

const (
	LogsText = "text"
	LogsJSON = "json"
)

// InitLogs set log format, output or logLevel
func InitLogs(config *Config) {
	logrus.SetReportCaller(true)

	switch config.LogFormat {
	case LogsText:
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.Debugf("set logs format %s", LogsText)
	case LogsJSON:
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.Debugf("set logs format %s", LogsJSON)
	}

	// 输出到日志目录和标准输出
	if config.LogDispatch {
		logrus.SetOutput(ioutil.Discard)
		logrus.AddHook(&WriterHook{
			Writer: MultiWriter(path.Join(config.LogPath, "error")),
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
			},
		})
		logrus.AddHook(&WriterHook{ // Send info and debug logs to stdout
			Writer: MultiWriter(path.Join(config.LogPath, "info")),
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
			},
		})
	} else {
		fileout, err := ConfigLocalFilesystemLogger(path.Join(config.LogPath, "log"))
		if err != nil {
			logrus.Errorf(err.Error())
		} else {
			mw := io.MultiWriter(os.Stdout, fileout)
			logrus.SetOutput(mw)
		}
	}

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Open Debug Level")
	}

	logrus.Debugf("logsFormat is: %+v", config.LogFormat)
}

// ConfigLocalFilesystemLogger 切割日志和清理过期日志
func ConfigLocalFilesystemLogger(filePath string) (io.Writer, error) {
	abs, _ := filepath.Abs(filePath)
	logrus.Infof(abs)
	return rotatelogs.New(
		abs+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(abs),           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*7),  // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
	)
}

// MultiWriter write to multi output
func MultiWriter(logPath string) io.Writer {
	mw := io.Writer(os.Stdout)
	fileout, err := ConfigLocalFilesystemLogger(logPath)
	if err != nil {
		logrus.Errorf(err.Error())
	} else {
		mw = io.MultiWriter(os.Stdout, fileout)
	}
	return mw
}

// RequestLogger logger for http response
func RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// 获取reqID
			reqID := middleware.GetReqID(r.Context())
			// 注入reqID到日志
			// logrus.AddHook(&RequestIDHook{
			// 	RequestID: reqID,
			// })

			// 获取 response body
			buffer := new(bytes.Buffer)
			ww := common.NewWrapResponseWriter(w, r.ProtoMajor, buffer)

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			requestLog := fmt.Sprintf("%s %s://%s%s %s from %s - ",
				r.Method, scheme, r.Host, r.RequestURI, r.Proto, r.RemoteAddr)

			requestBody := make([]byte, 0)
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				requestBody, _ = ioutil.ReadAll(r.Body)
				defer func() {
					if err := r.Body.Close(); err != nil {
						return
					}
				}()
				r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
			}

			t1 := time.Now()
			defer func() {
				responseLog := fmt.Sprintf("%03d %dB in %s", ww.Status(), ww.BytesWritten(), time.Since(t1))
				entry := logrus.WithFields(logrus.Fields{
					"responseBody": buffer.String(),
					"request_id":   reqID,
				})
				// 下面请求解析requestbody到日志
				if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
					if strings.Contains(r.Header.Get("Content-Type"), "json") {
						entry.Data["requestBody"] = string(requestBody)
					} else {
						entry.Data["requestBody"] = r.ContentLength
					}
				}
				entry.Info(requestLog + responseLog)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// 通过hook自定义日志字段,加入requestID
type RequestIDHook struct {
	RequestID string `json:"request_id"`
}

func (hook *RequestIDHook) Fire(entry *logrus.Entry) error {
	entry.Data["request_id"] = "[" + hook.RequestID + "]"
	return nil
}

func (hook *RequestIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}
