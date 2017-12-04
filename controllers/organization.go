package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
)

/*
机构相关接口
*/
type OrganizationController struct {
	BaseController
}

//新增机构和商务联系人
//@router /add [post]
func (c *OrganizationController) AddOrganization() {
	defer c.ServeJSON()
	var org models.Organization
	var bus models.BusinessLinkman
	if err := c.ParseForm(&org); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}
	if err := c.ParseForm(&bus); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}

	orgId, err := models.AddOrganization(&org)

	if err != nil || orgId == 0 {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "新增机构信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}
	busId, err := models.AddBusinessLinkman(&bus)
	if err != nil || busId == 0 {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "新增商务联系人信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "orgId": orgId, "busId": busId}
	return
}

//修改机构和商务联系人
//@router /update [post]
func (c *OrganizationController) UpdateOrganizationAndBusinessLinkman() {
	defer c.ServeJSON()
	var org models.Organization
	var bus models.BusinessLinkman

	if err := c.ParseForm(&org); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "读取机构信息参数失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "读取机构信息参数失败！", "err": err.Error()}
		return
	}
	if err := c.ParseForm(&bus); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "读取商务联系人参数失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "读取商务联系人参数失败！", "err": err.Error()}
		return
	}
	//更新机构
	if err := models.UpdateOrgById(&org); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新机构异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "更新机构失败！", "err": err.Error()}
		return
	}
	//更新商务联系人
	if err := models.UpdateBusById(&bus); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新商务联系人异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "更新商务联系人失败！", "err": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
	return
}
