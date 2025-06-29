package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		checkHeaders   bool
	}{
		{"GET request", "GET", "https://example.com", http.StatusOK, true},
		{"POST request", "POST", "https://api.example.com", http.StatusOK, true},
		{"OPTIONS request", "OPTIONS", "https://test.com", http.StatusNoContent, true},
		{"PUT request", "PUT", "", http.StatusOK, true},
		{"DELETE request", "DELETE", "https://localhost:3000", http.StatusOK, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)

			// Add CORS middleware
			router.Use(CORSMiddleware())

			// Add a test route
			router.Any("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			// Create request
			req := httptest.NewRequest(test.method, "/test", nil)
			if test.origin != "" {
				req.Header.Set("Origin", test.origin)
			}

			// Add headers for OPTIONS request
			if test.method == "OPTIONS" {
				req.Header.Set("Access-Control-Request-Method", "POST")
				req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")
			}

			c.Request = req
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}

			if test.checkHeaders {
				// Check CORS headers
				originHeader := w.Header().Get("Access-Control-Allow-Origin")
				if test.origin != "" && originHeader != test.origin {
					t.Errorf("Expected Access-Control-Allow-Origin %s, got %s", test.origin, originHeader)
				}

				credentialsHeader := w.Header().Get("Access-Control-Allow-Credentials")
				if credentialsHeader != "true" {
					t.Errorf("Expected Access-Control-Allow-Credentials true, got %s", credentialsHeader)
				}

				// Check OPTIONS-specific headers
				if test.method == "OPTIONS" {
					methodHeader := w.Header().Get("Access-Control-Allow-Methods")
					if methodHeader == "" {
						t.Error("Expected Access-Control-Allow-Methods header for OPTIONS request")
					}

					headersHeader := w.Header().Get("Access-Control-Allow-Headers")
					if headersHeader == "" {
						t.Error("Expected Access-Control-Allow-Headers header for OPTIONS request")
					}

					maxAgeHeader := w.Header().Get("Access-Control-Max-Age")
					if maxAgeHeader != "7200" {
						t.Errorf("Expected Access-Control-Max-Age 7200, got %s", maxAgeHeader)
					}
				}
			}
		})
	}
}

func TestCORSMiddlewareNoOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// Don't set Origin header

	c.Request = req
	router.ServeHTTP(w, req)

	// Should still work without origin
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Should still set credentials header
	credentialsHeader := w.Header().Get("Access-Control-Allow-Credentials")
	if credentialsHeader != "true" {
		t.Errorf("Expected Access-Control-Allow-Credentials true, got %s", credentialsHeader)
	}
}

func TestCORSMiddlewareNext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	middlewareCalled := false
	handlerCalled := false

	router.Use(CORSMiddleware())
	router.Use(func(c *gin.Context) {
		middlewareCalled = true
		c.Next()
	})
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	c.Request = req
	router.ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Expected middleware to be called")
	}
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}

func TestCORSMiddlewareAbortOnOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	handlerCalled := false

	router.Use(CORSMiddleware())
	router.OPTIONS("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	c.Request = req
	router.ServeHTTP(w, req)

	// Should abort and not reach the handler
	if handlerCalled {
		t.Error("Expected handler NOT to be called for OPTIONS request")
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func BenchmarkCORSMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		router.ServeHTTP(w, req)
	}
}
