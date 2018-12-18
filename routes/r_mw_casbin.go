package routes

import (
	"fmt"
	"moddns/app/http/context"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.NewContext(c)

		if b, err := enforcer.EnforceSafe(ctx.GetUserID(), c.Request.URL.Path, c.Request.Method); err != nil {
			ctx.ResError(errors.Wrap(err, "验证权限发生错误"), http.StatusInternalServerError)
			return
		} else if !b {
			ctx.ResError(fmt.Errorf("没有操作权限"), http.StatusUnauthorized, 9998)
			return
		}
		c.Next()
	}
}
