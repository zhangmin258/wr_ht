package controllers

import (
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/file"

	"github.com/tealeg/xlsx"
)

/**
产品数据模块控制器
*/
type ProductDataController struct {
	BaseController
}

//@router /jumpDataEntryHTML [get]
func (this *ProductDataController) JumpDataEntryHTML() {
	this.IsNeedTemplate()
	productId, _ := this.GetInt("productId")
	agentId, _ := this.GetInt("agentId")
	//获取代理产品
	agentProductId, err := models.GetAgentProduct(agentId, productId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取中间表id失败！", err.Error(), this.Ctx.Input)
		this.Data["err"] = "获取代理产品失败！"
		return
	}
	agentProInfo, err := models.GetAgentProInfo(agentProductId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取代理产品失败！", err.Error(), this.Ctx.Input)
		this.Data["err"] = "获取代理产品失败！"
		return
	}
	this.Data["cleaningState"] = 0
	this.Data["agentProInfo"] = agentProInfo
	this.Data["ProductId"] = productId
	this.Data["AgentId"] = agentId
	this.Data["agentProductId"] = agentProductId
	this.TplName = "agent-products/data_entry.html"
}

//@router /showDetailedData [get]
func (this *ProductDataController) ShowDetailedData() {
	this.IsNeedTemplate()
	id, _ := this.GetInt("id", 0)
	agentProductId, _ := this.GetInt("agentProductId")
	createDate := this.GetString("createDate")
	daliyData, err := models.QueryAgentDailyDataById(id, createDate)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取代理产品失败！", err.Error(), this.Ctx.Input)
		}
		this.Data["err"] = "获取代理产品失败！"
		return
	}
	agentProInfo, err := models.GetAgentProInfo(agentProductId)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取代理产品失败！", err.Error(), this.Ctx.Input)
		}
		this.Data["err"] = "获取代理产品失败！"
		return
	}
	cleaningState, _ := this.GetInt("cleaningState")
	this.Data["cleaningState"] = cleaningState
	this.Data["agentProInfo"] = agentProInfo
	this.Data["Id"] = id
	this.Data["daliyData"] = daliyData
	this.Data["agentProductId"] = agentProductId
	this.TplName = "agent-products/data_entry.html"
}

//填充数据录入信息
//@router /getproductdateinfo [post]
func (this *ProductDataController) GetProductDateInfo() {
	resultMap := make(map[string]interface{})
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	resultMap["ret"] = 403
	agentProductId, _ := this.GetInt("agentProductId")
	startDate := this.GetString("startDate")
	dailyData, err := models.GetDailyDataByDate(agentProductId, startDate)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			resultMap["err"] = "数据库中没有数据，请填写数据信息！"
			return
		}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取产品运营信息失败！", err.Error(), this.Ctx.Input)
		resultMap["err"] = "获取产品运营信息失败！"
		return
	}
	resultMap["ret"] = 200
	resultMap["dailyData"] = dailyData
}

