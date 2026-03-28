package middleware

import (
	"encoding/json"
	"strings"
	"telemetry_bridge/internal/router/model"

	log "github.com/dredfort42/go_logger"
	"github.com/gin-gonic/gin"
)

// Audit logs one structured JSON line per request with user/session/request correlation.
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		userID := c.GetString(model.CtxUserID)
		if userID == "" {
			userID = "guest"
		}

		sessionID := c.GetString(model.CtxSessionID)
		requestID := c.GetString(model.CtxRequestID)

		query := c.Request.URL.RawQuery
		if strings.Contains(query, "password=") {
			query = "password=REDACTED"
		}

		evt := map[string]any{
			"request_id": requestID,
			"session_id": sessionID,
			"user_id":    userID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		b, _ := json.Marshal(evt)
		log.Info.Println(string(b))
	}
}
