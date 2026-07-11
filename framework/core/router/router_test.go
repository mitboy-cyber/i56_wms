package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_GET(t *testing.T) {
	r := New()
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRouter_MethodMismatch(t *testing.T) {
	r := New()
	r.POST("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(201)
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 405 {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}

	// POST should work
	req2 := httptest.NewRequest("POST", "/hello", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != 201 {
		t.Errorf("expected 201, got %d", w2.Code)
	}
}

func TestRouter_AllMethods(t *testing.T) {
	r := New()
	called := false
	handler := func(w http.ResponseWriter, req *http.Request) {
		called = true
		w.WriteHeader(200)
	}

	r.GET("/test", handler)
	r.POST("/test", handler)
	r.PUT("/test", handler)
	r.PATCH("/test", handler)
	r.DELETE("/test", handler)

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, method := range methods {
		called = false
		req := httptest.NewRequest(method, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if !called {
			t.Errorf("handler not called for %s", method)
		}
		if w.Code != 200 {
			t.Errorf("expected 200 for %s, got %d", method, w.Code)
		}
	}
}

func TestRouter_WithPrefix(t *testing.T) {
	r := New().WithPrefix("/api/v1")
	r.GET("/orders", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	// Correct path
	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Wrong path
	req2 := httptest.NewRequest("GET", "/orders", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code == 200 {
		t.Error("expected non-200 for wrong path")
	}
}

func TestRouter_Handle(t *testing.T) {
	r := New()
	// Use WithPrefix so the sub-router strips the mount path
	sub := New().WithPrefix("/api")
	sub.GET("/nested", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	r.Handle("/api/", sub)

	req := httptest.NewRequest("GET", "/api/nested", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRouter_Use_Middleware(t *testing.T) {
	r := New()

	var order []string
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			order = append(order, "mw1")
			next.ServeHTTP(w, req)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			order = append(order, "mw2")
			next.ServeHTTP(w, req)
		})
	}

	r.Use(mw1, mw2)
	r.GET("/test", func(w http.ResponseWriter, req *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Middleware applied in order: mw1 is outermost (runs first), mw2 inner
	if len(order) != 3 || order[0] != "mw1" || order[1] != "mw2" || order[2] != "handler" {
		t.Errorf("unexpected order: %v", order)
	}
}