//更新H5的数据
//@router /updateProductData [post]
func (this *ProductDataController) UpdateProductData() {
	resultMap := make(map[string]interface{})
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	resultMap["ret"] = 403
	var proDailyData models.AgentDailyData
	dateTime := this.GetString("DateTime")
	//AgentProductId, _ := this.GetInt("AgentProductId")
	err := this.ParseForm(&proDailyData)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取代理产品失败！", err.Error(), this.Ctx.Input)
		resultMap["msg"] = "获取代理产品失败！"
		return
	}
	/*proInformation, err := models.GetProImformation(AgentProductId)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			resultMap["msg"] = "该产品结算信息未录入！"
			return
		}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据id查询产品结算信息异常！", err.Error(), this.Ctx.Input)
		resultMap["msg"] = "根据id查询产品结算信息异常！"
		return
	}*/

	switch proDailyData.JointMode {
	case 1:
		switch proDailyData.CpaDefine {
		case "注册":
			proDailyData.ProPrice = float64(proDailyData.RegisterCount) * proDailyData.CpaPrice
		case "认证":
			proDailyData.ProPrice = float64(proDailyData.ApplyCount) * proDailyData.CpaPrice
		case "授信":
			proDailyData.ProPrice = float64(proDailyData.CreditCount) * proDailyData.CpaPrice
		case "申请借款":
			proDailyData.ProPrice = float64(proDailyData.ApplyLoanCount) * proDailyData.CpaPrice
		case "放款":
			proDailyData.ProPrice = float64(proDailyData.MakeLoanCount) * proDailyData.CpaPrice
		}
	case 2:
		proDailyData.ProPrice = float64(proDailyData.MakeLoanAmount) * proDailyData.CpsFirstPer / 100
	case 3:
		switch proDailyData.CpaDefine {
		case "注册":
			proDailyData.ProPrice = float64(proDailyData.RegisterCount)*proDailyData.CpaPrice + float64(proDailyData.MakeLoanAmount)*proDailyData.CpsFirstPer/100
		case "认证":
			proDailyData.ProPrice = float64(proDailyData.ApplyCount)*proDailyData.CpaPrice + float64(proDailyData.MakeLoanAmount)*proDailyData.CpsFirstPer/100
		case "授信":
			proDailyData.ProPrice = float64(proDailyData.CreditCount)*proDailyData.CpaPrice + float64(proDailyData.MakeLoanAmount)*proDailyData.CpsFirstPer/100
		case "申请借款":
			proDailyData.ProPrice = float64(proDailyData.ApplyLoanCount)*proDailyData.CpaPrice + float64(proDailyData.MakeLoanAmount)*proDailyData.CpsFirstPer/100
		case "放款":
			proDailyData.ProPrice = float64(proDailyData.MakeLoanCount)*proDailyData.CpaPrice + float64(proDailyData.MakeLoanAmount)*proDailyData.CpsFirstPer/100
		}
	}

	if proDailyData.Id == 0 {
		err := models.InsertAgentDailyData(proDailyData, dateTime)

		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "修改产品运营信息失败！", err.Error(), this.Ctx.Input)
			resultMap["msg"] = "添加产品运营信息失败！"
			return
		}
	} else {
		err := models.UpdateAgentDailyData(proDailyData, dateTime)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "修改产品运营信息失败！", err.Error(), this.Ctx.Input)
			resultMap["msg"] = "修改产品运营信息失败！"
			return
		}
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "记录产品运营信息成功！"
}

/**
根据产品名字查询其代理商
*/
//@router /getAgentsByProductId [get]
func (this *ProductDataController) GetAgentsByProductId() {
	defer this.ServeJSON()
	name := this.GetString("name")
	flag := 1 //是否有该产品
	//根据产品名称查询产品id
	productId, err := models.GetProductIdByName(name)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			flag = 0
		}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据name查询产品id异常！", err.Error(), this.Ctx.Input)
	}
	//productId, _ := this.GetInt("ProductId")
	cooperationType, _ := this.GetInt("CooperationType")
	agents, err := models.QueryAgentsByProductId(productId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "数据不存在", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "数据不存在"}
	} else if cooperationType == 1 {
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": agents, "flag": flag}
	} else {
		agents = nil
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": agents, "flag": flag}
	}
}

//@router /getProductDataList [get]
func (this *ProductDataController) GetProductDataList() {
	this.IsNeedTemplate()
	products, err := models.GetProductH5IdAndNameFirst()
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取第一个h5产品异常", err.Error(), this.Ctx.Input)
	}
	this.Data["productsName"] = products.Name
	this.Data["productsId"] = products.Id
	this.Data["productsType"] = products.CooperationType
	this.TplName = "agent-products/product_data.html"
}

