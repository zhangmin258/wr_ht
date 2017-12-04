package controllers

import (
	//"encoding/json"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

//微融市场数据
func (c *ChannelController) GetAgentWrScDataList() {
	c.IsNeedTemplate()
	c.TplName = "channel-management/weirong_sc_data.html"
}

func (c *ChannelController) GetWrScDataByCondition() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	pkgParams := []int{}
	//获取渠道
	packageName, _ := c.GetInt("pkgName")
	source := c.GetString("source")
	//拼接条件
	if source != "" && source != "所有安卓市场" {
		condition += " AND u.source=? "
		params = append(params, source)
	} else if source == "所有安卓市场" {
		condition += " AND u.source!='AppStore' "
	} else {
		source = "所有市场"
	}
	if packageName != 99 {
		condition += " AND u.register_source=?"
		pkgParams = append(pkgParams, packageName)
	}
	var count, pageCount int
	dailyDataList, err := models.GetWrScHistoryData(packageName, source)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取每日微融市场数据异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取每日微融市场数据异常！"
		return
	}
	//获取今天的统计信息
	//今天微融注册用户
	//var dailyDataList []models.DailyData
	var dailyData models.DailyData
	todayRegisterCount, err := models.GetWrSCRegisterCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取今天微融注册用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取今天微融注册用户异常！"
		return
	}
	todayWrORCCount, err := models.GetWrSCORCCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融ocr用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融ocr用户异常！"
		return
	}
	todayWrIdentifyCountToday, err := models.GetWrSCIdentifyCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天完成认证用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天完成认证用户异常！"
		return
	}
	todayApplyNowUserCount, err := models.GetSCApplyNowUserCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融点击立即申请用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融点击立即申请用户异常！"
		return
	}
	todayProRegesiterCount, err := models.GetTodayProRegesiterCount(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天第三方导流量异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天第三方导流量异常！"
		return
	}
	// 按包和市场获取总收益
	todayPackageMoneyCount, err := models.GetTodayPackageMoney(condition, params, pkgParams)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天本包产生的所有收益异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天本包产生的所有收益异常！"
		return
	}
	dailyData.Date = time.Now()
	dailyData.RegisterCount = todayRegisterCount
	dailyData.OcrCount = todayWrORCCount
	dailyData.ApplynowCount = todayApplyNowUserCount
	dailyData.IdentifyCount = todayWrIdentifyCountToday
	dailyData.ProRegisterCount = todayProRegesiterCount
	if todayRegisterCount != 0 {
		dailyData.UserPerCount = float64(todayProRegesiterCount) / float64(todayRegisterCount)
		dailyData.UserPerCount = utils.SubFloatToFloat(dailyData.UserPerCount, 2)
		dailyData.UserPerProfit = todayPackageMoneyCount / float64(todayRegisterCount)
		dailyData.UserPerProfit = utils.SubFloatToFloat(dailyData.UserPerProfit, 2)
	}
	dailyData.TotalProfit = todayPackageMoneyCount
	dailyDataList = append(dailyDataList, dailyData)
	sort.Sort(models.DailyDataSort(dailyDataList))
	count = len(dailyDataList)
	pageCount, err = utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取需要页数失败", "获取需要页数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取需要页数失败"
		return
	}
	if pageNum == pageCount {
		dailyDataList = dailyDataList[(pageNum-1)*utils.PAGE_SIZE20:]
	} else {
		dailyDataList = dailyDataList[(pageNum-1)*utils.PAGE_SIZE20 : pageNum*utils.PAGE_SIZE20]
	}
	//注册到激活
	resultMap["dailyDataList"] = dailyDataList //以日期为单位，注册，登录，认证，点击立即申请人数等数据
	resultMap["pageNum"] = pageNum
	resultMap["pageCount"] = pageCount
	resultMap["count"] = count
	resultMap["source"] = source
	resultMap["pageSize"] = utils.PageSize20
	resultMap["ret"] = 200
	return
}

