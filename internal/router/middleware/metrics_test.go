package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"telemetry_bridge/internal/metrics"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetMetricsForTesting resets the global metrics state for testing.
func resetMetricsForTesting() {
	// Reset the metrics
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	// Initialize a fresh metrics instance
	metrics.GlobalMetrics = &metrics.Metrics{}

	// Initialize the metrics
	namespace := "test_service"
	labels := prometheus.Labels{"version": "test"}

	metrics.Add(namespace, labels)
}

func TestMetrics(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("Metrics records metrics for GET request", func(t *testing.T) {
		resetMetricsForTesting()

		// Create a new Gin router
		router := gin.New()
		router.Use(Metrics())

		// Add a test route
		router.GET("/about", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a test request
		req, err := http.NewRequest("GET", "/about", nil)
		require.NoError(t, err)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify HTTP metrics were recorded
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "/about", "200"))
		assert.Equal(t, float64(1), requestsTotal, "HTTP requests total should be incremented")

		// Verify in-flight requests returned to zero
		inFlightRequests := testutil.ToFloat64(metrics.GlobalMetrics.GetInFlight())
		assert.Equal(t, float64(0), inFlightRequests, "In-flight requests should return to zero")
	})

	t.Run("Metrics records metrics for POST request", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		router.POST("/streamshot/get", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "created"})
		})

		req, err := http.NewRequest("POST", "/streamshot/get", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify metrics were recorded with normalized endpoint
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("POST", "/streamshot/get", "201"))
		assert.Equal(t, float64(1), requestsTotal)
	})

	t.Run("Metrics skips metrics endpoint", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		router.GET("/metrics", func(c *gin.Context) {
			c.String(http.StatusOK, "prometheus metrics")
		})

		req, err := http.NewRequest("GET", "/metrics", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify NO metrics were recorded for /metrics endpoint
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "/metrics", "200"))
		assert.Equal(t, float64(0), requestsTotal, "Metrics endpoint should not record metrics for itself")
	})

	t.Run("Metrics handles nil GlobalMetrics gracefully", func(t *testing.T) {
		// Set GlobalMetrics to nil
		metrics.GlobalMetrics = nil

		router := gin.New()
		router.Use(Metrics())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		// Should not panic with nil GlobalMetrics
		assert.NotPanics(t, func() {
			router.ServeHTTP(w, req)
		})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Metrics records different status codes", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		router.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})
		router.GET("/error", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		})
		router.GET("/notfound", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})

		testCases := []struct {
			path           string
			expectedStatus int
			expectedCode   string
		}{
			{"/success", http.StatusOK, "200"},
			{"/error", http.StatusInternalServerError, "500"},
			{"/notfound", http.StatusNotFound, "404"},
		}

		for _, tc := range testCases {
			t.Run(tc.expectedCode, func(t *testing.T) {
				req, err := http.NewRequest("GET", tc.path, nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				// Verify metrics were recorded with correct status code
				requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "other", tc.expectedCode))
				assert.Equal(t, float64(1), requestsTotal)
			})
		}
	})

	t.Run("Metrics tracks in-flight requests correctly", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		// Add a slow route to test in-flight tracking
		router.GET("/slow", func(c *gin.Context) {
			time.Sleep(100 * time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"message": "slow response"})
		})

		// Start request in a goroutine
		var (
			wg          sync.WaitGroup
			maxInFlight float64
		)

		wg.Add(1)

		go func() {
			defer wg.Done()

			req, _ := http.NewRequest("GET", "/slow", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}()

		// Check in-flight requests during processing
		time.Sleep(50 * time.Millisecond)

		maxInFlight = testutil.ToFloat64(metrics.GlobalMetrics.GetInFlight())

		wg.Wait()

		// Verify in-flight requests were tracked
		assert.Equal(t, float64(1), maxInFlight, "Should have tracked in-flight request")

		// Verify in-flight requests returned to zero
		finalInFlight := testutil.ToFloat64(metrics.GlobalMetrics.GetInFlight())
		assert.Equal(t, float64(0), finalInFlight, "In-flight requests should return to zero")
	})

	t.Run("Metrics records response size", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		router.GET("/large-response", func(c *gin.Context) {
			// Create a response with known size
			response := map[string]string{
				"message": "This is a test response with some content",
				"data":    "Additional data to increase response size",
			}
			c.JSON(http.StatusOK, response)
		})

		req, err := http.NewRequest("GET", "/large-response", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Greater(t, w.Body.Len(), 0, "Response should have content")

		// Verify metrics were recorded
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "other", "200"))
		assert.Equal(t, float64(1), requestsTotal)
	})

	t.Run("Metrics handles concurrent requests", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		router.GET("/concurrent", func(c *gin.Context) {
			time.Sleep(10 * time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"message": "concurrent"})
		})

		const numRequests = 10

		var wg sync.WaitGroup

		// Send multiple concurrent requests
		for range numRequests {
			wg.Add(1)

			go func() {
				defer wg.Done()

				req, _ := http.NewRequest("GET", "/concurrent", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Code)
			}()
		}

		wg.Wait()

		// Verify all requests were recorded
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "other", "200"))
		assert.Equal(t, float64(numRequests), requestsTotal)

		// Verify in-flight requests returned to zero
		finalInFlight := testutil.ToFloat64(metrics.GlobalMetrics.GetInFlight())
		assert.Equal(t, float64(0), finalInFlight)
	})
}

func TestMetricsIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Metrics integration with multiple endpoints", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		// Add multiple routes
		router.GET("/streamshot/get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"type": "streamshot"})
		})
		router.GET("/avg_streamshot/get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"type": "avg_streamshot"})
		})
		router.GET("/about", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"type": "about"})
		})
		router.GET("/metrics", func(c *gin.Context) {
			c.String(http.StatusOK, "metrics data")
		})
		router.GET("/api/unknown", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"type": "unknown"})
		})

		// Test each endpoint
		endpoints := []struct {
			path                string
			normalizedPath      string
			shouldRecordMetrics bool
		}{
			{"/streamshot/get", "/streamshot/get", true},
			{"/avg_streamshot/get", "/streamshot/get", true},
			{"/about", "/about", true},
			{"/metrics", "/metrics", false}, // Should not record metrics
			{"/api/unknown", "other", true},
		}

		for _, endpoint := range endpoints {
			t.Run(endpoint.path, func(t *testing.T) {
				req, err := http.NewRequest("GET", endpoint.path, nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				if endpoint.shouldRecordMetrics {
					// Verify metrics were recorded
					requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", endpoint.normalizedPath, "200"))
					assert.Greater(t, requestsTotal, float64(0), "Metrics should be recorded for %s", endpoint.path)
				} else {
					// Verify metrics were NOT recorded
					requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", endpoint.normalizedPath, "200"))
					assert.Equal(t, float64(0), requestsTotal, "Metrics should NOT be recorded for %s", endpoint.path)
				}
			})
		}

		// Verify that both streamshot endpoints are normalized to the same metric
		streamshotTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "/streamshot/get", "200"))
		assert.Equal(t, float64(2), streamshotTotal, "Both streamshot endpoints should be counted together")
	})
}

func TestMetricsEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Metrics with middleware chain", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()

		// Add multiple middlewares
		router.Use(func(c *gin.Context) {
			c.Header("X-Custom-Header", "test")
			c.Next()
		})
		router.Use(Metrics())
		router.Use(func(c *gin.Context) {
			c.Header("X-Another-Header", "test2")
			c.Next()
		})

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test", w.Header().Get("X-Custom-Header"))
		assert.Equal(t, "test2", w.Header().Get("X-Another-Header"))

		// Verify metrics were recorded
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "other", "200"))
		assert.Equal(t, float64(1), requestsTotal)
	})

	t.Run("Metrics with panic recovery", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())      // HTTP metrics middleware first
		router.Use(gin.Recovery()) // Recovery middleware after

		router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		req, err := http.NewRequest("GET", "/panic", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// With recovery middleware, should return 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verify metrics were still recorded
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "other", "500"))
		assert.Equal(t, float64(1), requestsTotal)

		// Verify in-flight requests returned to zero
		finalInFlight := testutil.ToFloat64(metrics.GlobalMetrics.GetInFlight())
		assert.Equal(t, float64(0), finalInFlight)
	})

	t.Run("Metrics with very long request path", func(t *testing.T) {
		resetMetricsForTesting()

		router := gin.New()
		router.Use(Metrics())

		// Very long path
		longPath := "/api/very/long/path/with/many/segments/that/could/cause/issues"
		router.GET(longPath, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, err := http.NewRequest("GET", longPath, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Should be normalized to "other"
		requestsTotal := testutil.ToFloat64(metrics.GlobalMetrics.GetRequestsTotal().WithLabelValues("GET", "other", "200"))
		assert.Equal(t, float64(1), requestsTotal)
	})
}

// Benchmark tests.
func BenchmarkMetrics(b *testing.B) {
	gin.SetMode(gin.TestMode)
	resetMetricsForTesting()

	router := gin.New()
	router.Use(Metrics())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkMetricsWithoutMetrics(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// Set GlobalMetrics to nil for this benchmark
	metrics.GlobalMetrics = nil

	router := gin.New()
	router.Use(Metrics())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