//@router /getProducts [get]
func (this *ProductDataController) GetProducts() {
	defer this.ServeJSON()
	name := this.GetString("name")
	condition := " AND name like ? "
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
func (this *ProductDataController) GetProductDataInfo() {
	defer this.ServeJSON()
	//this.IsNeedTemplate()
	//选择合作产品
	//products, _ := models.GetProductsInfo()
	//this.Data["products"] = products
	name := this.GetString("name")
	flag := 1 //是否有该产品
	//根据产品名称查询产品id
	productId, err := models.GetProductIdByName(name)
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			flag = 0
			productId = 1
		}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据name查询产品id异常！", err.Error(), this.Ctx.Input)
	}
	//productId, _ := this.GetInt("productId", 1)
	cooperation_type, err := models.GetProductCooperationById(productId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据productId获取cooperation_type失败", err.Error(), this.Ctx.Input)
	}
	condition := ""
	params := make([]interface{}, 0, 2)
	var registerCount, authCount, creditExtensionCount, loanApplyCount, creditCount int
	//var err error
	if cooperation_type == 0 { //API
		//数据汇总(累计注册用户、累计申请用户、累计授信用户、累计申请借款用户、累计放款用户)
		//累计注册用户
		registerCount, err = models.GetRegisterCount(productId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计注册数据失败", err.Error(), this.Ctx.Input)
		}
		//累计认证
		authCount, err = models.GetAuthCount(productId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计认证数据失败", err.Error(), this.Ctx.Input)
		}
		/*//累计申请用户
		applyCount, err = models.GetApplyCount()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计申请数据失败", err.Error(), this.Ctx.Input)
		}*/
		//累计授信用户
		creditExtensionCount, err = models.GetCreditExtensionCount(productId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计授信数据失败", err.Error(), this.Ctx.Input)
		}
		//累计申请借款用户
		loanApplyCount, err = models.GetLoanTotal(productId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计申请借款数据失败", err.Error(), this.Ctx.Input)
		}
		//累计放款用户
		creditCount, err = models.GetCreditCount(productId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取累计放款数据失败", err.Error(), this.Ctx.Input)
		}
	} else if cooperation_type == 1 { //H5

		defaultAgentId, err := models.GetFirstProductAgentIdByProductId(productId)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据产品id获取第一个产品代理商id失败", err.Error(), this.Ctx.Input)
			}
			defaultAgentId = -1
		}
		agentId, err := this.GetInt("agentId", defaultAgentId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取产品代理商id失败", err.Error(), this.Ctx.Input)
		}
		condition += " AND p.id = ?"
		params = append(params, productId)
		if agentId != -1 {
			condition += " AND a.id = ?"
			params = append(params, agentId)
		}
		//数据汇总(累计注册用户、累计申请用户、累计授信用户、累计申请借款用户、累计放款用户)
		countH5, err := models.GetCountForH5(condition, params)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取汇总数据失败", err.Error(), this.Ctx.Input)
		}
		registerCount = countH5.RegisterCount
		authCount = countH5.ApplyCount
		creditExtensionCount = countH5.CreditCount
		loanApplyCount = countH5.ApplyLoanCount
		creditCount = countH5.MakeLoanCount
	}
	resultMap := make(map[string]interface{})
	resultMap["registerCount"] = registerCount
	resultMap["authCount"] = authCount
	resultMap["creditExtensionCount"] = creditExtensionCount
	resultMap["loanApplyCount"] = loanApplyCount
	resultMap["creditCount"] = creditCount
	resultMap["flag"] = flag
	this.Data["json"] = resultMap
}