func (c *ChannelController) GetWrScTotalData() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	condition := ""
	params := []string{}
	pkgParams := []int{}
	//获取渠道
	packageName, _ := c.GetInt("pkgName")
	source := c.GetString("source")
	//拼接条件
	if source != "" && source != "所有安卓市场" {
		condition += " AND u.source=? "
		params = append(params, source)
	} else if source == "所有安卓市场" {
		condition += " AND u.source!='AppStore' "
	}
	if packageName != 99 {
		condition += " AND u.register_source=?"
		pkgParams = append(pkgParams, packageName)
	}
	var registerCount, ORCCount, ApplyNowCount, creditCount int
	cacheCountStr := utils.WEIRONGCOUNTSOURCE
	//历史数据
	var weirongDataCount *models.WeiRongDataAll
	var err error
	if packageName != 99 || source != "" {
		weirongDataCount, err = models.GetWrSCRegisterAllCount(condition, params, pkgParams)
	} else {
		weirongDataCount, err = cache.GetWrSCDataCount(condition, cacheCountStr, params, pkgParams)
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取往日微融注册用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取往日微融注册用户异常！"
		return
	}
	//今天微融注册用户
	todayRegisterCount, err := models.GetWrSCRegisterCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取今天微融注册用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取今天微融注册用户异常！"
		return
	}
	todayWrORCCount, err := models.GetWrSCORCCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融ocr用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融ocr用户异常！"
		return
	}
	todayApplyNowUserCount, err := models.GetSCApplyNowUserCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融点击立即申请用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融点击立即申请用户异常！"
		return
	}
	todayTodayLoanUser, err := models.GetTodaySCLoanUser(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融借款用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融借款用户异常！"
		return
	}
	registerCount = weirongDataCount.RegisterCount + todayRegisterCount
	ORCCount = weirongDataCount.OcrCount + todayWrORCCount
	ApplyNowCount = weirongDataCount.ApplynowCount + todayApplyNowUserCount
	creditCount = weirongDataCount.LoanUser + todayTodayLoanUser
	//注册到激活
	resultMap["registerCount"] = registerCount //累计注册用户
	resultMap["ORCCount"] = ORCCount           //累计完成orc用户
	resultMap["ApplyNowCount"] = ApplyNowCount //累计点击立即申请用户
	resultMap["creditCount"] = creditCount     //累计放款用户
	resultMap["ret"] = 200
	return
}

//查询所有包名
func (c *ChannelController) GetAllPkgName() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	packages := []int{}
	packages = append(packages, 99)
	pkgs, err := models.GetAllPackageName()
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有包名异常", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取所有包名异常"
		return
	}
	for _, v := range pkgs {
		packages = append(packages, v)
	}
	var packageName []models.PackageName
	var pkgName models.PackageName
	for _, v := range packages {
		pkgName.PkgId = v
		if v == 99 {
			pkgName.PkgName = "所有分包"
		} else if v == 1 {
			pkgName.PkgName = "主包"
		} else {
			pkgName.PkgName = "分包" + strconv.Itoa(v)
		}
		packageName = append(packageName, pkgName)
	}
	resultMap["firstPkg"] = packageName[0].PkgName
	resultMap["pkgName"] = packageName
	resultMap["ret"] = 200
	return
}

//根据包名查询市场名
func (c *ChannelController) GetSourceByPkg() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	condition := ""
	params := []int{}
	var source []models.SourceName
	var sour models.SourceName
	pkgName, _ := c.GetInt("pkgName")
	if pkgName != 99 {
		condition = " AND register_source=?"
		params = append(params, pkgName)
	}
	sources, err := models.GetAllSourceByPackage(condition, params)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有包名异常", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取市场名称错误"
		return
	}
	sour.Source = "所有市场"
	sour.Name = ""
	source = append(source, sour)
	sour.Source = "所有安卓市场"
	sour.Name = "所有安卓市场"
	source = append(source, sour)
	for _, v := range sources {
		sour.Source = v
		sour.Name = v
		source = append(source, sour)
	}
	resultMap["source"] = source
	resultMap["ret"] = 200
	return
}

