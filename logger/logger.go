package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func RequestLogger(request *http.Request, loggerName string) *logrus.Entry {
	ipAddress := request.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = request.Header.Get("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = request.RemoteAddr
		}
	}

	redirectUrl := request.Host + request.URL.String()

	return Logger(loggerName).WithFields(logrus.Fields{
		"method":        request.Method,
		"uri":           request.RequestURI,
		"body":          request.Body,
		"form":          request.Form,
		"multipartFrom": request.MultipartForm,
		"postForm":      request.PostForm,
		"redirectUrl":   redirectUrl,
		"ipAddress":     ipAddress,
		"headers":       request.Header,
	})
}
