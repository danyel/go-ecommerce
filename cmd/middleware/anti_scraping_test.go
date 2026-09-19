package middleware

import (
	Http "net/http"
	HttpTest "net/http/httptest"
	Testing "testing"
	Time "time"
)

func TestRateLimiterUsesIPAndTokenAndReturnsRetryAfter(t *Testing.T) {
	limiter := NewRateLimiter(1, Time.Minute)
	handler := limiter.Middleware(Http.HandlerFunc(func(response Http.ResponseWriter, _ *Http.Request) {
		response.WriteHeader(Http.StatusNoContent)
	}))

	first := HttpTest.NewRecorder()
	request := HttpTest.NewRequest(Http.MethodGet, "/", nil)
	request.RemoteAddr = "192.0.2.1:1234"
	request.Header.Set(Authorization, "Bearer one")
	handler.ServeHTTP(first, request)
	if first.Code != Http.StatusNoContent {
		t.Fatalf("first request status = %d", first.Code)
	}

	second := HttpTest.NewRecorder()
	handler.ServeHTTP(second, request)
	if second.Code != Http.StatusTooManyRequests {
		t.Fatalf("second request status = %d", second.Code)
	}
	if second.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header")
	}

	differentToken := HttpTest.NewRequest(Http.MethodGet, "/", nil)
	differentToken.RemoteAddr = request.RemoteAddr
	differentToken.Header.Set(Authorization, "Bearer two")
	third := HttpTest.NewRecorder()
	handler.ServeHTTP(third, differentToken)
	if third.Code != Http.StatusNoContent {
		t.Fatalf("different token status = %d", third.Code)
	}
}

func TestSecurityHeadersRequireAllowedOrigin(t *Testing.T) {
	handler := SecurityHeaders([]string{"https://shop.example"})(Http.HandlerFunc(func(response Http.ResponseWriter, _ *Http.Request) {
		response.WriteHeader(Http.StatusNoContent)
	}))

	request := HttpTest.NewRequest(Http.MethodOptions, "/", nil)
	request.Header.Set("Origin", "https://evil.example")
	response := HttpTest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != Http.StatusForbidden {
		t.Fatalf("disallowed origin status = %d", response.Code)
	}

	request = HttpTest.NewRequest(Http.MethodOptions, "/", nil)
	request.Header.Set("Origin", "https://shop.example")
	response = HttpTest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != Http.StatusNoContent {
		t.Fatalf("allowed origin status = %d", response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "https://shop.example" {
		t.Fatal("expected explicit CORS origin")
	}
}
