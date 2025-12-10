package dino

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyMiddlewares_NoMiddlewares(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("original handler"))
		return nil
	})

	wrappedHandler := applyMiddlewares(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "original handler", rec.Body.String())
}

func TestApplyMiddlewares_SingleMiddleware(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("handler"))
		return nil
	})

	middleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("before-"))
			err := next(w, r)
			w.Write([]byte("-after"))
			return err
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	assert.Equal(t, "before-handler-after", rec.Body.String())
}

func TestApplyMiddlewares_MultipleMiddlewares(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("[handler]"))
		return nil
	})

	middleware1 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("[m1-before]"))
			err := next(w, r)
			w.Write([]byte("[m1-after]"))
			return err
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("[m2-before]"))
			err := next(w, r)
			w.Write([]byte("[m2-after]"))
			return err
		}
	}

	middleware3 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("[m3-before]"))
			err := next(w, r)
			w.Write([]byte("[m3-after]"))
			return err
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware1, middleware2, middleware3)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	// Middlewares are applied in reverse order (from right to left)
	// So middleware3 wraps handler, middleware2 wraps middleware3, middleware1 wraps middleware2
	// Execution order: m1-before -> m2-before -> m3-before -> handler -> m3-after -> m2-after -> m1-after
	assert.Equal(t, "[m1-before][m2-before][m3-before][handler][m3-after][m2-after][m1-after]", rec.Body.String())
}

func TestApplyMiddlewares_OrderMatters(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("H"))
		return nil
	})

	middlewareA := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("A"))
			return next(w, r)
		}
	}

	middlewareB := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("B"))
			return next(w, r)
		}
	}

	// Apply in order: A, B
	wrappedHandler1 := applyMiddlewares(handler, middlewareA, middlewareB)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec1 := httptest.NewRecorder()
	wrappedHandler1(rec1, req1)

	// Apply in reverse order: B, A
	wrappedHandler2 := applyMiddlewares(handler, middlewareB, middlewareA)

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec2 := httptest.NewRecorder()
	wrappedHandler2(rec2, req2)

	// Results should be different based on order
	assert.Equal(t, "ABH", rec1.Body.String())
	assert.Equal(t, "BAH", rec2.Body.String())
}

func TestApplyMiddlewares_MiddlewareModifiesRequest(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		value := r.Header.Get("X-Custom-Header")
		w.Write([]byte(value))
		return nil
	})

	middleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			r.Header.Set("X-Custom-Header", "modified-value")
			return next(w, r)
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	assert.Equal(t, "modified-value", rec.Body.String())
}

func TestApplyMiddlewares_MiddlewareModifiesResponse(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler response"))
		return nil
	})

	middleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("X-Middleware", "applied")
			return next(w, r)
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	assert.Equal(t, "applied", rec.Header().Get("X-Middleware"))
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	assert.Equal(t, "handler response", rec.Body.String())
}

func TestApplyMiddlewares_MiddlewareReturnsError(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("this should not be called"))
		return nil
	})

	expectedErr := errors.New("middleware error")

	middleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return expectedErr
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.ErrorIs(t, err, expectedErr)
	assert.Empty(t, rec.Body.String())
}

func TestApplyMiddlewares_MiddlewareShortCircuits(t *testing.T) {
	handlerCalled := false

	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		handlerCalled = true
		return nil
	})

	middleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			// Don't call next handler
			return nil
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	assert.False(t, handlerCalled, "Handler should not be called when middleware short-circuits")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "unauthorized", rec.Body.String())
}

func TestApplyMiddlewares_HandlerReturnsError(t *testing.T) {
	expectedErr := errors.New("handler error")

	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	middlewareCalled := false

	middleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			middlewareCalled = true
			err := next(w, r)
			// Middleware can inspect the error
			if err != nil {
				w.Header().Set("X-Error-Caught", "true")
			}
			return err
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, middlewareCalled)
	assert.Equal(t, "true", rec.Header().Get("X-Error-Caught"))
}

func TestApplyMiddlewares_ErrorPropagation(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("handler error")
	})

	middleware1 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			err := next(w, r)
			if err != nil {
				// Middleware1 sees the error and adds context
				return errors.New("middleware1: " + err.Error())
			}
			return nil
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			err := next(w, r)
			if err != nil {
				// Middleware2 sees the wrapped error
				return errors.New("middleware2: " + err.Error())
			}
			return nil
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware1, middleware2)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.Error(t, err)
	// applyMiddlewares applies in reverse order, so middleware2 wraps middleware1
	// Error flows: handler -> middleware1 -> middleware2
	assert.Equal(t, "middleware1: middleware2: handler error", err.Error())
}

func TestApplyMiddlewares_MiddlewareChainWithEarlyReturn(t *testing.T) {
	middleware1Called := false
	middleware2Called := false
	middleware3Called := false
	handlerCalled := false

	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		handlerCalled = true
		return nil
	})

	middleware1 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			middleware1Called = true
			return next(w, r)
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			middleware2Called = true
			// Early return - don't call next
			return errors.New("stopped at middleware2")
		}
	}

	middleware3 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			middleware3Called = true
			return next(w, r)
		}
	}

	wrappedHandler := applyMiddlewares(handler, middleware1, middleware2, middleware3)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.Error(t, err)
	assert.True(t, middleware1Called, "First middleware should be called")
	assert.True(t, middleware2Called, "Second middleware should be called")
	assert.False(t, middleware3Called, "Third middleware should not be called")
	assert.False(t, handlerCalled, "Handler should not be called")
}

func TestApplyMiddlewares_StatePassedThroughMiddlewares(t *testing.T) {
	handler := Handler(func(w http.ResponseWriter, r *http.Request) error {
		count := r.Header.Get("X-Count")
		w.Write([]byte("count: " + count))
		return nil
	})

	incrementMiddleware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			current := r.Header.Get("X-Count")
			if current == "" {
				current = "0"
			}
			// Simulate incrementing a counter
			newCount := current + "1"
			r.Header.Set("X-Count", newCount)
			return next(w, r)
		}
	}

	wrappedHandler := applyMiddlewares(handler, incrementMiddleware, incrementMiddleware, incrementMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	err := wrappedHandler(rec, req)

	assert.NoError(t, err)
	assert.Equal(t, "count: 0111", rec.Body.String())
}