//产品数据-数据明细
func (this *ProductDataController) GetProductDataDetailList() {
	resultMap := make(map[string]interface{})
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	productId, _ := this.GetInt("productId", 1)
	pageNum, _ := this.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	cooperationType, err := models.GetProductCooperationById(productId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据productId获取cooperationType失败", err.Error(), this.Ctx.Input)
	}
	detailedDatas := make([]models.DetailedData, 0)
	defaultAgentId, err := models.GetFirstProductAgentIdByProductId(productId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "根据产品id获取第一个产品代理商id失败", err.Error(), this.Ctx.Input)
	}
	agentId, err := this.GetInt("agentId", defaultAgentId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取产品代理商id失败", err.Error(), this.Ctx.Input)
	}
	var pageCount, count int
	var params []int
	var condition = ""
	if cooperationType == 0 { //api
		//注册用户
		ru, err := models.GetRegisterUsersByCondition(productId, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取注册用户失败", err.Error(), this.Ctx.Input)
		}
		count, err = models.GetRegisterUsersCountByCondition(productId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取注册用户失败", err.Error(), this.Ctx.Input)
		}
		//认证用户
		lus, err := models.GetAuthsersByCondition(productId, "", nil)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取申请用户失败", err.Error(), this.Ctx.Input)
		}
		//授信用户
		ceus, err := models.GetCreditExtensionUsersByCondition(productId, "", nil)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取授信用户失败", err.Error(), this.Ctx.Input)
		}
		//申请借款用户
		ltus, err := models.GetLoanTotalUsersByCondition(productId, "", nil)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取申请借款用户失败", err.Error(), this.Ctx.Input)
		}
		//放款用户
		cus, err := models.GetCreditUsersByCondition(productId, "", nil)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取放款用户失败", err.Error(), this.Ctx.Input)
		}
		//放款金额
		cmus, err := models.GetCreditMoneyUsersByCondition(productId, "", nil)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取放款金额失败", err.Error(), this.Ctx.Input)
		}
		for _, v := range ru {
			dData := models.DetailedData{CreateDate: v.CreateDate, RegisterCount: v.Count}
			detailedDatas = append(detailedDatas, dData)
		}
		registerMap := make(map[string]int)
		for _, v := range ru {
			registerMap[v.CreateDate] = v.Count
		}
		loanMap := make(map[string]int)
		for _, v := range lus {
			loanMap[v.CreateDate] = v.Count
		}
		creditExtensionMap := make(map[string]int)
		for _, v := range ceus {
			creditExtensionMap[v.CreateDate] = v.Count
		}
		loanTotalMap := make(map[string]int)
		for _, v := range ltus {
			loanTotalMap[v.CreateDate] = v.Count
		}
		creditMap := make(map[string]int)
		for _, v := range cus {
			creditMap[v.CreateDate] = v.Count
		}
		creditMoneyMap := make(map[string]float32)
		for _, v := range cmus {
			creditMoneyMap[v.CreateDate] = v.Count
		}
		for k, v := range detailedDatas {
			if _, ok := loanMap[v.CreateDate]; ok {
				detailedDatas[k].ApplyCount = loanMap[v.CreateDate]
			} else {
				detailedDatas[k].ApplyCount = 0
			}
			if _, ok := creditExtensionMap[v.CreateDate]; ok {
				detailedDatas[k].CreditExtensionCount = creditExtensionMap[v.CreateDate]
			} else {
				detailedDatas[k].CreditExtensionCount = 0
			}
			if _, ok := loanTotalMap[v.CreateDate]; ok {
				detailedDatas[k].LoanApplyCount = loanTotalMap[v.CreateDate]
			} else {
				detailedDatas[k].LoanApplyCount = 0
			}
			if _, ok := creditMap[v.CreateDate]; ok {
				detailedDatas[k].CreditCount = creditMap[v.CreateDate]
			} else {
				detailedDatas[k].CreditCount = 0
			}
			if _, ok := creditMoneyMap[v.CreateDate]; ok {
				detailedDatas[k].CreditMoney = creditMoneyMap[v.CreateDate]
			} else {
				detailedDatas[k].CreditMoney = 0.0
			}
		}
		pageCount, _ = utils.GetPageCount(count, utils.PageSize20)
		resultMap["items"] = detailedDatas
	} else { //h5
		agentProId, err := models.GetAgentProduct(agentId, productId)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取下级代理商品失败", err.Error(), this.Ctx.Input)
			}
		} else {
			condition += " AND agent_product_id = ?"
			params = append(params, agentProId)
		}
		dailyDatas, err := models.GetCountForH5GroupByDate(condition, params, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
		count, err = models.GetDaiDataCount(condition, params)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取下级代理商品总数失败", err.Error(), this.Ctx.Input)
		}
		pageCount, _ = utils.GetPageCount(count, utils.PageSize20)
		resultMap["items"] = dailyDatas
	}
	resultMap["cooperationType"] = cooperationType
	resultMap["pageCount"] = pageCount
	resultMap["pageNum"] = pageNum
	resultMap["pageSize"] = utils.PageSize20
	resultMap["count"] = count
}

/**
选择排序
*/
func (this *ProductDataController) SelectSort(target []models.DetailedData) {
	length := len(target)
	for i := 0; i < length-1; i++ {
		index := i
		for j := i + 1; j < length; j++ {
			if target[j].CreateDate > target[index].CreateDate {
				index = j
			}
		}
		//如果找最大的数据时候，该位置上的数据不是最大，才需要交换
		if index != i {
			temp := target[index]
			target[index] = target[i]
			target[i] = temp
		}
	}
	//fmt.Println("======", target)
}

