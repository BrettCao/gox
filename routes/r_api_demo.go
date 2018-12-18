package routes

import (
	"github.com/gin-gonic/gin"
	"moddns/app/http/context"
	"moddns/app/http/ctl"
)

// APIDemoRouter 注册/demos路由
func APIDemoRouter(g *gin.RouterGroup, demo *ctl.Demo) {
	g.GET("/demos", context.WrapContext(demo.Query, "查询示例数据"))
	g.GET("/demos/:id", context.WrapContext(demo.Get, "查询指定示例数据"))
	g.POST("/demos", context.WrapContext(demo.Create, "创建示例数据"))
	g.PUT("/demos/:id", context.WrapContext(demo.Update, "更新示例数据"))
	g.DELETE("/demos/:id", context.WrapContext(demo.Delete, "删除示例数据"))
	g.DELETE("/demos", context.WrapContext(demo.DeleteMany, "删除多条示例数据"))
}
