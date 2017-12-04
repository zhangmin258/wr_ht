package controllers

import (
	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

/*
代理相关接口
*/
type AgentController struct {
	BaseController
}

//获取下级代理商列表
func (c *AgentController) GetAgentList() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//代理商名称
	if orgName := c.GetString("orgName"); orgName != "" {
		condition += " and org_name = ?"
		params = append(params, orgName)
	}
	//查询
	agentList, err := models.GetAgentList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询代理商列表异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetAgentCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询代理商总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
	}
	c.Data["agentList"] = agentList
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "channel-management/lower_level_management.html"
}

//保存下级代理商
func (c *AgentController) AddAgent() {
	defer c.ServeJSON()
	//接受参数
	var agent models.Agent
	err := c.ParseForm(&agent)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析数据异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "解析数据异常"}
		return
	}
	if agent.Id == 0 { //插入下级代理数据
		_, err := models.AddLowerLevel(&agent)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "插入下级代理数据异常"}
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入下级代理数据异常！", err.Error(), c.Ctx.Input)
			return
		}
	} else { //修改下级代理数据
		err := models.EditLowerLevel(&agent)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改下级代理数据异常"}
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改下级代理数据异常！", err.Error(), c.Ctx.Input)
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
}

//跳转到编辑或者添加下级代理商页面
func (c *AgentController) EditAgent() {
	c.IsNeedTemplate()
	idStr := c.GetString("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			agent, err := models.GetAgentById(id)
			if err != nil && err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询下级代理数据异常！", err.Error(), c.Ctx.Input)
			}
			c.Data["agent"] = agent
		}
	}
	c.TplName = "channel-management/add_agent.html"
}

//获取商品列表
func (c *AgentController) GetProductList() {
	defer c.ServeJSON()
	name := c.GetString("name")
	condition := " AND name LIKE ? "
	param := "%" + name + "%"
	products, err := models.GetProductIdAndName(condition, param)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询产品异常"}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品异常！", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "products": products}
	return
}

//获取上線商品列表
func (c *AgentController) GetProductListUse() {
	defer c.ServeJSON()
	name := c.GetString("name")
	condition := " AND name LIKE ? AND is_use = 0 "
	param := "%" + name + "%"
	products, err := models.GetProductIdAndName(condition, param)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询产品异常"}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品异常！", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "products": products}
	return
}
