package middlewares

import (
	"log"
	"net/http"
	"os"
	"products-api-with-jwt/global"
	"products-api-with-jwt/models"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// Rate limit settings
	rateLimiter = rate.NewLimiter(1, 3) // 1 request per second with a burst of 3
	blockedIPs  = make(map[string]bool)
	mu          sync.Mutex
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL)
		c.Next()
	}
}

// RateLimiterMiddleware limits requests from the same IP address
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if blockedIPs[ip] {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusTooManyRequests,
				Message: "Too many requests",
				Data:    nil,
			})
			c.Abort()
			return
		}
		mu.Unlock()

		if !rateLimiter.Allow() {
			mu.Lock()
			blockedIPs[ip] = true
			mu.Unlock()

			go func() {
				// Get the duration from the environment variable
				durationStr := os.Getenv(global.ENVRateLimitDur)
				duration, err := strconv.Atoi(durationStr)
				if err != nil {
					duration = 60
				}

				// Determine sleep duration based on the value of times
				times := os.Getenv(global.ENVRateLimitTime)
				var sleepDuration time.Duration
				switch times {
				case "second":
					sleepDuration = time.Duration(duration) * time.Second
				case "minute":
					sleepDuration = time.Duration(duration) * time.Minute
				default:
					// Default to seconds if the value is unrecognized
					sleepDuration = time.Duration(duration) * time.Second
				}

				// Sleep for the determined duration
				time.Sleep(sleepDuration)

				// Unblock the IP
				mu.Lock()
				delete(blockedIPs, ip)
				mu.Unlock()
			}()

			c.JSON(http.StatusTooManyRequests, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusTooManyRequests,
				Message: "Too many requests",
				Data:    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
