package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type CreditCardController struct {
	BaseController
}

//显示信用卡列表
func (this *CreditCardController) GetCreditCardList() {
	this.IsNeedTemplate()
	condition := ""
	params := []string{}
	pageNum, _ := this.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	count, _ := models.GetCreditCardCount(condition, params)
	creditCardList, err := models.ShowCreditCardList(condition, params, utils.StartIndex(pageNum, utils.PageSize5), utils.PageSize5)
	pageCount, _ := utils.GetPageCount(count, utils.PageSize5)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询信用卡异常！", err.Error(), this.Ctx.Input)
	}
	this.Data["creditCardList"] = creditCardList
	this.Data["pageNum"] = pageNum
	this.Data["pageCount"] = pageCount
	this.Data["count"] = count
	this.TplName = "agent-products/creditCardList.html"
}

//跳转到新增信用卡页面
func (this *CreditCardController) JumpToAddCreditCard() {
	this.IsNeedTemplate()
	this.TplName = "agent-products/add_creditCard.html"
}

//保存信用卡信息
func (this *CreditCardController) SaveCreditCard() {
	defer this.ServeJSON()
	var card models.CreditCard
	err := this.ParseForm(&card)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "参数解析异常", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取异常"}
		return
	}
	err = models.AddCreditCard(&card)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "新增信用卡信息异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "careditcard/savecreditcard 保存异常!"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "保存成功！"}
	return
}

//查询信用卡并跳转到修改页面
func (this *CreditCardController) JumpToUpdateCreditCard() {
	this.IsNeedTemplate()
	id, err := this.GetInt("Id")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取id失败", err.Error(), this.Ctx.Input)
	}
	creditCard, err := models.GetCreditCardById(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据ID查询信用卡信息异常！", err.Error(), this.Ctx.Input)
	}
	this.Data["card"] = creditCard
	this.TplName = "agent-products/edit_creditCard.html"
}

//修改信用卡
func (this *CreditCardController) UpdateCreditCard() {
	defer this.ServeJSON()
	var card models.CreditCard
	err := this.ParseForm(&card)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "数据解析失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}
	err = models.UpdateCreditCard(&card)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "修改信用卡信息异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "保存成功！"}
	return
}

//根据ID删除信用卡信息
func (this *CreditCardController) DelCreditCard() {
	defer this.ServeJSON()
	pid, err := this.GetInt("cardId")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "参数解析异常", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "creditcard/delcreditcard|参数解析异常"}
	}
	err = models.DeleteCreditCard(pid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "信用卡删除异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "creditcard/delcreditcard|信用卡删除异常"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200}
}
