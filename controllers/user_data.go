package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type UserDataController struct {
	BaseController
}

//获取累计注册用户
//@router /queryusersregistercount [post]
func (this *UserDataController) QueryUsersRegisterCount() {
	defer this.ServeJSON()
	count, err := models.QueryUsersRegisterCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计注册用户数量失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取累计注册用户数量失败"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "data": count}
}

//获取累计认证用户
//@router /queryusersauthcount [post]
func (this *UserDataController) QueryUsersAuthCount() {
	defer this.ServeJSON()
	count, err := models.QueryUsersAuthCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计认证用户数量失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取累计认证用户数量失败"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "data": count}
}

//获取累计放款用户
//@router /queryusersloancount [post]
func (this *UserDataController) QueryUsersLoanCount() {
	defer this.ServeJSON()
	count, err := models.QueryUsersLoanCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取放款用户数量失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取累计放款用户数量失败"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "data": count}
}

//获取用户数据-数据明细
//@router /queryusersdatadetail [post]
func (this *UserDataController) QueryUsersDataDetail() {
	defer this.ServeJSON()
	pageNum, _ := this.GetInt("pageNum", 1)
	data := make(map[string]interface{})
	data["count"] = 0
	data["pageNum"] = 0
	data["data"] = nil
	data["ret"] = 403
	count, err := models.QueryUsersDataDetailCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取数据明细总条数失败", err.Error(), this.Ctx.Input)
		data["err"] = "获取数据明细总条数失败"
		return
	}
	count += 1
	start := utils.StartIndex(pageNum, utils.PageSize20)
	end := utils.PageSize20
	if pageNum == 1 {
		end = end - 1
	} else {
		start = start - 1
	}
	//获取用户数据详情
	usersDataDetail, err := models.QueryUsersDataDetail(start, end)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取用户数据详情", err.Error(), this.Ctx.Input)
		data["err"] = "获取用户数据详情异常"
		return
	}

	var dataDetail []models.UsersDataDetail
	//获取今日数据明细
	if pageNum == 1 {
		var todayDetailData models.UsersDataDetail
		var todayActiveUsersCount int
		todayDetailData, err = models.QueryTodayDailyDatas()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取今日数据明细异常", err.Error(), this.Ctx.Input)
			data["err"] = "获取今日数据明细异常"
			return
		}
		todayActiveUsersCount, err = models.QueryTodayActiveData()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取今日数据明细异常", err.Error(), this.Ctx.Input)
			data["err"] = "获取今日数据明细异常"
			return
		}
		todayDetailData.ActiveCount = todayActiveUsersCount
		dataDetail = append(dataDetail, todayDetailData)
	}
	dataDetail = append(dataDetail, usersDataDetail[:]...)

	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取需要页数失败", "获取需要页数失败", err.Error(), this.Ctx.Input)
		data["err"] = "获取需要页数失败"
		return
	}
	data["count"] = count
	data["data"] = dataDetail
	data["ret"] = 200
	data["pageNum"] = pageNum
	data["pageCount"] = pageCount
	data["pageSize"] = utils.PageSize20
	this.Data["json"] = data
	return
}
