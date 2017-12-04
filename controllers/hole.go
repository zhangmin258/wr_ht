package controllers

/*
*口子接口
 */

import (
	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type HoleController struct {
	BaseController
}

//口子列表页
func (c *HoleController) HoleList() {
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page", 1)
	phoneType, err := c.GetInt("phoneType", 1) //1安卓，2ios
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询口子列表异常！", err.Error(), c.Ctx.Input)
	}
	name := c.GetString("name")
	count := 0
	var holeList []*models.HoleInfo
	condition := ""
	params := []interface{}{}
	if name != "" {
		condition += " AND name = ? "
		params = append(params, name)
	}
	if phoneType == 1 {
		holeList, err = models.GetAndroidHoleList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询安卓口子列表异常！", err.Error(), c.Ctx.Input)
		}
		count, err = models.GetAndroidHoleCount(condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询安卓口子总数异常！", err.Error(), c.Ctx.Input)
		}
	} else {
		holeList, err = models.GetIosHoleList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询ios口子列表异常！", err.Error(), c.Ctx.Input)
		}
		count, err = models.GetIosHoleCount(condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询ios口子总数异常！", err.Error(), c.Ctx.Input)
		}
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询需要页数异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["name"] = name
	c.Data["phoneType"] = phoneType
	c.Data["holeList"] = holeList
	c.TplName = "pay-management/newCut_management.html"
}

//保存口子
func (c *HoleController) SaveHole() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var hole models.HoleInfo
	err := c.ParseForm(&hole)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析口子出错！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析口子出错"
		return
	}
	if hole.Id == 0 {
		err = models.SaveHole(hole)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存口子出错！", err.Error(), c.Ctx.Input)
			resultMap["err"] = "保存口子出错"
			return
		}
	} else {
		err = models.UpdateHole(hole)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改口子出错！", err.Error(), c.Ctx.Input)
			resultMap["err"] = "修改口子出错"
			return
		}
	}

	resultMap["ret"] = 200
	resultMap["msg"] = "保存口子成功"
}

//删除口子
func (c *HoleController) DelHole() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	id, err := c.GetInt("hid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "接收口子id异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "接收口子id异常"
		return
	}
	phoneType, err := c.GetInt("phoneType") //1安卓，2ios
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询口子列表异常！", err.Error(), c.Ctx.Input)
		return
	}
	//删除口子
	err = models.DelHoleById(id, phoneType)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除口子异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "删除口子异常"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "删除口子成功"
}

//新口子配置
func (c *HoleController) HoleConfig() {
	c.IsNeedTemplate()
	id := 2 // 新口子的配置的信息id
	holeConfig, err := models.GetHoleConfig(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取新口子配置异常！", err.Error(), c.Ctx.Input)
		c.Data["err"] = "获取新口子配置异常"
	}
	c.Data["holeConfig"] = holeConfig
	c.TplName = "pay-management/newCut_config.html"
}

//保存新口子配置
func (c *HoleController) SaveHoleConfig() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var holeConfig models.HoleConfig
	err := c.ParseForm(&holeConfig)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析口子配置出错！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析口子配置出错"
		return
	}
	//修改新口子配置
	err = models.UpdateHoleConfig(holeConfig)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改口子配置异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "修改口子配置异常"
		return
	}
	//删除缓存
	if utils.Re == nil && utils.Rc.IsExist(utils.WR_CACHE_KEY_MONEY_SERVICE+strconv.Itoa(holeConfig.Id)) {
		utils.Rc.Delete(utils.WR_CACHE_KEY_MONEY_SERVICE + strconv.Itoa(holeConfig.Id))
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "修改口子配置成功"
}
