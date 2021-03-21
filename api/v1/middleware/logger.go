package middleware

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LoggerToFile Logs to file
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process Request
		c.Next()

		//End time
		endTime := time.Now()

		//Execution time
		duration := endTime.Sub(startTime)

		entry := log.WithFields(log.Fields{
			"method":     c.Request.Method,
			"duration":   duration,
			"requestUri": c.Request.RequestURI,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"clientIP":   c.ClientIP(),
			"referrer":   c.Request.Referer(),

			// "client_ip":  util.GetClientIP(c),
			// "user_id":    util.GetUserID(c),
			"request_id": c.Writer.Header().Get("Request-Id"),
			// "api_version": util.ApiVersion,

		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("")
		}
	}
}

// writeFunc convert func to io.Writer.
type writeFunc func([]byte) (int, error)

func (fn writeFunc) Write(data []byte) (int, error) {
	return fn(data)
}

// NewDefaultWriter Returns a new default write
func NewDefaultWriter() io.Writer {
	return writeFunc(func(data []byte) (int, error) {
		log.Debugf("%s", data)
		return 0, nil
	})
}
