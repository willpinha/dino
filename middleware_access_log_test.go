package httpbox_test

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/willpinha/httpbox"
)

// mockLogHandler captures log records for testing
type mockLogHandler struct {
	records []mockLogRecord
}

type mockLogRecord struct {
	level   slog.Level
	message string
	attrs   map[string]any
}

func (h *mockLogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *mockLogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]any)
	record.Attrs(func(a slog.Attr) bool {
		h.extractAttrs(attrs, "", a)
		return true
	})

	h.records = append(h.records, mockLogRecord{
		level:   record.Level,
		message: record.Message,
		attrs:   attrs,
	})
	return nil
}

func (h *mockLogHandler) extractAttrs(attrs map[string]any, prefix string, attr slog.Attr) {
	if attr.Value.Kind() == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		groupPrefix := attr.Key
		if prefix != "" {
			groupPrefix = prefix + "." + attr.Key
		}
		for _, a := range groupAttrs {
			h.extractAttrs(attrs, groupPrefix, a)
		}
	} else {
		key := attr.Key
		if prefix != "" {
			key = prefix + "." + attr.Key
		}
		attrs[key] = attr.Value.Any()
	}
}

func (h *mockLogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *mockLogHandler) WithGroup(_ string) slog.Handler {
	return h
}

func TestAccessLogMiddleware_DefaultOptions(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	err := handler(w, req)

	assert.NoError(t, err)
	require.Len(t, mockHandler.records, 1)

	record := mockHandler.records[0]

	assert.Equal(t, httpbox.LevelAccess, record.level)
	assert.Equal(t, "Access", record.message)

	// Verify request attributes
	assert.Equal(t, "GET", record.attrs["req.method"])
	assert.Equal(t, "/test?foo=bar", record.attrs["req.url"])
	assert.Equal(t, "192.168.1.1:1234", record.attrs["req.remote_addr"])

	// Verify response attributes
	assert.Equal(t, int64(200), record.attrs["res.status"])
	assert.Equal(t, int64(13), record.attrs["res.body_size"])
}

func TestAccessLogMiddleware_WithCustomLevel(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(
		httpbox.WithAccessLogger(logger),
		httpbox.WithAccessLevel(slog.LevelInfo),
	)

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	require.Len(t, mockHandler.records, 1)
	assert.Equal(t, slog.LevelInfo, mockHandler.records[0].level)
}

func TestAccessLogMiddleware_WithSkipFunc(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	// Skip health check requests
	middleware := httpbox.AccessLogMiddleware(
		httpbox.WithAccessLogger(logger),
		httpbox.WithAccessSkipFunc(func(r *http.Request) bool {
			return r.URL.Path == "/health"
		}),
	)

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})

	// Request to /health should be skipped
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Empty(t, mockHandler.records, "Expected no log records for /health")

	// Request to other path should be logged
	req = httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w = httptest.NewRecorder()
	handler(w, req)

	assert.Len(t, mockHandler.records, 1, "Expected 1 log record for /api/users")
}

func TestAccessLogMiddleware_WithSkipFunc_MultipleConditions(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	// Skip health check and metrics endpoints
	middleware := httpbox.AccessLogMiddleware(
		httpbox.WithAccessLogger(logger),
		httpbox.WithAccessSkipFunc(func(r *http.Request) bool {
			return r.URL.Path == "/health" || r.URL.Path == "/metrics"
		}),
	)

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})

	// Both /health and /metrics should be skipped
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	handler(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	handler(httptest.NewRecorder(), req)

	assert.Empty(t, mockHandler.records, "Expected no log records for skipped paths")

	// Regular request should be logged
	req = httptest.NewRequest(http.MethodGet, "/api/data", nil)
	handler(httptest.NewRecorder(), req)

	assert.Len(t, mockHandler.records, 1, "Expected 1 log record for /api/data")
}

func TestAccessLogMiddleware_CustomStatusCode(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	require.Len(t, mockHandler.records, 1)

	record := mockHandler.records[0]
	assert.Equal(t, int64(404), record.attrs["res.status"])
	assert.Equal(t, int64(9), record.attrs["res.body_size"])
}

func TestAccessLogMiddleware_ErrorPropagation(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))
	expectedErr := errors.New("handler error")

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	w := httptest.NewRecorder()

	err := handler(w, req)

	assert.ErrorIs(t, err, expectedErr)

	// Log should still be recorded even with error
	assert.Len(t, mockHandler.records, 1, "Expected 1 log record even with error")
}

func TestAccessLogMiddleware_NoWriteHeader(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))

	// Handler that doesn't call WriteHeader explicitly
	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("OK"))
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/implicit", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	record := mockHandler.records[0]

	// Should default to 200
	assert.Equal(t, int64(200), record.attrs["res.status"])
}

func TestAccessLogMiddleware_MultipleWrites(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("Hello"))
		w.Write([]byte(", "))
		w.Write([]byte("World!"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/multi", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	record := mockHandler.records[0]

	// Should sum all writes
	assert.Equal(t, int64(13), record.attrs["res.body_size"])
}

func TestAccessLogMiddleware_DifferentMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			mockHandler := &mockLogHandler{}
			logger := slog.New(mockHandler)

			middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))

			handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
				return nil
			})

			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			record := mockHandler.records[0]
			assert.Equal(t, method, record.attrs["req.method"])
		})
	}
}

func TestAccessLogMiddleware_EmptyResponse(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	middleware := httpbox.AccessLogMiddleware(httpbox.WithAccessLogger(logger))

	// Handler that writes nothing
	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusNoContent)
		return nil
	})

	req := httptest.NewRequest(http.MethodDelete, "/resource/123", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	record := mockHandler.records[0]

	assert.Equal(t, int64(0), record.attrs["res.body_size"])
	assert.Equal(t, int64(204), record.attrs["res.status"])
}

func TestAccessLogMiddleware_CombinedOptions(t *testing.T) {
	mockHandler := &mockLogHandler{}
	logger := slog.New(mockHandler)

	// Test combining multiple options
	middleware := httpbox.AccessLogMiddleware(
		httpbox.WithAccessLogger(logger),
		httpbox.WithAccessLevel(slog.LevelWarn),
		httpbox.WithAccessSkipFunc(func(r *http.Request) bool {
			return r.URL.Path == "/internal"
		}),
	)

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})

	// Request to /internal should be skipped
	req := httptest.NewRequest(http.MethodGet, "/internal", nil)
	handler(httptest.NewRecorder(), req)

	assert.Empty(t, mockHandler.records, "Expected no logs for skipped request")

	// Regular request should be logged with custom level
	req = httptest.NewRequest(http.MethodGet, "/api", nil)
	handler(httptest.NewRecorder(), req)

	require.Len(t, mockHandler.records, 1)
	assert.Equal(t, slog.LevelWarn, mockHandler.records[0].level)
}

func TestAccessLogMiddleware_NoOptions(t *testing.T) {
	// Test that middleware works without any options (uses defaults)
	middleware := httpbox.AccessLogMiddleware()

	handler := middleware(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	err := handler(w, req)

	assert.NoError(t, err)

	// Should have written the response
	assert.Equal(t, http.StatusOK, w.Code)
}
