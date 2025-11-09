package httpbox

import (
	"log/slog"
	"net/http"
)

const LevelAccess = slog.Level(1)

type AccessSkipFunc func(r *http.Request) bool
type AccessRequestAttrsFunc func(r *http.Request) []any
type AccessResponseAttrsFunc func(w http.ResponseWriter) []any

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
	logger            *slog.Logger
	level             slog.Level
	skipFunc          AccessSkipFunc
	requestAttrsFunc  AccessRequestAttrsFunc
	responseAttrsFunc AccessResponseAttrsFunc
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

func WithAccessSkipFunc(skipFunc AccessSkipFunc) AccessLogOption {
	return func(options *accessLogConfig) {
		options.skipFunc = skipFunc
	}
}

func WithAccessRequestAttrs(fn AccessRequestAttrsFunc) AccessLogOption {
	return func(options *accessLogConfig) {
		options.requestAttrsFunc = fn
	}
}

func WithAccessResponseAttrs(fn AccessResponseAttrsFunc) AccessLogOption {
	return func(options *accessLogConfig) {
		options.responseAttrsFunc = fn
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

			reqAttrs := []any{
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
			}
			if options.requestAttrsFunc != nil {
				reqAttrs = append(reqAttrs, options.requestAttrsFunc(r)...)
			}
			reqGroup := slog.Group("req", reqAttrs...)

			resAttrs := []any{
				slog.Int("status", arw.statusCode),
				slog.Int("body_size", arw.bodySize),
			}
			if options.responseAttrsFunc != nil {
				resAttrs = append(resAttrs, options.responseAttrsFunc(arw)...)
			}
			resGroup := slog.Group("res", resAttrs...)

			options.logger.Log(r.Context(), options.level, "Access log", reqGroup, resGroup)

			return err
		}
	}
}