/**
用户信息趋势分析

*/
//@router /getUserStatisticsData [get]
func (this *ProductDataController) GetUserStatisticsData() {
	defer this.ServeJSON()
	productId, _ := this.GetInt("productId", 1)
	cooperation_type, _ := models.GetProductCooperationById(productId)
	startDate := this.GetString("startDate")
	endDate := this.GetString("stopDate")
	code, _ := this.GetInt("identifyCode", 1) //1：日 2：周
	state, _ := this.GetInt("state", 1)       //1，2，3，4，5，6，7
	if startDate == "" {
		startDate = "2017-05-01"
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "选择开始日期异常", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "选择开始日期异常;err:" + err.Error()}
		return
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "选择截止日期异常", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "选择截止日期异常;err:" + err.Error()}
		return
	}
	condition := ""
	params := make([]string, 0, 2)
	if cooperation_type == 0 { //API
		if startDate != "" {
			if state == 6 { //统计放款金额
				condition += " AND DATE_FORMAT(bl.real_time,'%Y-%m-%d') >= ? "
				params = append(params, endDate)
			} else {
				condition += " AND DATE_FORMAT(pru.create_time,'%Y-%m-%d') >= ?"
				params = append(params, startDate)
			}
		}
		if endDate != "" {
			if state == 6 { //统计放款金额
				condition += " AND DATE_FORMAT(bl.real_time,'%Y-%m-%d') <= ? "
				params = append(params, endDate)
			} else {
				condition += " AND DATE_FORMAT(pru.create_time,'%Y-%m-%d') <= ? "
				params = append(params, endDate)
			}
		}
		var us []models.RegisterUser
		var cmu []models.CreditMoneyUser
		var err error
		switch state {
		case 1:
			us, err = models.GetRegisterUsersByCondition2(productId, condition, params) //注册用户
		case 2:
			us, err = models.GetAuthsersByCondition(productId, "", nil) //认证用户
		case 3:
			us, err = models.GetCreditExtensionUsersByCondition(productId, condition, params) //授信用户
		case 4:
			us, err = models.GetLoanTotalUsersByCondition(productId, condition, params) //申请借款用户
		case 5:
			us, err = models.GetCreditUsersByCondition(productId, condition, params) //放款用户
		case 6:
			cmu, err = models.GetCreditMoneyUsersByCondition(productId, condition, params) //放款金额
		default:
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取状态失败。"}
			return
		}
		if err != nil {
			if err.Error() == utils.ErrNoRow() {
				this.Data["json"] = map[string]interface{}{"ret": 304, "err": "数据不存在"}
				return
			}
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取注册数据失败", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取注册数据失败;err:" + err.Error()}
			return
		}
		umap := make(map[string]interface{})
		if state == 6 {
			for _, u := range cmu {
				umap[u.CreateDate] = u.Count
			}
		} else {
			for _, u := range us {
				umap[u.CreateDate] = u.Count
			}
		}

		sm := make(map[string]interface{}, 0)
		if code == 1 {
			if state == 6 {
				sm = utils.GetSeriesMonth(startTime, endTime, 1)
			} else {
				sm = utils.GetSeriesMonth(startTime, endTime, 0)
			}
		} else if code == 2 {
			if state == 6 {
				sm = utils.GetSeriesDay(startDate, endDate, 1)
			} else {
				sm = utils.GetSeriesDay(startDate, endDate, 0)
			}
		}
		//这天有数据的就填充
		for k, _ := range sm {
			if _, ok := umap[k]; ok {
				sm[k] = umap[k]
				//fmt.Println(k, umap[k])
			}
		}

		sortedKeys := make([]string, 0)
		for k, _ := range sm {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		var resultUser []models.RegisterUser
		var creditUser []models.CreditMoneyUser
		if state == 6 {
			for _, k := range sortedKeys {
				user := new(models.CreditMoneyUser)
				user.CreateDate = k
				user.Count = sm[k].(float32)
				creditUser = append(creditUser, *user)
			}
		} else {
			for _, k := range sortedKeys {
				user := new(models.RegisterUser)
				user.CreateDate = k
				user.Count = sm[k].(int)
				resultUser = append(resultUser, *user)
			}
		}
		if code == 1 {
			if state == 6 {
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": creditUser}
				return
			} else {
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": resultUser}
				return
			}
		} else if code == 2 {
			if state == 6 {
				realResultUser := this.GetCreditMoneyPerWeek(creditUser)
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": realResultUser}
				return
			} else {
				realResultUser := this.GetXXUserPerWeek(resultUser)
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": realResultUser}
				return
			}
		}
	} else if cooperation_type == 1 { //H5
		defaultAgentId, _ := models.GetFirstProductAgentIdByProductId(productId)
		agentId, _ := this.GetInt("agentId", defaultAgentId)
		if productId != 0 {
			condition += " AND p.id = ?"
			params = append(params, strconv.Itoa(productId))
		}
		if agentId != 0 {
			condition += " AND a.id = ?"
			params = append(params, strconv.Itoa(agentId))
		}
		if startDate != "" {
			condition += " AND DATE_FORMAT(dd.date,'%Y-%m-%d') >= ?"
			params = append(params, startDate)
		}
		if endDate != "" {
			condition += " AND DATE_FORMAT(dd.date,'%Y-%m-%d') <= ? "
			params = append(params, endDate)
		}
		var us []models.RegisterUser
		var cmu []models.CreditMoneyUser
		var err error
		switch state {
		case 1:
			us, err = models.GetRegisterUsersByConditionForH5(condition, params) //注册用户
		case 2:
			us, err = models.GetLoanUsersByConditionForH5(condition, params) //认证用户
		case 3:
			us, err = models.GetCreditExtensionUsersByConditionForH5(condition, params) //授信用户
		case 4:
			us, err = models.GetLoanTotalUsersByConditionForH5(condition, params) //申请借款用户
		case 5:
			us, err = models.GetCreditUsersByConditionForH5(condition, params) //放款用户
		case 6:
			cmu, err = models.GetCreditMoneyUsersByConditionForH5(condition, params) //放款金额
		case 7:
			cmu, err = models.GetProfitH5(condition, params) //h5收益
		default:
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取状态失败。"}
			return
		}
		if err != nil {
			if err.Error() == utils.ErrNoRow() {
				this.Data["json"] = map[string]interface{}{"ret": 304, "err": "数据不存在"}
				return
			}
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取注册数据失败", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取注册数据失败;err:" + err.Error()}
			return
		}
		umap := make(map[string]interface{})
		if state == 6 || state == 7 {
			for _, u := range cmu {
				umap[u.CreateDate] = u.Count
			}
		} else {
			for _, u := range us {
				umap[u.CreateDate] = u.Count
			}
		}
		sm := make(map[string]interface{}, 0)
		if code == 1 {
			if state == 6 || state == 7 {
				sm = utils.GetSeriesMonth(startTime, endTime, 1)
			} else {
				sm = utils.GetSeriesMonth(startTime, endTime, 0)
			}
		} else if code == 2 {
			if state == 6 || state == 7 {
				sm = utils.GetSeriesDay(startDate, endDate, 1)
			} else {
				sm = utils.GetSeriesDay(startDate, endDate, 0)
			}
		}
		//这天有数据的就填充
		for k, _ := range sm {
			if _, ok := umap[k]; ok {
				sm[k] = umap[k]
			}
		}
		sortedKeys := make([]string, 0)
		for k, _ := range sm {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		var resultUser []models.RegisterUser
		var creditUser []models.CreditMoneyUser
		if state == 6 || state == 7 {
			for _, k := range sortedKeys {
				user := new(models.CreditMoneyUser)
				user.CreateDate = k
				user.Count = sm[k].(float32)
				creditUser = append(creditUser, *user)
			}
		} else {
			for _, k := range sortedKeys {
				user := new(models.RegisterUser)
				user.CreateDate = k
				user.Count = sm[k].(int)
				resultUser = append(resultUser, *user)
			}
		}
		if code == 1 {
			if state == 6 || state == 7 {
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": creditUser}
				return
			} else {
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": resultUser}
				return
			}
		} else if code == 2 {
			if state == 6 || state == 7 {
				realResultUser := this.GetCreditMoneyPerWeek(creditUser)
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": realResultUser}
				return
			} else {
				realResultUser := this.GetXXUserPerWeek(resultUser)
				this.Data["json"] = map[string]interface{}{"ret": 200, "data": realResultUser}
				return
			}
		}
	}
}