/**
用户信息趋势分析
*/
//@router /getWrSCStatisticsData [get]
func (this *ChannelController) GetWrSCStatisticsData() {
	defer this.ServeJSON()
	startDate := this.GetString("startDate")
	endDate := this.GetString("stopDate")
	code, _ := this.GetInt("identifyCode") //1：日 2：周
	state, _ := this.GetInt("state")       //1，2，3，4
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
	//fmt.Println(starTime, endTime, endTime.Sub(starTime).Hours()/24+1)
	condition := ""
	params := make([]string, 0, 3)
	pkgParams := []int{}

	//获取渠道
	packageName, _ := this.GetInt("pkgName")
	source := this.GetString("source")
	//拼接条件
	if source != "" && source != "所有安卓市场" {
		condition += " AND u.source=?"
		params = append(params, source)
	} else if source == "所有安卓市场" {
		condition += " AND u.source!='AppStore' "
	}
	var us []models.WrRegisterUser
	//var err error
	switch state {
	case 1:
		if startDate != "" {
			condition += " AND u.create_date >= ?"
			params = append(params, startDate)
		} else {
			startDate = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		}
		if endDate != "" {
			condition += " AND u.create_date<= ? "
			params = append(params, endDate)
		} else {
			endDate = time.Now().Format("2006-01-02")
		}
		if packageName != 99 {
			condition += " AND u.register_source=?"
			pkgParams = append(pkgParams, packageName)
		}
		us, err = models.GetWrSCRegisterUsersByCondition(condition, params, pkgParams) //注册用户
	case 2:
		if startDate != "" {
			condition += " AND a.real_name_time >= ?"
			params = append(params, startDate)
		} else {
			startDate = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		}
		if endDate != "" {
			condition += " AND DATE_ADD(a.real_name_time, INTERVAL -1 DAY) <= ? "
			params = append(params, endDate)

		} else {
			endDate = time.Now().Format("2006-01-02")
		}
		if packageName != 99 {
			condition += " AND u.register_source=?"
			pkgParams = append(pkgParams, packageName)
		}
		us, err = models.GetWrSCOcrUsersByCondition(condition, params, pkgParams) //OCR用户
	case 3:
		if startDate != "" {
			condition += " AND GREATEST(ua.user_data_time,ua.zm_auth_time) >= ?"
			params = append(params, startDate)
		} else {
			startDate = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		}
		if endDate != "" {
			condition += " AND DATE_ADD(GREATEST(ua.user_data_time,ua.zm_auth_time), INTERVAL -1 DAY) <= ? "
			params = append(params, endDate)

		} else {
			endDate = time.Now().Format("2006-01-02")
		}
		if packageName != 99 {
			condition += " AND u.register_source=?"
			pkgParams = append(pkgParams, packageName)
		}
		us, err = models.GetWrSCIdentifyUsersByCondition(condition, params, pkgParams) //完成认证用户
	case 4:
		if startDate != "" {
			condition += " AND t.first_loan_date >= ?"
			params = append(params, startDate)
		} else {
			startDate = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		}

		if endDate != "" {
			condition += " AND DATE_ADD(t.first_loan_date, INTERVAL -1 DAY) <= ? "
			params = append(params, endDate)

		} else {
			endDate = time.Now().Format("2006-01-02")
		}
		if packageName != 99 {
			condition += " AND u.register_source=?"
			pkgParams = append(pkgParams, packageName)
		}
		us, err = models.GetSCApplyNowUsersByCondition(condition, params, pkgParams) //点击立即申请用户
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

	for _, u := range us {
		umap[u.CreateDate] = u.Count
	}

	sm := make(map[string]interface{}, 0)
	if code == 1 {
		sm = utils.GetSeriesMonth(startTime, endTime, 0)
	} else if code == 2 {
		sm = utils.GetSeriesDay(startDate, endDate, 0)
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
	var resultUser []models.WrRegisterUser
	//var creditUser []models.CreditMoneyUser

	for _, k := range sortedKeys {
		user := new(models.WrRegisterUser)
		user.CreateDate = k
		user.Count = sm[k].(int)
		resultUser = append(resultUser, *user)
	}

	if code == 1 {
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": resultUser}
		return
	} else if code == 2 {
		realResultUser := this.GetWrUserPerWeek(resultUser)
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": realResultUser}
		return
	}
}

//数据导出到EXCEL
//@router /scdatatoexcel [get]
func (c *ChannelController) SCDataToExcel() {
	condition := ""
	params := []string{}
	pkgParams := []int{}
	//获取渠道
	packageName, _ := c.GetInt("pkgName")
	source := c.GetString("source")
	var exl [][]string
	exl = [][]string{{"日期", "注册用户", "完成OCR用户", "认证完成用户", "点击立即申请用户数", "第三方导流量", "人均注册平台", "人均创造收益", "总收益"}}
	//获取所有的渠道列表
	//拼接条件
	condition += "  AND u.source !='' "
	if source != "" && source != "所有安卓市场" {
		condition += " AND u.source=?"
		params = append(params, source)
	} else if source == "所有安卓市场" {
		condition += " AND u.source!='AppStore' "
	} else {
		source = "所有市场"
	}
	if packageName != 99 {
		condition += " AND u.register_source=?"
		pkgParams = append(pkgParams, packageName)
	}
	var exportToExcel []models.DailyData
	var err error
	dailyList, err := models.GetWrScHistoryData(packageName, source)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "累计微融注册用户异常！", err.Error(), c.Ctx.Input)
	}
	//获取今天的统计信息
	//今天微融注册用户
	var dailyData models.DailyData
	todayRegisterCount, err1 := models.GetWrSCRegisterCountToday(condition, params, pkgParams)
	if err1 != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融注册用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrORCCount1, err := models.GetWrSCORCCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融ocr用户异常！", err.Error(), c.Ctx.Input)
	}
	todayApplyNowUserCount1, err := models.GetSCApplyNowUserCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融点击立即申请用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrIdentifyCountToday, err := models.GetWrSCIdentifyCountToday(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天完成认证用户异常！", err.Error(), c.Ctx.Input)
	}
	todayProRegesiterCount, err := models.GetTodayProRegesiterCount(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天第三方导流量异常！", err.Error(), c.Ctx.Input)
	}
	// 按包和市场获取总收益
	todayPackageMoneyCount, err := models.GetTodayPackageMoney(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天本包产生的所有收益异常！", err.Error(), c.Ctx.Input)
	}
	dailyData.Date = time.Now()
	dailyData.RegisterCount = todayRegisterCount
	dailyData.OcrCount = todayWrORCCount1
	dailyData.ApplynowCount = todayApplyNowUserCount1
	dailyData.IdentifyCount = todayWrIdentifyCountToday
	dailyData.ProRegisterCount = todayProRegesiterCount
	//添加总收益
	dailyData.TotalProfit = todayPackageMoneyCount
	if todayRegisterCount != 0 {
		dailyData.UserPerCount = float64(todayProRegesiterCount) / float64(todayRegisterCount)
		dailyData.UserPerCount = utils.SubFloatToFloat(dailyData.UserPerCount, 2)
		dailyData.UserPerProfit = todayPackageMoneyCount / float64(todayRegisterCount)
		dailyData.UserPerProfit = utils.SubFloatToFloat(dailyData.UserPerProfit, 2)
	}
	exportToExcel = append(exportToExcel, dailyData)
	exportToExcel = append(exportToExcel, dailyList[:]...)
	for _, v := range exportToExcel {
		exp := []string{v.Date.Format("2006-01-02"), strconv.Itoa(v.RegisterCount), strconv.Itoa(v.OcrCount), strconv.Itoa(v.IdentifyCount), strconv.Itoa(v.ApplynowCount), strconv.Itoa(v.ProRegisterCount), strconv.FormatFloat(v.UserPerCount, 'f', -1, 64), strconv.FormatFloat(v.UserPerProfit, 'f', -1, 64), strconv.FormatFloat(v.TotalProfit, 'f', -1, 64)}
		exl = append(exl, exp)
	}

	filename, err := utils.ExportToExcel(exl)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存文件错误", err.Error(), c.Ctx.Input)
	}
	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除文件错误", err.Error(), c.Ctx.Input)
	}
}

