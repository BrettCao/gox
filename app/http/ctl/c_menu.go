package ctl

import (
	"moddns/app/util"
	"strings"

	"moddns/app/bll"
	"moddns/app/schema"
	"moddns/app/http/context"
)

// Menu 菜单管理
type Menu struct {
	MenuBll *bll.Menu `inject:""`
}

// Query 查询数据
func (a *Menu) Query(ctx *context.Context) {
	switch ctx.Query("type") {
	case "page":
		a.QueryPage(ctx)
	case "tree":
		a.QueryTree(ctx)
	default:
		ctx.ResBadRequest(nil)
	}
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	params := schema.MenuQueryParam{
		Name:     ctx.Query("name"),
		ParentID: ctx.Query("parent_id"),
		Status:   util.S(ctx.Query("status")).Int(),
		Type:     util.S(ctx.Query("mtype")).Int(),
	}

	total, items, err := a.MenuBll.QueryPage(ctx.NewContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResPage(total, items)
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(ctx *context.Context) {
	params := schema.MenuSelectQueryParam{
		Name:   ctx.Query("name"),
		Status: util.S(ctx.Query("status")).Int(),
	}

	if util.S(ctx.Query("is_menu")).Int() == 1 {
		params.Types = []int{10, 20, 30}
	}

	treeData, err := a.MenuBll.QueryTree(ctx.NewContext(), params)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResList(treeData)
}

// Get 查询指定数据
func (a *Menu) Get(ctx *context.Context) {
	item, err := a.MenuBll.Get(ctx.NewContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *Menu) Create(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	item.Creator = ctx.GetUserID()
	err := a.MenuBll.Create(ctx.NewContext(), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	newItem, err := a.MenuBll.Get(ctx.NewContext(), item.RecordID)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResSuccess(newItem)
}

// Update 更新数据
func (a *Menu) Update(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.MenuBll.Update(ctx.NewContext(), ctx.Param("id"), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *Menu) Delete(ctx *context.Context) {
	err := a.MenuBll.Delete(ctx.NewContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *Menu) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")

	for _, id := range ids {
		err := a.MenuBll.Delete(ctx.NewContext(), id)
		if err != nil {
			ctx.ResInternalServerError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
func (a *Menu) Enable(ctx *context.Context) {
	err := a.MenuBll.UpdateStatus(ctx.NewContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
func (a *Menu) Disable(ctx *context.Context) {
	err := a.MenuBll.UpdateStatus(ctx.NewContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}
