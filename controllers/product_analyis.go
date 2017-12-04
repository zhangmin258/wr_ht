package controllers

import (
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type ProductAnalysisController struct {
	BaseController
}

//@router /getProductAnalysis [get]
func (c *ProductAnalysisController) GetProductAnalysis() {
	c.IsNeedTemplate()
	products, err := models.GetProductH5IdAndNameFirst()
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取第一个产品异常", err.Error(), c.Ctx.Input)
	}
	c.Data["productsName"] = products.Name
	c.Data["productsId"] = products.Id
	c.Data["productsType"] = products.CooperationType
	c.TplName = "agent-products/agent_data_analysis.html"
}

//各平台数据对比
//@router /getproductdetaildata [get]
func (c *ProductAnalysisController) GetProductDetailData() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	resultMap["ret"] = 403

	code := c.GetString("code")
	params := []string{}
	condition := ""
	//开始时间1
	if startDate1 := c.GetString("startDate"); startDate1 != "" {
		condition += " AND bl.create_time>=?"
		params = append(params, startDate1)
	}
	//结束时间1
	if endDate1 := c.GetString("endDate"); endDate1 != "" {
		condition += " AND bl.create_time<=?"
		params = append(params, endDate1+" 23:59:59")
	}
	//产品列表
	productlsit, err := models.GetAnalyisProduct("") //产品名称列表
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品列表失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取产品列表失败"
		return
	}

	//1:申请通过率，2：放款金额，3：获客成本，4：点击次数-注册人数，5：点击-注册转化率
	var AnalysisProduct []models.AnalysisProduct
	switch code {
	case "1":
		loancount, err := models.GetLoan_count(condition, params) //放款人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取放款人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取放款人数失败"
			return
		}
		applycount, err := models.GetApplycount(condition, params) //申请人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取申请人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取申请人数失败"
			return
		}
		h5loancount, err := models.GetH5Loan_count(params) //H5放款人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5放款人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取H5放款人数失败"
			return
		}
		h5applycount, err := models.GetH5Applycount(params) //H5申请人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5申请人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取H5申请人数失败"
			return
		}
		for k, _ := range h5loancount {
			loancount = append(loancount, h5loancount[k])
		}
		for k, _ := range h5applycount {
			applycount = append(applycount, h5applycount[k])
		}
		for _, v := range productlsit {
			var r models.AnalysisProduct
			x := 0
			y := 0
			r.Name = v.Name
			for _, v1 := range loancount {
				if v.Id == v1.Id {
					x = v1.Count
					break
				}
			}
			for _, v1 := range applycount {
				if v.Id == v1.Id {
					y = v1.Count
					break
				}
			}
			if y == 0 {
				r.ApplyRate = 0
			} else {
				r.ApplyRate = float64(x) / float64(y)
			}
			AnalysisProduct = append(AnalysisProduct, r)
		}
	case "2":
		loanmoney, err := models.GetLoanmoney(condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品放款金额失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取产品放款金额失败"
			return
		}
		h5loanmoney, err := models.GetH5Loanmoney(params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5产品放款金额失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取H5产品放款金额失败"
			return
		}
		for k, _ := range h5loanmoney {
			loanmoney = append(loanmoney, h5loanmoney[k])
		}
		for _, v := range productlsit {
			var r models.AnalysisProduct
			r.Name = v.Name
			for _, v1 := range loanmoney {
				if v.Id == v1.Id {
					r.Loan_money = v1.Money
					break
				}
			}
			AnalysisProduct = append(AnalysisProduct, r)
		}
	case "3":
		allloanmoney, err := models.GetAllloanmoney(condition, params) //总支出金额
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取总支出金额失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取总支出金额失败"
			return
		}
		allloancount, err := models.GetAllloancount(condition, params) //新老用户总放款人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取新老用户总放款人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取新老用户总放款人数失败"
			return
		}
		h5allloanmoney, err := models.GetH5Allloanmoney(params) //H5总支出金额
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5总支出金额失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取H5总支出金额失败"
			return
		}
		h5allloancount, err := models.GetH5Allloancount(params) //H5新老用户总放款人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5新老用户总放款人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取H5新老用户总放款人数失败"
			return
		}
		for k, _ := range h5allloanmoney {
			allloanmoney = append(allloanmoney, h5allloanmoney[k])
		}
		for k, _ := range h5allloancount {
			allloancount = append(allloancount, h5allloancount[k])
		}
		for _, v := range productlsit {
			var r models.AnalysisProduct
			x := 0.0
			y := 0
			r.Name = v.Name
			for _, v1 := range allloanmoney {
				if v.Id == v1.Id {
					x = v1.Money
					break
				}
			}
			for _, v1 := range allloancount {
				if v.Id == v1.Id {
					y = v1.Count
					break
				}
			}
			if y == 0 {
				r.Cac = 0
			} else {
				r.Cac = float64(x) / float64(y)
			}
			AnalysisProduct = append(AnalysisProduct, r)
		}
	/*case "4":
	cradit, err := models.GetCredit(condition, params) //信用
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取信用评分失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取信用评分失败"
		return
	}
	for _, v := range productlsit {
		var r models.AnalysisProduct
		r.Name = v.Name
		for _, v1 := range cradit {
			if v.Id == v1.Id {
				r.Credit = v1.Money
				break
			}
		}
		AnalysisProduct = append(AnalysisProduct, r)
	}*/
	/*case "5":
	loanmoney, err := models.GetLoanmoney(condition, params) //放款金额
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取放款金额失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取放款金额失败"
		return
	}
	loan_count, err := models.GetLoan_count(condition, params) //放款人数
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取放款人数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取放款人数失败"
		return
	}
	h5loancount, err := models.GetH5Loan_count(params) //H5放款人数
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5放款人数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取H5放款人数失败"
		return
	}
	h5loanmoney, err := models.GetH5Loanmoney(params) //H5放款金额
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取H5放款金额失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取H5放款金额失败"
		return
	}
	for k, _ := range h5loanmoney {
		loanmoney = append(loanmoney, h5loanmoney[k])
	}
	for k, _ := range h5loancount {
		loan_count = append(loan_count, h5loancount[k])
	}
	for _, v := range productlsit {
		var r models.AnalysisProduct
		x := 0.0
		y := 0
		r.Name = v.Name
		for _, v1 := range loanmoney {
			if v.Id == v1.Id {
				x = v1.Money
				break
			}
		}
		for _, v1 := range loan_count {
			if v.Id == v1.Id {
				y = v1.Count
				break
			}
		}
		if y == 0 {
			r.Amount = 0
		} else {
			r.Amount = float64(x) / float64(y)
		}
		AnalysisProduct = append(AnalysisProduct, r)
	}*/
	case "5":
		hits, err := models.GetHits(params) //点击次数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取点击次数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取点击次数失败"
			return
		}
		register, err := models.GetRegister(params) //注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取注册人数失败"
			return
		}
		platformRegister, err := models.GetPlatformRegister(params) //平台注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取平台注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取平台注册人数失败"
			return
		}
		for _, v := range productlsit {
			var r models.AnalysisProduct
			x := 0 //点击次数
			y := 0 //我司统计注册
			z := 0 //平台注册
			r.Name = v.Name
			for _, v1 := range hits {
				if v.Id == v1.Id {
					x = v1.Count
					break
				}
			}
			for _, v1 := range register {
				if v.Id == v1.Id {
					if v.Type == 0 {
						y = v1.Count
					}
					break
				}
			}
			for _, v1 := range platformRegister {
				if v.Id == v1.Id {
					if v.Type == 1 {
						z = v1.Count
					}
					break
				}
			}
			if x == 0 {
				r.HitsRegister = 0
				r.PlatformHitsRegister = 0
			} else {
				r.HitsRegister = float64(y) / float64(x)
				r.PlatformHitsRegister = float64(z) / float64(x)
			}
			AnalysisProduct = append(AnalysisProduct, r)
		}
	case "4":
		hits, err := models.GetHits(params) //点击次数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取点击次数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取点击次数失败"
			return
		}
		register, err := models.GetRegister(params) //注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取注册人数失败"
			return
		}
		platformRegister, err := models.GetPlatformRegister(params) //平台注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取平台注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取平台注册人数失败"
			return
		}
		for _, v := range productlsit {
			var r models.AnalysisProduct
			r.Name = v.Name
			for _, v1 := range hits {
				if v.Id == v1.Id {
					r.HitsCount = v1.Count
				}
			}
			for _, v1 := range register {
				if v.Id == v1.Id {
					r.RegisterCount = v1.Count
				}
			}
			for _, v1 := range platformRegister {
				if v.Id == v1.Id {
					r.PlatformRegister = v1.Count
				}
			}
			AnalysisProduct = append(AnalysisProduct, r)
		}

	}
	resultMap["ret"] = 200
	resultMap["msg"] = AnalysisProduct
}

