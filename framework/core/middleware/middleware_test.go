package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any) {}
func (l testLogger) Info(msg string, args ...any)  {}
func (l testLogger) Warn(msg string, args ...any)  {}
func (l testLogger) Error(msg string, args ...any) {}
func (l testLogger) With(args ...any) Logger       { return l }
func (l testLogger) WithGroup(name string) Logger  { return l }

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

func TestChain_Order(t *testing.T) {
	order := []string{}

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	chained := Chain(handler, mw1, mw2)
	chained.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	// mw2 is outermost (applied last), so it wraps mw1
	expected := []string{"mw2-before", "mw1-before", "handler", "mw1-after", "mw2-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d steps, got %d: %v", len(expected), len(order), order)
	}
	for i, e := range expected {
		if order[i] != e {
			t.Errorf("step %d: expected %q, got %q", i, e, order[i])
		}
	}
}

func TestRecovery_CatchesPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	rec := Recovery(testLogger{})(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	rec.ServeHTTP(w, req)
	// Should return 500
	if w.Code < 500 {
		t.Errorf("expected 5xx, got %d", w.Code)
	}
}

func TestRequestID_SetsHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mw := RequestID()(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header")
	}
}

func TestRequestID_RespectsIncoming(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mw := RequestID()(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "my-custom-id")
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)
	if w.Header().Get("X-Request-ID") != "my-custom-id" {
		t.Errorf("expected 'my-custom-id', got %q", w.Header().Get("X-Request-ID"))
	}
}

func TestLogger_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mw := Logger(testLogger{})(handler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)
	// Should not panic or error
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCORS_Defaults(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mw := CORS(nil)(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("expected CORS origin, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_Options(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for OPTIONS")
	})
	mw := CORS(nil)(handler)

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

func TestCORS_CustomOptions(t *testing.T) {
	opts := &CORSOptions{
		AllowedOrigins:   []string{"https://app.example.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mw := CORS(opts)(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("expected Allow-Credentials header")
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	opts := &CORSOptions{
		AllowedOrigins: []string{"https://app.example.com"},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mw := CORS(opts)(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("expected no CORS header for disallowed origin")
	}
}

func TestJoin(t *testing.T) {
	result := join([]string{"a", "b", "c"})
	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got %q", result)
	}
}
