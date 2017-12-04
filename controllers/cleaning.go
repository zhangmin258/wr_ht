package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
)

/*
清算信息相关接口
*/
type CleaningController struct {
	BaseController
}

//添加清算信息
//@router /addCleaning [post]
func (c *CleaningController) AddCleaning() {
	defer c.ServeJSON()
	c.IsNeedTemplate()
	//接受参数
	var cleaning models.ProductCleaning
	err := c.ParseForm(&cleaning)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "cleaning/AddCleaning|参数解析异常!"}
		return
	}
	err = models.AddCleaning(&cleaning)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入结算数据异常！", err.Error(), c.Ctx.Input)
		//插入清算信息异常`
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "cleaning/AddCleaning|插入结算数据异常!"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}

}

//更新结算信息  H5
//@router /updateH5 [post]
func (c *CleaningController) UpdateCleaningH5() {
	defer c.ServeJSON()
	var cleaning models.ProductCleaning
	if err := c.ParseForm(&cleaning); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "H5产品结算信息参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "cleaning/update|H5产品参数解析异常！"}
		return
	}
	if err := models.UpdateCleaningH5(&cleaning); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "H5产品更新清算信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "cleaning/update|H5产品更新清算信息异常！"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//更新结算信息  API
//@router /updateAPI [post]
func (c *CleaningController) UpdateCleaningAPI() {
	defer c.ServeJSON()
	var cleaning models.ProductCleaning
	if err := c.ParseForm(&cleaning); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "API产品结算参数解析异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "cleaning/updateSmall|API产品结算参数解析异常！"}
		return
	}
	if err := models.UpdateCleaningAPI(&cleaning); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "API产品更新清算信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "cleaning/updateSmall|API产品更新清算信息异常！"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}
