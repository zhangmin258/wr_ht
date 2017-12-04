package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

//数据运营接口
type DataController struct {
	BaseController
}

//获取服务收益页面数据
//@router /getServicesIncomeData [get]
func (c *DataController) GetServicesIncomeData() {
	//设置整体加载
	c.IsNeedTemplate()
	condition := ""
	params := make([]interface{}, 0, 2)

	//新口子累计收入
	newCount, err := models.GetCounts(condition, params, 2)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取新口子累计收入失败", err.Error(), c.Ctx.Input)
	}

	//网贷征信累计收入
	loanCount, err := models.GetCreditCounts(condition, params,3)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取网贷征信累计收入失败", err.Error(), c.Ctx.Input)
	}

	//平台征信累计收入
	platformCount, err := models.GetCreditCounts(condition, params,5)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取平台征信累计收入失败", err.Error(), c.Ctx.Input)
	}

	//一个月会员累计收入
	oneCount, err := models.GetCounts(condition, params, 6)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取一个月VIP累计收入失败", err.Error(), c.Ctx.Input)
	}
	//两个月会员累计收入
	twoCount, err := models.GetCounts(condition, params, 7)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取两个月VIP累计收入失败", err.Error(), c.Ctx.Input)
	}
	VIPCount := oneCount + twoCount

	//累计总收入
	totalCount := newCount + loanCount + platformCount + VIPCount

	//数据明细
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}

	//交易类型
	Source := c.GetString("Source")
	if Source != "" {
		condition += " AND u.deal_type IN (" + Source + ")"
	}
	if Source == "" {
		condition += " AND u.deal_type IN (2,3,5,6,7)"
	}
	//开始时间
	if startTime := c.GetString("startTime"); startTime != "" {
		condition += " AND u.create_time>=?"
		params = append(params, startTime)
	}
	//结束时间
	if endTime := c.GetString("endTime"); endTime != "" {
		condition += " AND DATE_ADD(u.create_time, INTERVAL -1 DAY)<=?"
		params = append(params, endTime)
	}

	//查询
	IncomeList, err := models.GetIncomeList(condition, params, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取服务收益数据明细失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取数据总数失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	//总计人数
	usersCount, err := models.GetUsersCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取人数失败", err.Error(), c.Ctx.Input)
	}
	//总计收入金额
	totalMoney, err := models.GetTotalMoney(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取收入金额失败", err.Error(), c.Ctx.Input)
	}
	money, err := models.GetMoney(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取失败或执行中金额失败", err.Error(), c.Ctx.Input)
	}
	totalMoney = totalMoney - money

	c.Data["count"] = count             //总数据数
	c.Data["usersCount"] = usersCount   //总计人数
	c.Data["totalMoney"] = totalMoney   //收入金额
	c.Data["IncomeList"] = IncomeList   //数据明细
	c.Data["pageNum"] = pageNum         //当前页
	c.Data["pageCount"] = pageCount     //总页数
	c.Data["newCount"] = newCount       //新口子累计收入
	c.Data["loanCount"] = loanCount      //网贷征信口子累计收入
	c.Data["platformCount"] = platformCount //平台征信口子累计收入
	c.Data["vipCount"] = VIPCount       //VIP累计收入
	c.Data["totalCount"] = totalCount   //累计总收入
	c.TplName = "data-operation/service_revenue.html"
}

//获取活动消耗页面数据
//@router /getActivityUseData [get]
func (c *DataController) GetActivityUseData() {
	//设置整体加载
	c.IsNeedTemplate()
	condition := ""
	params := make([]interface{}, 0, 2)

	//一元话费消耗数
	oneCount, err := models.GetOneBillCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取一元话费消耗数失败", err.Error(), c.Ctx.Input)
	}
	//十元话费消耗数
	tenCount, err := models.GetTenBillCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取十元话费消耗数失败", err.Error(), c.Ctx.Input)
	}
	billCount := 1*oneCount + 10*tenCount

	//抽奖奖励现金消耗
	moneyCount, err := models.GetMoneyCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取抽奖奖励现金消耗失败", err.Error(), c.Ctx.Input)
	}

	//累计总消耗
	totalCount := billCount + moneyCount

	//数据明细
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	//开始时间
	if startTime := c.GetString("startTime"); startTime != "" {
		condition += " AND u.create_time>=?"
		params = append(params, startTime)
	}
	//结束时间
	if endTime := c.GetString("endTime"); endTime != "" {
		condition += " AND DATE_ADD(u.create_time, INTERVAL -1 DAY)<=?"
		params = append(params, endTime)
	}

	Source := c.GetString("Source")
	if Source == "1,4" {
		condition += " AND u.service_states = 0 AND u.deal_type IN (" + Source + ") "
		totalMoney, err := models.GetBillCount(condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取指定时间内话费数失败", err.Error(), c.Ctx.Input)
		}
		c.Data["totalMoney"] = totalMoney //话费支出金额
	} else if Source == "13" {
		condition += " AND u.deal_type IN (" + Source + ") AND u.money_amount > 0"
		totalMoney, err := models.GetTotalMoney(condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取现金支出金额失败", err.Error(), c.Ctx.Input)
		}
		c.Data["totalMoney"] = totalMoney //现金支出金额
	} else {
		totalMoney, err := models.GetAllActivitymoney(condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有项目支出金额失败", err.Error(), c.Ctx.Input)
		}
		c.Data["totalMoney"] = totalMoney //所有项目支出金额
	}
	//获取活动消耗数据明细
	activityList, err := models.GetActivityList(Source,condition, params, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取活动明细列表失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetActivityCount(Source,condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取数据总数失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	//总计人数
	usersCount, err := models.GetActivityUsersCount(Source,condition,params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取人数失败", err.Error(), c.Ctx.Input)
	}

	c.Data["count"] = count               //数据量
	c.Data["pageCount"] = pageCount       //总页数
	c.Data["usersCount"] = usersCount     //人数
	c.Data["pageNum"] = pageNum           //当前页
	c.Data["activityList"] = activityList //活动消耗数据明细
	c.Data["billCount"] = billCount       //话费消耗
	c.Data["moneyCount"] = moneyCount     //现金消耗
	c.Data["totalCount"] = totalCount     //累计总消耗
	c.TplName = "data-operation/activity_consume.html"
}

