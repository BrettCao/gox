package http

import (
	"fmt"
	"moddns/app/http/context"
	"moddns/app/http/ctl"
	"moddns/app/service/mysql"
	"moddns/routes"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Init 初始化所有服务
func Init(db *mysql.DB, enforcer *casbin.Enforcer, ctlCommon *ctl.Common) *gin.Engine {
	gin.SetMode(viper.GetString("run_mode"))
	app := gin.New()

	// 注册中间件
	apiPrefixes := []string{"/api/"}
	app.Use(routes.TraceMiddleware(apiPrefixes...))
	app.Use(routes.LoggerMiddleware(apiPrefixes, "/api/v1/loggers"))
	app.Use(routes.RecoveryMiddleware())
	app.Use(routes.SessionMiddleware(db, apiPrefixes...))

	app.NoMethod(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("方法不允许"), 405)
	}))

	app.NoRoute(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("资源不存在"), 404)
	}))

	// 注册/api/v1路由
	routes.APIV1Handler(app, enforcer, ctlCommon)

	// 加载casbin策略数据
	err := loadCasbinPolicyData(ctlCommon)
	if err != nil {
		panic("加载casbin策略数据发生错误：" + err.Error())
	}

	return app
}

// 加载casbin策略数据，包括角色权限数据、用户角色数据
func loadCasbinPolicyData(ctlCommon *ctl.Common) error {
	err := ctlCommon.RoleAPI.RoleBll.LoadAllPolicy()
	if err != nil {
		return err
	}

	err = ctlCommon.UserAPI.UserBll.LoadAllPolicy()
	if err != nil {
		return err
	}
	return nil
}
