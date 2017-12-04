package controllers

import (
	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type CleanHisController struct {
	BaseController
}

func (c *CleanHisController) CleanHisPage() {
	c.IsNeedTemplate()
	pageNum1 := 1 //已结算分页信息（第几页）
	pageNum2 := 1 //未结算分页信息（第几页）
	condition := ""
	params := []string{}
	//已结算产品信息
	order := " ORDER BY ch.begin_time "
	cleanHistoryInfo, err := models.GetCleanHisByTime(order, condition, params, utils.StartIndex(pageNum1, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有已结算产品历史信息异常！", err.Error(), c.Ctx.Input)
	}
	for k, v := range cleanHistoryInfo {
		if v.JointMode == 1 {
			cleanHistoryInfo[k].ProPrice = strconv.FormatFloat(v.CpaPrice, 'f', -1, 64)
		} else if v.JointMode == 2 {
			cleanHistoryInfo[k].ProPrice = strconv.FormatFloat(v.CpsPrice, 'f', -1, 64)
		} else if v.JointMode == 3 {
			cleanHistoryInfo[k].ProPrice = strconv.FormatFloat(v.CpaPrice, 'f', -1, 64) + "+" + strconv.FormatFloat(v.CpsPrice, 'f', -1, 64) + "%"
		}
	}
	count1, err := models.GetCleanHisCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有已结算产品历史总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount1, err := utils.GetPageCount(count1, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	//未结算产品信息
	order = " ORDER BY date "
	cleaning, err := models.GetNotSettleProCleaning(order, condition, params, utils.StartIndex(pageNum1, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有未结算产品历史信息异常！", err.Error(), c.Ctx.Input)
	}
	for k, v := range cleaning {
		switch v.CpaDefine {
		case "注册":
			if v.JointMode == 1 {
				cleaning[k].QuantityCount = v.RegisterCount
			} else if v.JointMode == 2 {
				cleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				cleaning[k].QuantityCount = v.RegisterCount
			}
		case "认证":
			if v.JointMode == 1 {
				cleaning[k].QuantityCount = v.ApplyCount
			} else if v.JointMode == 2 {
				cleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				cleaning[k].QuantityCount = v.ApplyCount
			}
		case "授信":
			if v.JointMode == 1 {
				cleaning[k].QuantityCount = v.CreditCount
			} else if v.JointMode == 2 {
				cleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				cleaning[k].QuantityCount = v.CreditCount
			}
		case "申请借款":
			if v.JointMode == 1 {
				cleaning[k].QuantityCount = v.ApplyLoanCount
			} else if v.JointMode == 2 {
				cleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				cleaning[k].QuantityCount = v.ApplyLoanCount
			}
		case "放款":
			if v.JointMode == 1 {
				cleaning[k].QuantityCount = v.MakeLoanCount
			} else if v.JointMode == 2 {
				cleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				cleaning[k].QuantityCount = v.MakeLoanCount
			}
		}
	}
	count2, err := models.NotSettleCleaning(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有未结算产品历史总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount2, err := utils.GetPageCount(count2, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["pageCount1"] = pageCount1 //总共多少页
	c.Data["pageNum1"] = pageNum1     //第几页
	c.Data["count1"] = count1
	c.Data["cleanHistoryInfo"] = cleanHistoryInfo
	c.Data["pageCount2"] = pageCount2 //总共多少页
	c.Data["pageNum2"] = pageNum2     //第几页
	c.Data["count2"] = count2
	c.Data["cleaning"] = cleaning
	c.TplName = "settlement-management/settlement_monthSummary.html"
}

func (c *CleanHisController) Cleaned() {
	date := c.GetString("date")
	order := c.GetString("order")
	condition := ""
	params := []string{}
	if date != "" {
		endDate := ""
		date += "-01"
		d, ok := utils.MaxDay(date)
		if ok {
			endDate = d
		}
		condition += " AND ch.begin_time>=? AND ch.begin_time<=? "
		params = append(params, date)
		params = append(params, endDate)
	}
	switch order {
	case "1":
		order = " ORDER BY ch.begin_time "
	case "2":
		order = " ORDER BY ch.begin_time DESC "
	case "3":
		order = " ORDER BY ch.settle_money DESC "
	}
	pageNum, err := c.GetInt("page", 1) //分页信息（第几页）
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取分页信息（第几页）异常!", err.Error(), c.Ctx.Input)
	}
	cleanHistoryInfo, err := models.GetCleanHisByTime(order, condition, params, utils.StartIndex(pageNum, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有已结算产品历史信息异常！", err.Error(), c.Ctx.Input)
	}
	for k, v := range cleanHistoryInfo {
		if v.JointMode == 1 {
			cleanHistoryInfo[k].ProPrice = strconv.FormatFloat(v.CpaPrice, 'f', -1, 64)
		} else if v.JointMode == 2 {
			cleanHistoryInfo[k].ProPrice = strconv.FormatFloat(v.CpsPrice, 'f', -1, 64)
		} else if v.JointMode == 3 {
			cleanHistoryInfo[k].ProPrice = strconv.FormatFloat(v.CpaPrice, 'f', -1, 64) + "+" + strconv.FormatFloat(v.CpsPrice, 'f', -1, 64) + "%"
		}
	}
	count, err := models.GetCleanHisCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有已结算产品历史总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["pageCount"] = pageCount //总共多少页
	c.Data["pageNum"] = pageNum     //第几页
	c.Data["count"] = count
	c.Data["cleanHistoryInfo"] = cleanHistoryInfo
	c.TplName = "settlement-management/settlement_monthSummary_clean.html"
}

func (c *CleanHisController) NoClean() {
	date := c.GetString("date")
	order := c.GetString("order")
	if order == "" {
		order = c.GetString("order1")
	}
	condition := ""
	params := []string{}
	if date != "" {
		endDate := date + "-31"
		date += "-01"
		condition += " AND dd.date>=? AND dd.date<=? "
		params = append(params, date)
		params = append(params, endDate)
	}
	switch order {
	case "1":
		order = " ORDER BY date "
	case "2":
		order = " ORDER BY date DESC "
	case "3":
		order = "ORDER BY cpa_money DESC "
	}
	pageNum, err := c.GetInt("page", 1) //分页信息（第几页）
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取分页信息（第几页）异常!", err.Error(), c.Ctx.Input)
	}
	notSettleProCleaning, err := models.GetNotSettleProCleaning(order, condition, params, utils.StartIndex(pageNum, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有未结算产品历史信息异常！", err.Error(), c.Ctx.Input)
	}
	for k, v := range notSettleProCleaning {
		switch v.CpaDefine {
		case "注册":
			if v.JointMode == 1 {
				notSettleProCleaning[k].QuantityCount = v.RegisterCount
			} else if v.JointMode == 2 {
				notSettleProCleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				notSettleProCleaning[k].QuantityCount = v.RegisterCount
			}
		case "认证":
			if v.JointMode == 1 {
				notSettleProCleaning[k].QuantityCount = v.ApplyCount
			} else if v.JointMode == 2 {
				notSettleProCleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				notSettleProCleaning[k].QuantityCount = v.ApplyCount
			}
		case "授信":
			if v.JointMode == 1 {
				notSettleProCleaning[k].QuantityCount = v.CreditCount
			} else if v.JointMode == 2 {
				notSettleProCleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				notSettleProCleaning[k].QuantityCount = v.CreditCount
			}
		case "申请借款":
			if v.JointMode == 1 {
				notSettleProCleaning[k].QuantityCount = v.ApplyLoanCount
			} else if v.JointMode == 2 {
				notSettleProCleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				notSettleProCleaning[k].QuantityCount = v.ApplyLoanCount
			}
		case "放款":
			if v.JointMode == 1 {
				notSettleProCleaning[k].QuantityCount = v.MakeLoanCount
			} else if v.JointMode == 2 {
				notSettleProCleaning[k].MakeLoanAmount = v.MakeLoanAmount
			} else {
				notSettleProCleaning[k].QuantityCount = v.MakeLoanCount
			}
		}
	}
	count, err := models.NotSettleCleaning(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有已结算产品历史总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["pageCount"] = pageCount //总共多少页
	c.Data["pageNum"] = pageNum     //第几页
	c.Data["count"] = count
	c.Data["notSettleProCleaning"] = notSettleProCleaning
	c.TplName = "settlement-management/settlement_monthSummary_no_clean.html"
}
