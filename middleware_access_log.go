package httpbox

import (
	"log/slog"
	"net/http"
)

const LevelAccess = slog.Level(1)

type accessResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bodySize   int
}

func newAccessResponseWriter(w http.ResponseWriter) *accessResponseWriter {
	return &accessResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (arw *accessResponseWriter) WriteHeader(statusCode int) {
	arw.statusCode = statusCode
	arw.ResponseWriter.WriteHeader(statusCode)
}

func (arw *accessResponseWriter) Write(b []byte) (int, error) {
	size, err := arw.ResponseWriter.Write(b)
	arw.bodySize += size
	return size, err
}

type accessLogConfig struct {
	logger   *slog.Logger
	level    slog.Level
	skipFunc func(r *http.Request) bool
}

func newAccessLogConfig(opts ...AccessLogOption) accessLogConfig {
	options := accessLogConfig{
		logger: slog.Default(),
		level:  LevelAccess,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type AccessLogOption func(*accessLogConfig)

func WithAccessLogger(logger *slog.Logger) AccessLogOption {
	return func(options *accessLogConfig) {
		options.logger = logger
	}
}

func WithAccessLevel(level slog.Level) AccessLogOption {
	return func(options *accessLogConfig) {
		options.level = level
	}
}

func WithAccessSkipFunc(skipFunc func(r *http.Request) bool) AccessLogOption {
	return func(options *accessLogConfig) {
		options.skipFunc = skipFunc
	}
}

func AccessLogMiddleware(opts ...AccessLogOption) Middleware {
	options := newAccessLogConfig(opts...)

	return func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if options.skipFunc != nil && options.skipFunc(r) {
				return h(w, r)
			}

			arw := newAccessResponseWriter(w)

			err := h(arw, r)

			reqGroup := slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
			)

			resGroup := slog.Group("res",
				slog.Int("status", arw.statusCode),
				slog.Int("body_size", arw.bodySize),
			)

			options.logger.Log(r.Context(), options.level, "Access", reqGroup, resGroup)

			return err
		}
	}
}
