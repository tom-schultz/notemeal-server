package internal

import (
	"log/slog"
	"net/http"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func LogRequestEnd(request *http.Request, status int) {
	if status == 0 {
		status = 200
	}

	slog.Info("END", "method", request.Method, "url", request.URL.Path, "Status", status)
}

func LogRequestStart(request *http.Request) {
	LogRequestMsg("START", request)
}

func LogRequestMsg(msg string, request *http.Request) {
	slog.Info(msg, "method", request.Method, "url", request.URL.Path)
}

func LogRequestError(err error, request *http.Request) {
	slog.Info(err.Error(), "method", request.Method, "url", request.URL.Path)
}
