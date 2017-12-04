package controllers

/*
*结算接口
 */
import (
	"strconv"
	"strings"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type SettlementController struct {
	BaseController
}

//结算页面
func (c *SettlementController) SettlePage() {
	c.IsNeedTemplate()
	pid, _ := c.GetInt("pid")
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	condition := ""
	params := []string{}
	if startDate != "" && endDate == "" {
		d, ok := utils.MaxDay(startDate)
		if ok {
			endDate = d
		}
	} else if startDate == "" && endDate != "" {
		d := strings.SplitAfter(endDate, "-")
		if len(d) > 1 {
			startDate += d[0]
			startDate += d[1]
			startDate += "01"
		}
	} else if startDate == "" && endDate == "" {
		minTime, err := models.GetMinTime(pid, condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取最小时间异常！", err.Error(), c.Ctx.Input)
		}
		startDate = minTime
		d, ok := utils.MaxDay(minTime)
		if ok {
			endDate = d
		}
	}
	if startDate != "" && endDate != "" {
		condition += " AND dd.date>=? AND dd.date<=?"
		params = append(params, startDate)
		params = append(params, endDate)
	}
	//获取产品标题
	proInfo, err := models.GetProInfoByPid(pid)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			proInfo = new(models.ProductInfo)
		} else {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品信息异常！", err.Error(), c.Ctx.Input)
		}
	}
	cleaningData, err := models.GetProNotSettle(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息异常！", err.Error(), c.Ctx.Input)
	}
	//获取指定时间内的最小时间
	minDate, err := models.GetMinTime(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取最小时间异常！", err.Error(), c.Ctx.Input)
	}
	if minDate != "" {
		startDate = minDate
	}
	//获取指定时间各项数据总数
	cleanCount, err := models.GetCleanCount(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.OneNotSettleCleaning(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["pid"] = pid
	c.Data["startDate"] = startDate
	c.Data["endDate"] = endDate
	c.Data["proInfo"] = proInfo
	c.Data["cleanCount"] = cleanCount
	c.Data["count"] = count
	c.Data["cleaningData"] = cleaningData
	c.TplName = "settlement-management/settlement_detail.html"
}

//结算产品
func (c *SettlementController) SettlePro() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	remark := c.GetString("remark")
	isBilling, err := c.GetInt("isBilling")
	pid, _ := c.GetInt("pid")
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	condition := ""
	params := []string{}
	if startDate != "" && endDate != "" {
		condition += " AND dd.date>=? AND dd.date<=?"
		params = append(params, startDate)
		params = append(params, endDate)
	} else {
		resultMap["err"] = "获取时间失败"
		resultMap["ret"] = 403
		return
	}
	//获取指定时间未结算数据的id
	ids, err := models.GetDailyDataId(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取结算信息异常"
		resultMap["ret"] = 403
		return
	}
	DailyDataIds := ""
	for k, v := range ids {
		if k != len(ids)-1 {
			DailyDataIds = DailyDataIds + strconv.Itoa(v) + ","
		} else {
			DailyDataIds = DailyDataIds + strconv.Itoa(v)
		}
	}
	//获取指定时间内未结算总数
	cleanCount, err := models.GetCleanCount(pid, condition, params)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
		}
		resultMap["err"] = "保存结算信息异常"
		resultMap["ret"] = 403
	}
	//计算指定时间内的有效数据或者放款金额
	quantityCount, err := models.GetQuantityCount(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取有效数据异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "保存结算信息异常"
		resultMap["ret"] = 403
	}
	switch quantityCount.CpaDefine {
	case "注册":
		if quantityCount.JointMode == 1 {
			cleanCount.QuantityCount = quantityCount.RegisterCount
		} else if quantityCount.JointMode == 2 {
			cleanCount.MakeLoanAmount = quantityCount.MakeLoanAmount
		} else {
			cleanCount.QuantityCount = quantityCount.RegisterCount
		}
	case "认证":
		if quantityCount.JointMode == 1 {
			cleanCount.QuantityCount = quantityCount.ApplyCount
		} else if quantityCount.JointMode == 2 {
			cleanCount.MakeLoanAmount = quantityCount.MakeLoanAmount
		} else {
			cleanCount.QuantityCount = quantityCount.ApplyCount
		}
	case "授信":
		if quantityCount.JointMode == 1 {
			cleanCount.QuantityCount = quantityCount.CreditCount
		} else if quantityCount.JointMode == 2 {
			cleanCount.MakeLoanAmount = quantityCount.MakeLoanAmount
		} else {
			cleanCount.QuantityCount = quantityCount.CreditCount
		}
	case "申请借款":
		if quantityCount.JointMode == 1 {
			cleanCount.QuantityCount = quantityCount.ApplyLoanCount
		} else if quantityCount.JointMode == 2 {
			cleanCount.MakeLoanAmount = quantityCount.MakeLoanAmount
		} else {
			cleanCount.QuantityCount = quantityCount.ApplyLoanCount
		}
	case "放款":
		if quantityCount.JointMode == 1 {
			cleanCount.QuantityCount = quantityCount.MakeLoanCount
		} else if quantityCount.JointMode == 2 {
			cleanCount.MakeLoanAmount = quantityCount.MakeLoanAmount
		} else {
			cleanCount.QuantityCount = quantityCount.MakeLoanCount
		}
	}
	//保存结算信息
	err = models.SaveCleaningHistory(pid, isBilling, DailyDataIds, remark, condition, cleanCount, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存结算信息异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "保存结算信息异常"
		resultMap["ret"] = 403
		return
	}
	resultMap["msg"] = "保存结算信息成功"
	resultMap["ret"] = 200
	return
}