/**
按每周统计XX用户
*/
func (this *ProductDataController) GetXXUserPerWeek(users []models.RegisterUser) (realResultUser []models.RegisterUser) {
	length := len(users)
	realResultUser = make([]models.RegisterUser, 0)
	count := 0
	for i := 1; i <= length; i++ {
		count += users[i-1].Count
		if i%7 == 0 {
			user := new(models.RegisterUser)
			user.CreateDate = users[i-7].CreateDate
			user.Count = count
			realResultUser = append(realResultUser, *user)
			count = 0
		}
	}
	//fmt.Println("GetXXUserPerWeek:", realResultUser)
	return
}

/**
按每周统计放款金额
*/
func (this *ProductDataController) GetCreditMoneyPerWeek(users []models.CreditMoneyUser) (realResultUser []models.CreditMoneyUser) {
	length := len(users)
	realResultUser = make([]models.CreditMoneyUser, 0)
	var totalMoney float32 = 0.0
	for i := 1; i <= length; i++ {
		totalMoney += users[i-1].Count
		if i%7 == 0 {
			user := new(models.CreditMoneyUser)
			user.CreateDate = users[i-7].CreateDate
			user.Count = totalMoney
			realResultUser = append(realResultUser, *user)
			totalMoney = 0.0
		}
	}
	//fmt.Println("GetCreditMoneyPerWeek:", realResultUser)
	return
}