//@router /getTerraceProducts [get]
func (this *ProductAnalysisController) GetTerraceProducts() {
	defer this.ServeJSON()
	name := this.GetString("name")
	condition := " AND name like ?"
	params := "%" + name + "%"
	products, err := models.GetProductsInfo(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取API和H5产品异常", err.Error(), this.Ctx.Input)
	}
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 200
	resultMap["products"] = products
	this.Data["json"] = resultMap
}

//平台数据明细
//@router /getplatformdata [get]
func (c *ProductAnalysisController) GetPlatformData() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	resultMap["ret"] = 403
	code := c.GetString("code")
	params := []string{}
	condition := ""
	//产品列表
	productlsit, err := models.GetAnalyisProduct("") //产品名称列表
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品列表失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取产品列表失败"
		return
	}
	//代理产品
	condition += " AND bl.product_id = ?"
	name := c.GetString("name")
	flag := 1 //是否有该产品
	//根据产品名称查询产品id
	productId, err := models.GetProductIdByName(name)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			flag = 0
		} else {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据name查询产品id异常！", err.Error(), c.Ctx.Input)
		}
	}
	//productId, _ := c.GetInt("productId")
	if productId == 0 {
		productId = productlsit[0].Id
	}
	//
	// 获取产品类型    0:Api,1:H5
	proType, err := models.GetProductCooperationTypeById(productId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品合作类型错误", err.Error(), c.Ctx.Input)
	}
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	st, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取时间错误", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取时间错误"
		return
	}
	et, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取时间错误", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取时间错误"
		return
	}
	//开始时间2
	if startDate != "" {
		condition += " AND bl.create_time>=? "
		params = append(params, startDate)
	}
	//结束时间2
	if endDate != "" {
		condition += " AND DATE_FORMAT(bl.create_time,'%Y-%m-%d') <=? "
		params = append(params, endDate)
	}
	timediff := et.Sub(st)
	res2 := make([]models.AnalysisProduct, int(timediff.Hours()/24)+1)
	x := st
	for k, _ := range res2 {
		res2[k].Data = x
		x = x.AddDate(0, 0, 1)
	}
	switch code {
	case "1":
		var loan_count []models.ProductCountDate
		var applycount []models.ProductCountDate
		if proType == 0 {
			loan_count, err = models.GetLoan_countByTime(condition, productId, params) //放款人数
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取放款人数失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取放款人数失败"
			}
			applycount, err = models.GetApplycountByTime(condition, productId, params) //申请人数
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取申请人数失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取申请人数失败"
			}
		}
		if proType == 1 {
			loan_count, err = models.GetH5Loan_countByTime(productId, params) //放款人数
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取放款人数失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取放款人数失败"
			}
			applycount, err = models.GetH5ApplycountByTime(productId, params) //申请人数
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取申请人数失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取申请人数失败"
			}
		}
		for k, v := range res2 {
			x := 0
			y := 0
			for _, v1 := range loan_count {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					x = v1.Count
				}
			}
			for _, v1 := range applycount {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					y = v1.Count
				}
			}
			if y == 0 {
				res2[k].ApplyRate = 0
			} else {
				res2[k].ApplyRate = float64(x) / float64(y)
			}
		}
	case "2":
		var loanmoney []models.ProductMoneyData
		if proType == 0 {
			loanmoney, err = models.GetLoanmoneyByTime(condition, productId, params)
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品按时间取放款金额失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取产品按时间取放款金额失败"
			}
		}
		if proType == 1 {
			loanmoney, err = models.GetH5LoanmoneyByTime(productId, params)
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取按时间获取H5产品放款金额失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取按时间获取H5产品放款金额失败"
			}
		}
		for k, v := range res2 {
			for _, v1 := range loanmoney {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					res2[k].Loan_money = v1.Money
				}
			}
		}
	case "3":
		var allloanmoney []models.ProductMoneyData
		var allloancount []models.ProductCountDate
		if proType == 0 {
			allloanmoney, err = models.GetAllloanmoneyByTime(condition, productId, params) //总支出金额
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取总支出金额失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取总支出金额失败"
			}
			allloancount, err = models.GetAllloancountByTime(condition, productId, params) //新老用户总放款人数
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取新老用户总放款人数失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取新老用户总放款人数失败"
			}
		}
		if proType == 1 {
			allloanmoney, err = models.GetH5AllloanmoneyByTime(productId, params) //h5总支出金额
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取h5总支出金额失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取h5总支出金额失败"
			}
			allloancount, err = models.GetH5AllloancountByTime(productId, params) //h5新老用户总放款人数
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取h5新老用户总放款人数失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "获取h5新老用户总放款人数失败"
			}
		}
		for k, v := range res2 {
			x := 0.0
			y := 0
			for _, v1 := range allloanmoney {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					x = v1.Money
				}
			}
			for _, v1 := range allloancount {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					y = v1.Count
				}
			}
			if y == 0 {
				res2[k].Cac = 0
			} else {
				res2[k].Cac = float64(x) / float64(y)
			}
		}
	case "5":
		hits, err := models.GetHitsByTime(productId, params) //点击次数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取点击次数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取点击次数失败"
			return
		}
		register, err := models.GetRegisterByTime(productId, params) //我司统计注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取我司统计注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取我司统计注册人数失败"
			return
		}
		platformRegister, err := models.GetPlatformRegisterByTime(productId, params) //平台注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取平台注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取平台注册人数失败"
			return
		}
		for k, v := range res2 {
			for _, v1 := range hits {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					res2[k].HitsCount = v1.Count
				}
			}
			for _, v1 := range register {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					res2[k].RegisterCount = v1.Count
				}
			}
			for _, v1 := range platformRegister {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					res2[k].PlatformRegister = v1.Count
				}
			}
		}
	case "4":
		hits, err := models.GetHitsByTime(productId, params) //点击次数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取点击次数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取点击次数失败"
			return
		}
		register, err := models.GetRegisterByTime(productId, params) //我司统计注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取注册人数失败"
			return
		}
		platformRegister, err := models.GetPlatformRegisterByTime(productId, params) //平台注册人数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取平台注册人数失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取平台注册人数失败"
			return
		}
		for k, v := range res2 {
			x := 0
			y := 0
			z := 0
			for _, v1 := range hits {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					x = v1.Count
					break
				}
			}
			for _, v1 := range register {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					y = v1.Count
					break
				}
			}
			for _, v1 := range platformRegister {
				if v.Data.Format("2006-01-02") == v1.Data.Format("2006-01-02") {
					z = v1.Count
					break
				}
			}
			if x == 0 {
				res2[k].HitsRegister = 0
				res2[k].PlatformHitsRegister = 0
			} else {
				res2[k].HitsRegister = float64(y) / float64(x)
				res2[k].HitsRegister = float64(z) / float64(x)
			}
		}
	}
	resultMap["ret"] = 200
	resultMap["msg"] = res2
	resultMap["flag"] = flag
}

//获取产品列表
//@router /getproductlist [get]
func (c *ProductAnalysisController) GetProductlist() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	resultMap["ret"] = 403
	productlsit, err := models.GetAnalyisProduct("") //产品名称列表
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品列表失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取产品列表失败"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = productlsit
	return
}
