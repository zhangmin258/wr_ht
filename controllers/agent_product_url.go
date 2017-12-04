package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
)

type AgentProUrlController struct {
	BaseController
}

//获取商品外放链接列表
func (c *AgentProUrlController) GetProductUrl() {
	defer c.ServeJSON()
	id, _ := c.GetInt("id")
	urls, err := models.GetProductUrlList(id)

	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询商品的外放链接异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询产品异常" }
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "urls": urls }
}
