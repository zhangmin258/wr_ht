package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type MakeConnController struct {
	BaseController
}

//查询并展示预约合作顾问的信息
//router /getmsgofmakeconn [get]
func (this *MakeConnController) GetMsgOfMakeConn() {
	defer this.IsNeedTemplate()
	pageNum, _ := this.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	makeConn, err := models.GetAllMakeConn(utils.StartIndex(pageNum, utils.PAGE_SIZE), utils.PAGE_SIZE)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询预约合作信息异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "makeconn/getmsgofmakeconn|查询预约信息异常"}
		return
	}
	pageCount, err := utils.GetPageCount(len(makeConn), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询分页信息异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "makeconn/getmsgofmakeconn|查询预约信息异常"}

	}
	this.Data["Count"] = len(makeConn)
	this.Data["PageCount"] = pageCount
	this.Data["MakeConn"] = makeConn
	this.TplName = ""
}