//柱状图
func (c *ChannelController) HistogramData() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var histogramData []models.HistogramData
	var histogram models.HistogramData
	condition := ""
	params := []string{}
	pkgParams := []int{}
	//获取渠道
	packageName, _ := c.GetInt("pkgName")
	source := c.GetString("source")
	//拼接条件
	condition += " AND u.source!='' "
	if source != "" && source != "所有安卓市场" {
		condition += " AND u.source=?"
		params = append(params, source)
	} else if source == "所有安卓市场" {
		condition += " AND u.source!='AppStore' "
	}
	if packageName != 99 {
		condition += " AND u.register_source=?"
		pkgParams = append(pkgParams, packageName)
	}
	var registerCount, ocrCount, identifyCount, applyNowUserCount int
	var err error
	registerCount, err = models.GetWrSCRegisterCount(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取累计微融注册用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取累计微融注册用户异常！"
		return
	}
	ocrCount, err = models.GetWrSCORCCount(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取累计OCR用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取累计OCR用户异常！"
		return
	}
	identifyCount, err = models.GetWrSCIdentifyCount(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取累计认证用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取累计认证用户异常！"
		return
	}
	applyNowUserCount, err = models.GetSCApplyNowUserCount(condition, params, pkgParams)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取累计申请用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取累计申请用户异常！"
		return
	}
	histogram.Name = "注册到OCR"
	if registerCount != 0 {
		histogram.Count = utils.SubFloatToFloat(100*float64(ocrCount)/float64(registerCount), 2)
	}
	histogramData = append(histogramData, histogram)
	histogram.Name = "OCR到认证"
	if ocrCount != 0 {
		histogram.Count = utils.SubFloatToFloat(100*float64(identifyCount)/float64(ocrCount), 2)
	}

	histogramData = append(histogramData, histogram)
	histogram.Name = "认证到申请"
	if identifyCount != 0 {
		histogram.Count = utils.SubFloatToFloat(100*float64(applyNowUserCount)/float64(identifyCount), 2)
	}
	histogramData = append(histogramData, histogram)
	resultMap["histogramData"] = histogramData
	resultMap["ret"] = 200
	return
}
