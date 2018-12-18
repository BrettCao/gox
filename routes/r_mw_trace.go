package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"moddns/app/util"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(allowPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !util.CheckPrefix(c.Request.URL.Path, allowPrefixes...) {
			c.Next()
			return
		}

		traceID := c.Query("X-Request-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set(util.ContextKeyTraceID, traceID)
		c.Next()
	}
}
