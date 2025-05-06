package api

import (
"log"
"time"

"github.com/gin-gonic/gin"
)

// Logger is a middleware function that logs the request
func Logger() gin.HandlerFunc {
return func(c *gin.Context) {
// Start timer
start := time.Now()

// Process request
c.Next()

// Calculate latency
latency := time.Since(start)

// Log request
log.Printf(
"[%s] %s %s %s %d %s",
c.Request.Method,
c.Request.URL.Path,
c.Request.Proto,
latency,
c.Writer.Status(),
c.ClientIP(),
)
}
}

// CORS is a middleware function that handles CORS
func CORS() gin.HandlerFunc {
return func(c *gin.Context) {
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

if c.Request.Method == "OPTIONS" {
c.AbortWithStatus(204)
return
}

c.Next()
}
}
