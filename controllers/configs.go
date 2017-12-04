package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

/*
*配置接口
 */
type ConfigsController struct {
	BaseController
}

//config信息列表
func (c *ConfigsController) GetConfigsList() {
	c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	condition := ""
	params := []string{}

	count, _ := models.GetConfigsCount(condition, params)
	list, _ := models.GetConfigsList(condition, params, utils.StartIndex(page, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	pageCount, _ := utils.GetPageCount(count, utils.PAGE_SIZE20)
	c.Data["configList"] = list
	c.Data["pageNum"] = page
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "show-management/config.html"
}

//添加config页面
func (c *ConfigsController) AddConfig() {
	c.IsNeedTemplate()
	c.TplName = "show-management/addConfig.html"
}

//新增或者编辑config的页面
func (c *ConfigsController) EditConfig() {
	c.IsNeedTemplate()
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取id失败", err.Error(), c.Ctx.Input)
	}
	if id != 0 {
		config, err := models.GetConfigById(id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id获取config信息异常！", err.Error(), c.Ctx.Input)
		}
		c.Data["config"] = config
	}
	c.TplName = "show-management/addConfig.html"
}

//保存config信息
func (c *ConfigsController) SaveConfig() {
	defer c.ServeJSON()
	//接受参数
	var config models.Config
	err := c.ParseForm(&config)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "处理异常"}
		return
	}
	if config.Id == 0 { //插入config数据
		err := models.AddConfig(&config)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入config数据异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "插入config数据异常"}
			return
		}
	} else { //修改config数据
		configOld, err := models.GetConfigById(config.Id)
		err = models.UpdateConfig(&config, configOld, c.User)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改config数据异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改config数据异常"}
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
	return
}

//删除config信息
func (c *ConfigsController) DelConfig() {
	defer c.ServeJSON()
	//接受参数
	id, err := c.GetInt("id")
	if err != nil {
		//解析参数异常
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "处理异常"}
		return
	}
	configOld, err := models.GetConfigById(id)
	err = models.DelConfig(configOld, c.User)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除config数据异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改config数据异常"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
	return
}
