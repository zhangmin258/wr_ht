package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"strings"
	"strconv"
)

/*
*H5产品结算接口
 */
type ProCleaningController struct {
	BaseController
}

func (c *ProCleaningController) ProCleaningPage() {
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	//获取默认产品
	product, err := models.GetProductH5IdAndNameFirst()
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品异常！", err.Error(), c.Ctx.Input)
	}
	//默认产品的结算信息详情
	proCleaning, err := models.GetProCleaningDay(product.Id, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询产品结算信息异常！", err.Error(), c.Ctx.Input)
	}
	cleaningData, err := models.GetAllProCleaning()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有产品结算信息异常！", err.Error(), c.Ctx.Input)
	}

	//已结算和未结算产品的总金额
	var settlementCount, notSettlementCount float64
	for _, v := range cleaningData {
		if v.State == 0 {
			notSettlementCount += v.ProPrice
		} else {
			settlementCount += v.ProPrice
		}
	}
	count, err := models.GetProCleaningCount(product.Id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["count"] = count                           //数据总共多少条
	c.Data["pageCount"] = pageCount                   //总共多少页
	c.Data["pageNum"] = pageNum                       //第几页
	c.Data["settlementCount"] = settlementCount       //已结费用
	c.Data["notSettlementCount"] = notSettlementCount //未结费用
	c.Data["product"] = product
	c.Data["proCleaning"] = proCleaning
	c.TplName = "settlement-management/products_settlement.html"
}