//从excel导入
//@router  /datafromexcel [post]
func (this *ProductDataController) DataFromExcel() {
	resultMap := make(map[string]interface{})
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	resultMap["ret"] = 403
	xlsxUrl := this.GetString("filesPath")
	xlsxUrl = strings.Replace(xlsxUrl, "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,", "", -1)
	xlsxUrl = strings.Replace(xlsxUrl, "data:;base64,", "", -1)
	arr := strings.Split(xlsxUrl, "||")
	length := len(arr[0])
	xlFileUrl := ""
	deleteFile := []string{}
	if length > 0 {
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		fileName := "xlData" + timeStr + ".xlsx"
		filePath := "./static/" + fileName
		err := file.SaveBase64ToFile(arr[0], filePath)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "上传Excel文件失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "上传Excel文件失败"
			return
		}
		deleteFile = append(deleteFile, fileName)
		xlFileUrl = filePath
	}
	xlFile, err := xlsx.OpenFile(xlFileUrl)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "打开Excel文件失败", err.Error(), this.Ctx.Input)
		resultMap["err"] = "打开Excel文件失败"
		return
	}
	//遍历Excel文件并将数据存入到结构体中
	var dailyData models.ExcelData
	var dailyDatas []models.ExcelData
	defer func() {
		if len(deleteFile) > 0 {
			for i := 0; i < len(deleteFile); i++ {
				os.Remove("./static/" + deleteFile[i])
			}
		}
	}()
	sheet := xlFile.Sheets[0]
	for key, row := range sheet.Rows {
		if key == 0 {
			var strs []string
			for _, v := range row.Cells {
				strs = append(strs, v.Value)
			}
			str := strings.Join(strs, ",")
			if str != utils.HEARDER {
				resultMap["err"] = "请按照模板导入数据"
				return
			}
		}
		if key > 0 {
			for k, cell := range row.Cells {
				if k == 0 && cell.Value == "" {
					break
				}
				switch k {
				case 0:
					if strings.Index(cell.Value, "-") == -1 {
						dailyData.Date, err = cell.GetTime(false)
						if err != nil {
							cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取日期数据异常！", err.Error(), this.Ctx.Input)
							resultMap["err"] = "获取日期数据异常1"
							return
						}
					} else {
						dailyData.Date, err = time.Parse("2006-01-02", cell.Value)
						if err != nil {
							cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取日期数据异常！", err.Error(), this.Ctx.Input)
							resultMap["err"] = "获取日期数据异常2"
							return
						}
					}
				case 1:
					if cell.Value == "" {
						dailyData.RegisterCount = 0
					} else {
						dailyData.RegisterCount, err = strconv.Atoi(cell.Value)
						if err != nil {
							resultMap["err"] = "获取注册人数数据异常"
							return
						}
					}

				case 2:
					if cell.Value == "" {
						dailyData.ApplyCount = 0
					} else {
						dailyData.ApplyCount, err = strconv.Atoi(cell.Value)
						if err != nil {
							resultMap["err"] = "获取完成认证人数数据异常"
							return
						}
					}

				case 3:
					if cell.Value == "" {
						dailyData.CreditCount = 0
					} else {
						dailyData.CreditCount, err = strconv.Atoi(cell.Value)
						if err != nil {
							resultMap["err"] = "获取授信人数数据异常"
							return
						}
					}

				case 4:
					if cell.Value == "" {
						dailyData.ApplyLoanCount = 0
					} else {
						dailyData.ApplyLoanCount, err = strconv.Atoi(cell.Value)
						if err != nil {
							resultMap["err"] = "获取申请借款人数据异常"
							return
						}
					}

				case 5:
					if cell.Value == "" {
						dailyData.MakeLoanCount = 0
					} else {
						dailyData.MakeLoanCount, err = strconv.Atoi(cell.Value)
						if err != nil {
							resultMap["err"] = "获取放款人数据异常"
							return
						}
					}

				case 6:
					if cell.Value == "" {
						dailyData.MakeLoanAmount = 0
					} else {
						dailyData.MakeLoanAmount, err = strconv.ParseFloat(cell.Value, 64)
						if err != nil {
							resultMap["err"] = "获取放款金额数据异常"
							return
						}
					}

				case 7:
					jointMode := cell.Value
					if jointMode == "CPA" {
						dailyData.JointMode = 1
					} else if jointMode == "CPS" {
						dailyData.JointMode = 2
					} else if jointMode == "CPA+CPS" {
						dailyData.JointMode = 3
					} else {
						cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取合作模式异常！", "", this.Ctx.Input)
						resultMap["err"] = "获取合作模式异常，其值可以为-> CPA、CPS、CPA+CPS"
						return
					}
				case 8:
					if cell.Value != "注册" && cell.Value != "认证" && cell.Value != "授信" && cell.Value != "申请放款" && cell.Value != "放款" {
						cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取有效性定义异常！", "", this.Ctx.Input)
						resultMap["err"] = "获取有效性定义异常，其值可以为-> 注册、认证、授信、申请放款、放款"
						return
					}
					dailyData.CpaDefine = cell.Value
				case 9:
					if cell.Value == "" {
						dailyData.CpaPrice = 0
					} else {
						dailyData.CpaPrice, err = strconv.ParseFloat(cell.Value, 64)
						if err != nil {
							resultMap["err"] = "获取CPA价格数据异常"
							return
						}
					}

				case 10:
					if cell.Value == "" {
						dailyData.CpsFirstPer = 0
					} else {
						dailyData.CpsFirstPer, err = strconv.ParseFloat(cell.Value, 64)
						if err != nil {
							resultMap["err"] = "获取CPS的首借百分比异常"
							return
						}
					}

				case 11:
					if cell.Value == "" {
						dailyData.CpsAgainPer = 0
					} else {
						dailyData.CpsAgainPer, err = strconv.ParseFloat(cell.Value, 64)
						if err != nil {
							resultMap["err"] = "获取CPS的复借百分比异常"
							return
						}
					}

				}
			}
			dailyDatas = append(dailyDatas, dailyData)
		}
	}
	agentProductId, _ := this.GetInt("AgentProductId")
	//根据agentProductId和日期判断数据是否已存在，存在就更新，不存在就插入
	date, err := models.FindDataFromExcel(agentProductId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "判断数据是否存在失败", err.Error(), this.Ctx.Input)
		resultMap["err"] = "判断数据是否存在失败"
		return
	}
	var xlDatas []models.ExcelData
	var xlData []models.ExcelData
	if len(date) != 0 {
		xlDatas, xlData = EqualQueryData(date, dailyDatas)
	} else {
		err = models.InsertDataFromExcel(dailyDatas, agentProductId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存数据失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "保存数据失败"
			return
		}
	}
	if len(xlDatas) != 0 {
		err = models.UpdateDataFromExcel(xlDatas, agentProductId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "更新数据失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "更新数据失败"
			return
		}
	}
	if len(xlData) != 0 {
		err = models.InsertDataFromExcel(xlData, agentProductId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存数据失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "保存数据失败"
			return
		}
	}

	resultMap["ret"] = 200
	resultMap["msg"] = "记录产品运营信息成功！"
}

//对比数据是否存在
func EqualQueryData(date []time.Time, dailyDatas []models.ExcelData) (xlDatas []models.ExcelData, xlData []models.ExcelData) {
	reMap := make(map[string]time.Time)
	for _, v := range date {
		reMap[v.Format("2006-01-02")] = v
	}
	for k, v := range dailyDatas {
		s := v.Date.Format("2006-01-02")
		if _, ok := reMap[s]; ok {
			xlDatas = append(xlDatas, dailyDatas[k])
		} else {
			xlData = append(xlData, dailyDatas[k])
		}
	}
	return xlDatas, xlData
}

//下载Excel模板
//@router /downloadexcel [get]
func (this *ProductDataController) DownLoadExcel() {
	filename := "./static/xlFile.xlsx"
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename=Excel模板.xlsx")
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
}
