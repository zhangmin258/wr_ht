package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type UsersBehaviorController struct {
	BaseController
}

//跳转至意愿用户页面
func (this *UsersBehaviorController) GetUsersBehavior() {
	this.IsNeedTemplate()
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	condition := ""
	params := []string{}
	phoneNum := this.GetString("phone_number")
	if phoneNum != "" {
		params = append(params, phoneNum)
		condition += " AND um.account=?"
	}
	var m []models.UsersBehavior
	cacheStr := utils.CACHE_KEY_USERS_BEHAVIOR_LIST
	data, err := utils.Rc.RedisBytes(cacheStr)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取历史意愿用户列表错误", err.Error(), this.Ctx.Input)
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "解析历史意愿用户列表错误", err.Error(), this.Ctx.Input)
	}
	list, err := models.GetBehaviorUsersList(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取意愿用户列表错误", err.Error(), this.Ctx.Input)
	}
	for _, v := range m {
		if phoneNum == "" || v.Account == phoneNum {
			list = append(list, v)
		}
	}
	count := len(list)
	pageCount := utils.PageCount(count, utils.PageSize10)
	if count != 0 {
		if page == pageCount {
			list = list[(page-1)*utils.PageSize10:]
		} else {
			list = list[(page-1)*utils.PageSize10 : page*utils.PageSize10]
		}
	}
	this.Data["behaviorUsersList"] = list
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = page
	this.TplName = "loan-under-steady/aspiration_user.html"
}

//意愿用户导出Excel
func (this *UsersBehaviorController) BehaviorUsersExportExcel() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	condition := ""
	params := []string{}
	startTime := this.GetString("startTime")
	endTime := this.GetString("endTime")
	if startTime != "" {
		condition += " AND DATE_FORMAT(ub.create_time,'%Y-%m-%d') >=?"
		params = append(params, startTime)
	}
	if endTime != "" {
		condition += " AND DATE_FORMAT(ub.create_time,'%Y-%m-%d')<=?"
		params = append(params, endTime)
	}
	list, err := models.ExportBehaviorUsersList(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "导出Excel时查询数据异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "导出Excel时查询数据异常！"
		return
	}
	behaviorData := [][]string{{"创建时间", "手机号", "姓名", "意愿类型"}}
	colWidth := []float64{24.0, 20.0, 20.0, 20.0}
	action := ""
	for _, v := range list {
		if v.Action == 1 {
			action = "页面停留时长"
		} else if v.Action == 2 {
			action = "咨询客服"
		} else if v.Action == 3 {
			action = "进入支付流程"
		}
		listData := []string{v.CreateTime, v.Account, v.Name, action}
		behaviorData = append(behaviorData, listData)
	}
	fileName, err := utils.SupportLoanExportToExcel(behaviorData, colWidth, "意愿用户列表"+time.Now().Format("2006-01-02"))
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "导出Excel异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "导出Excel异常！"
		return
	}
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fileName)
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, fileName)
	err = os.Remove(fileName)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}
