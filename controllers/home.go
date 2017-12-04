package controllers

import (

	//"wr_v1/cache"
	"wr_v1/cache"
	"wr_v1/services"
	"wr_v1/utils"
)

// HomeController 主页
type HomeController struct {
	BaseController
}

// Get 主页Get
func (c *HomeController) Get() {
	c.IsNeedTemplate()
	c.TplName = "index.html"
}

// Post 主页获取数据
func (c *HomeController) Post() {
	//m, err := services.GetSysMenuTreeByRoleId(c.User.RoleId) //c.User.RoleId
	m, err := services.GetSysMenuTreeByRoleId(c.User.StationId) //根据岗位ID获取菜单
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户菜单失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}

	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "data": m}
	}
	c.ServeJSON()
}
