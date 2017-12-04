package controllers

/*
*征信口子接口
 */

import (
	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type PlateformHoleController struct {
	BaseController
}

//征信口子列表
func (c *PlateformHoleController) PlaCreditList() {
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page", 1)
	name := c.GetString("name")
	condition := ""
	params := []interface{}{}
	if name != "" {
		condition += " AND plateform_name = ? "
		params = append(params, name)
	}
	plaCreditList, err := models.GetPlaCreditList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询口子列表异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.PlaCreditCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询口子总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询需要页数异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["plaCreditList"] = plaCreditList
	c.TplName = "pay-management/zxCut_management.html"
}

//保存征信口子
func (c *PlateformHoleController) SavePlaCredit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var plaCredit models.PlaCredit
	err := c.ParseForm(&plaCredit)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析征信异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析征信口子出错"
		return
	}
	if plaCredit.Id == 0 {
		//查出所有征信口子名称
		plaCredits, err := models.GetAllPlaCredit()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询征信口子异常！", err.Error(), c.Ctx.Input)
			resultMap["err"] = "查询征信口子出错"
			return
		}
		for _, v := range plaCredits {
			if v.PlateformName == plaCredit.PlateformName {
				resultMap["err"] = "已存在该名称的征信口子！"
				return
			}
		}
		err = models.SavePlaCredit(&plaCredit)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存征信口子异常！", err.Error(), c.Ctx.Input)
			resultMap["err"] = "保存征信口子出错"
			return
		}
	} else {
		err = models.UpdatePlaCredit(&plaCredit)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存征信口子异常！", err.Error(), c.Ctx.Input)
			resultMap["err"] = "保存征信口子出错"
			return
		}
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "保存征信口子成功"
}

//删除征信口子
func (c *PlateformHoleController) DelPlaCredit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	pcid, err := c.GetInt("pcid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取征信id异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取征信id异常"
		return
	}
	err = models.DelPlaCreditById(pcid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除征信口子异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "删除征信口子异常"
		return
	}
	resultMap["ret"] = 200
	resultMap["err"] = "删除征信口子成功"
}

//征信口子配置
func (c *PlateformHoleController) PlaCreditConfig() {
	c.IsNeedTemplate()
	id := 5 // 征信口子的配置的信息id
	holeConfig, err := models.GetHoleConfig(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取新口子配置异常！", err.Error(), c.Ctx.Input)
		c.Data["err"] = "获取新口子配置异常"
	}
	c.Data["holeConfig"] = holeConfig
	c.TplName = "pay-management/zxCut_config.html"
}

func (c *PlateformHoleController) SavePlaCreditConfig() {
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
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改征信口子异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "修改征信口子异常"
		return
	}
	//删除缓存
	if utils.Re == nil && utils.Rc.IsExist(utils.WR_CACHE_KEY_MONEY_SERVICE+strconv.Itoa(holeConfig.Id)) {
		utils.Rc.Delete(utils.WR_CACHE_KEY_MONEY_SERVICE + strconv.Itoa(holeConfig.Id))
	}

	resultMap["ret"] = 200
	resultMap["msg"] = "修改口子配置成功"
}