//数据明细
func (c *ProCleaningController) GetProCleanData() {
	// name := c.GetString("name")
	// id, err := models.GetPidByPname(name)
	// if err != nil && err.Error() != utils.ErrNoRow() {
	// 	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据产品名称查询id异常！", err.Error(), c.Ctx.Input)
	// }
	// if err.Error() == utils.ErrNoRow() {
	// 	c.Data["msg"] = "查无此产品！"
	// 	return
	// }
	id, _ := c.GetInt("id")
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	proCleaning, err := models.GetProCleaningDay(id, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询产品结算异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetProCleaningCount(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	for k, v := range proCleaning {
		switch v.JointMode {
		case 1:
			switch v.CpaDefine {
			case "注册":
				proCleaning[k].CpaMoney = float64(v.RegisterCount) * v.CpaPrice
			case "认证":
				proCleaning[k].CpaMoney = float64(v.ApplyCount) * v.CpaPrice
			case "授信":
				proCleaning[k].CpaMoney = float64(v.CreditCount) * v.CpaPrice
			case "申请借款":
				proCleaning[k].CpaMoney = float64(v.ApplyLoanCount) * v.CpaPrice
			case "放款":
				proCleaning[k].CpaMoney = float64(v.MakeLoanCount) * v.CpaPrice
			}
		case 2:
			proCleaning[k].CpaMoney = float64(v.MakeLoanAmount) * v.CpsPrice
		case 3:
			switch v.CpaDefine {
			case "注册":
				proCleaning[k].CpaMoney = float64(v.RegisterCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "认证":
				proCleaning[k].CpaMoney = float64(v.ApplyCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "授信":
				proCleaning[k].CpaMoney = float64(v.CreditCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "申请借款":
				proCleaning[k].CpaMoney = float64(v.ApplyLoanCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "放款":
				proCleaning[k].CpaMoney = float64(v.MakeLoanCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			}
		}
	}
	c.Data["id"] = id
	c.Data["count"] = count         //数据总共多少条
	c.Data["pageCount"] = pageCount //总共多少页
	c.Data["pageNum"] = pageNum
	c.Data["flag"] = 1
	c.Data["proCleaning"] = proCleaning
	c.TplName = "settlement-management/products_settlement_info.html"
}

//数据明细
func (c *ProCleaningController) GetProCleanDataByName() {
	name := c.GetString("name")
	flag := 1 //是否有该产品
	//根据产品名称查询产品id
	id, err := models.GetProductIdByName(name)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			flag = 0
		} else {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据name查询产品id异常！", err.Error(), c.Ctx.Input)
		}
	}
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	proCleaning, err := models.GetProCleaningDay(id, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询产品结算异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetProCleaningCount(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	for k, v := range proCleaning {
		switch v.JointMode {
		case 1:
			switch v.CpaDefine {
			case "注册":
				proCleaning[k].CpaMoney = float64(v.RegisterCount) * v.CpaPrice
			case "认证":
				proCleaning[k].CpaMoney = float64(v.ApplyCount) * v.CpaPrice
			case "授信":
				proCleaning[k].CpaMoney = float64(v.CreditCount) * v.CpaPrice
			case "申请借款":
				proCleaning[k].CpaMoney = float64(v.ApplyLoanCount) * v.CpaPrice
			case "放款":
				proCleaning[k].CpaMoney = float64(v.MakeLoanCount) * v.CpaPrice
			}
		case 2:
			proCleaning[k].CpaMoney = float64(v.MakeLoanAmount) * v.CpsPrice
		case 3:
			switch v.CpaDefine {
			case "注册":
				proCleaning[k].CpaMoney = float64(v.RegisterCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "认证":
				proCleaning[k].CpaMoney = float64(v.ApplyCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "授信":
				proCleaning[k].CpaMoney = float64(v.CreditCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "申请借款":
				proCleaning[k].CpaMoney = float64(v.ApplyLoanCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			case "放款":
				proCleaning[k].CpaMoney = float64(v.MakeLoanCount)*v.CpaPrice + float64(v.MakeLoanAmount)*v.CpsPrice
			}
		}
	}
	c.Data["id"] = id
	c.Data["count"] = count         //数据总共多少条
	c.Data["pageCount"] = pageCount //总共多少页
	c.Data["pageNum"] = pageNum
	c.Data["proCleaning"] = proCleaning
	c.Data["flag"] = flag
	c.TplName = "settlement-management/products_settlement_info.html"
}

//结算信息、结算历史、开票信息、后台信息
func (c *ProCleaningController) GetProCleanDetail() {
	id, err := c.GetInt("pid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取id错误！", err.Error(), c.Ctx.Input)
	}
	flag, err := c.GetInt("flag") //1结算信息2结算历史3开票信息4后台信息
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取flag错误！", err.Error(), c.Ctx.Input)
	}
	switch flag {
	case 1: //1结算信息
		cleaning, err := models.GetProCleaningNow(id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询产品结算信息异常！", err.Error(), c.Ctx.Input)
		}
		c.Data["cleaning"] = cleaning
		c.TplName = "settlement-management/settlement_info.html"
	case 2:                                 //2结算历史
		pageNum, err := c.GetInt("page", 1) //分页信息（第几页）
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取分页信息（第几页）错误！", err.Error(), c.Ctx.Input)
		}
		cleanHisTory, err := models.GetCleanHisByPid(id, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算历史错误！", err.Error(), c.Ctx.Input)
		}
		count, err := models.CleanHisByPidCount(id)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算历史总数错误！", err.Error(), c.Ctx.Input)
		}
		pageCount, err := utils.GetPageCount(count, utils.PageSize20)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
		}
		c.Data["pid"] = id
		c.Data["pageNum"] = pageNum
		c.Data["count"] = count
		c.Data["pageCount"] = pageCount //总共多少页
		c.Data["cleanHisTory"] = cleanHisTory
		c.TplName = "settlement-management/settlement_history.html"
	case 3: //3开票信息
		bilInfor, err := models.GetBilInformation(id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取开票信息失败", err.Error(), c.Ctx.Input)
		}
		c.Data["bilInfor"] = bilInfor
		c.TplName = "settlement-management/invoice_info.html"
	case 4: //4后台信息
		backGround, err := models.GetBilInformation(id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取后台信息失败", err.Error(), c.Ctx.Input)
		}
		c.Data["backGround"] = backGround
		c.TplName = "settlement-management/background_info.html"
	}
}

//获取商品列表
func (c *ProCleaningController) GetH5ProductList() {
	defer c.ServeJSON()
	name := c.GetString("name")
	condition := " and name like ?"
	param := "%" + name + "%"
	products, err := models.GetProductH5IdAndName(condition, param)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询产品异常"}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品异常！", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "products": products}
	return
}

//撤销结算
func (c *ProCleaningController) CancelSettlement() {
	defer c.ServeJSON()
	clId,err:=c.GetInt("clId")
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取清算历史id异常"}
		return
	}
	dailyDataId := c.GetString("dailyDataId")
	ids:=strings.Split(dailyDataId,",")
	var idsNew []int
	//更新该时间段的数据
	for _,v:=range ids{
		i,err:=strconv.Atoi(v)
		if err!=nil {
			c.Data["json"] = map[string]interface{}{"ret": 403, "err": "dailyDataId转换异常"}
			return
		}
		idsNew=append(idsNew,i)
	}
	err = models.CancelSettlement(clId,idsNew)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "更新指定时间段的数据异常"}
		return
	}
	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "撤销结算！", "结算历史ids:"+dailyDataId, c.Ctx.Input)
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "撤销结算成功"}
}

//查看结算历史明细
func (c *ProCleaningController)ShowSettlementHistory()  {
	c.IsNeedTemplate()
	pid, err := c.GetInt("pid")
	if err!=nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品id异常！", err.Error(), c.Ctx.Input)
	}
	clId, err := c.GetInt("clId")
	if err!=nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品id异常！", err.Error(), c.Ctx.Input)
	}
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	condition := ""
	params := []string{}
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
		}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品信息异常！", err.Error(), c.Ctx.Input)
	}
	cleaningData, err := models.GetProSettle(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息异常！", err.Error(), c.Ctx.Input)
	}
	//获取指定时间各项数据总数
	cleanCount, err := models.GetCleanedCount(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.OneSettleCleaning(pid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结算信息总数异常！", err.Error(), c.Ctx.Input)
	}
	//获取开票信息
	remark,err:=models.GetCleaningHistoryInfo(clId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取开票信息异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["pid"] = pid
	c.Data["proInfo"] = proInfo
	c.Data["cleanCount"] = cleanCount
	c.Data["count"] = count
	c.Data["remark"] = remark
	c.Data["cleaningData"] = cleaningData
	c.TplName = "settlement-management/show_settlement.html"
}
