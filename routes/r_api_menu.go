package routes

import (
	"github.com/gin-gonic/gin"
	"moddns/app/http/context"
	"moddns/app/http/ctl"
)

// APIMenuRouter 注册/menus路由
func APIMenuRouter(g *gin.RouterGroup, menu *ctl.Menu) {
	g.GET("/menus", context.WrapContext(menu.Query, "查询菜单数据"))
	g.GET("/menus/:id", context.WrapContext(menu.Get, "查询指定菜单数据"))
	g.POST("/menus", context.WrapContext(menu.Create, "创建菜单数据"))
	g.PUT("/menus/:id", context.WrapContext(menu.Update, "更新菜单数据"))
	g.DELETE("/menus/:id", context.WrapContext(menu.Delete, "删除菜单数据"))
	g.DELETE("/menus", context.WrapContext(menu.DeleteMany, "删除多条菜单数据"))
	g.PATCH("/menus/:id/enable", context.WrapContext(menu.Enable, "启用菜单数据"))
	g.PATCH("/menus/:id/disable", context.WrapContext(menu.Disable, "禁用菜单数据"))
}
