package controllers

import (
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type ChannelController struct {
	BaseController
}

//微融外链
func (c *ChannelController) GetAgentWrWlDataList() {
	c.IsNeedTemplate()
	c.TplName = "channel-management/weirong_wl_data.html"
}

//外链数据明细
func (c *ChannelController) WlDataStructure() {
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
	//获取渠道
	source := c.GetString("source")
	if source == "所有渠道" {
		source = ""
	}
	//拼接条件
	if source != "" {
		condition += " AND u.out_put_source=?"
		params = append(params, source)
	}
	var count, pageCount int
	//按照日期分组查询数据明细
	var dailyDataList []models.DailyData
	dailyList, err := models.GetChannelHistoryData(params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取历史数据明细异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取历史数据明细异常！"
		return
	}
	for k, v := range dailyList {
		if v.RegisterCount != 0 {
			dailyList[k].UserPerCount = utils.SubFloatToFloat(float64(v.ProRegisterCount)/float64(v.RegisterCount), 2)
			dailyList[k].UserPerProfit = utils.SubFloatToFloat(v.TotalProfit/float64(v.RegisterCount), 2)
		}
	}
	//获取今天的统计信息
	//今天微融注册用户
	var dailyData models.DailyData
	todayRegisterCount, err := models.GetWrRegisterCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融注册用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融注册用户异常！"
		return
	}
	todayWrFirstUserCount1, err := models.GetWrFirstUserCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融登录用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融登录用户异常！"
		return
	}
	todayWrORCCount1, err := models.GetWrORCCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融ocr用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融ocr用户异常！"
		return
	}
	todayApplyNowUserCount1, err := models.GetApplyNowUserCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融点击立即申请用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天微融点击立即申请用户异常！"
		return
	}
	todayWrIdentifyCountToday, err := models.GetWrIdentifyCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天完成认证用户异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "今天完成认证用户异常！"
		return
	}
	//获取今日导流量
	todayProRegesiterCount, err := models.GetProRegitsterCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取获取今日导流量异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取获取今日导流量异常！"
		return
	}
	//总收益
	todayProfit, err := models.GetTodayProfit(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取获取今日总收益异常！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取获取今日总收益异常！"
		return
	}
	if todayRegisterCount != 0 {
		//获取人均注册平台
		dailyData.UserPerCount = float64(todayProRegesiterCount) / float64(todayRegisterCount)
		dailyData.UserPerCount = utils.SubFloatToFloat(dailyData.UserPerCount, 2)
		//获取人均收益
		dailyData.UserPerProfit = todayProfit / float64(todayRegisterCount)
		dailyData.UserPerProfit = utils.SubFloatToFloat(dailyData.UserPerProfit, 2)
	}
	dailyData.TotalProfit = todayProfit
	dailyData.ProRegisterCount = todayProRegesiterCount //导流量
	dailyData.Date = time.Now()
	dailyData.RegisterCount = todayRegisterCount
	dailyData.LoginCount = todayWrFirstUserCount1
	dailyData.OcrCount = todayWrORCCount1
	dailyData.ApplynowCount = todayApplyNowUserCount1
	dailyData.IdentifyCount = todayWrIdentifyCountToday
	dailyDataList = append(dailyDataList, dailyData)
	dailyDataList = append(dailyDataList, dailyList[:]...)
	count = len(dailyDataList)
	pageCount, err = utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取需要页数失败", "获取需要页数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取需要页数失败"
		return
	}
	if count != 0 {
		if pageNum == pageCount {
			dailyDataList = dailyDataList[(pageNum-1)*utils.PAGE_SIZE20:]
		} else {
			dailyDataList = dailyDataList[(pageNum-1)*utils.PAGE_SIZE20 : pageNum*utils.PAGE_SIZE20]
		}
	}
	resultMap["dailyDataList"] = dailyDataList //以日期为单位，注册，登录，认证，点击立即申请人数等数据
	resultMap["pageNum"] = pageNum
	resultMap["pageCount"] = pageCount
	resultMap["count"] = count
	resultMap["source"] = source
	resultMap["pageSize"] = utils.PageSize20
	resultMap["ret"] = 200
	return
}

//模糊查询获取渠道列表
func (c *ChannelController) GetAgentSourceList() {
	defer c.ServeJSON()
	name := c.GetString("source")
	if name == "所有渠道" {
		name = ""
	}
	condition := ""
	var param []string
	if name != "" {
		condition = " AND a.out_put_source LIKE ? "
		param = append(param, "%"+name+"%")
	} else {
		condition = "AND a.out_put_source != ''"
	}
	//模糊查询渠道列表
	outPutSourceList, err := models.GetOutPutSourceListByName(condition, param)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": "获取渠道列表失败!"}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取渠道列表异常", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "sourceList": outPutSourceList}
	return
}

//外链数据累计
func (c *ChannelController) AgentWlTotalData() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	condition := ""
	params := []string{}
	//获取渠道
	source := c.GetString("source")
	if source == "所有渠道" {
		source = ""
	}
	flag := true //用于判断是否有该渠道
	//拼接条件
	if source != "" {
		count, err := models.GetOutPutSourceByName(source)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取渠道异常!", err.Error(), c.Ctx.Input)
		}
		if count == 0 {
			flag = false
		}
		condition += " AND u.out_put_source=?"
		params = append(params, source)
	}
	var registerCount, firstuserCount, ORCCount, ApplyNowCount int
	var err error
	//按照日期分组查询数据明细
	weirongDataCount, err := models.GetTotalHistoryData(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取累计历史数据异常！", err.Error(), c.Ctx.Input)
	}
	//今天微融注册用户
	todayRegisterCount, err := models.GetWrRegisterCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融注册用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrFirstUserCount, err := models.GetWrFirstUserCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融登录用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrORCCount, err := models.GetWrORCCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融ocr用户异常！", err.Error(), c.Ctx.Input)
	}
	todayApplyNowUserCount, err := models.GetApplyNowUserCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融点击立即申请用户异常！", err.Error(), c.Ctx.Input)
	}
	registerCount = weirongDataCount.RegisterCount + todayRegisterCount
	firstuserCount = weirongDataCount.FirstCount + todayWrFirstUserCount
	ORCCount = weirongDataCount.OcrCount + todayWrORCCount
	ApplyNowCount = weirongDataCount.ApplynowCount + todayApplyNowUserCount

	resultMap["registerCount"] = registerCount   //累计注册用户
	resultMap["firstuserCount"] = firstuserCount //累计激活用户
	resultMap["ORCCount"] = ORCCount             //累计完成orc用户
	resultMap["ApplyNowCount"] = ApplyNowCount   //累计点击立即申请用户
	resultMap["flag"] = flag                     //用于判断是否有该渠道
	resultMap["ret"] = 200
	return
}

/**
用户信息趋势分析
*/
//@router /getWrStatisticsData [get]
func (this *ChannelController) GetWrStatisticsData() {
	defer this.ServeJSON()
	startDate := this.GetString("startDate")
	endDate := this.GetString("stopDate")
	code, _ := this.GetInt("identifyCode") //1：日 2：周
	state, _ := this.GetInt("state")       //1，2，3，4，5，6
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

	//获取渠道
	source := this.GetString("source")
	if source == "所有渠道" {
		source = ""
	}
	//拼接条件
	if source != "" {
		condition += " AND u.out_put_source=?"
		params = append(params, source)
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
		us, err = models.GetWrRegisterUsersByCondition(condition, params) //注册用户
	case 2:
		if startDate != "" {
			condition += " AND u.active_time >= ?"
			params = append(params, startDate)
		} else {
			startDate = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		}
		if endDate != "" {
			condition += " AND DATE_ADD(u.active_time, INTERVAL -1 DAY) <= ? "
			params = append(params, endDate)

		} else {
			endDate = time.Now().Format("2006-01-02")
		}
		us, err = models.GetWrFirstUsersByCondition(condition, params) //激活用户
	case 3:
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
		us, err = models.GetWrOcrUsersByCondition(condition, params) //OCR用户
	case 4:
		if startDate != "" {
			condition += " AND GREATEST(a.user_data_time,a.zm_auth_time) >= ?"
			params = append(params, startDate)
		} else {
			startDate = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		}
		if endDate != "" {
			condition += " AND DATE_ADD(GREATEST(a.user_data_time,a.zm_auth_time), INTERVAL -1 DAY) <= ? "
			params = append(params, endDate)

		} else {
			endDate = time.Now().Format("2006-01-02")
		}
		us, err = models.GetWrIdentifyUsersByCondition(condition, params) //完成认证用户
	case 5:
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
		us, err = models.GetApplyNowUsersByCondition(condition, params) //点击立即申请用户
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

/**
按每周统计XX用户
*/
func (this *ChannelController) GetWrUserPerWeek(users []models.WrRegisterUser) (realResultUser []models.WrRegisterUser) {
	length := len(users)
	realResultUser = make([]models.WrRegisterUser, 0)
	count := 0
	for i := 1; i <= length; i++ {
		count += users[i-1].Count
		if i%7 == 0 {
			user := new(models.WrRegisterUser)
			user.CreateDate = users[i-7].CreateDate
			user.Count = count
			realResultUser = append(realResultUser, *user)
			count = 0
		}
	}
	//fmt.Println("GetXXUserPerWeek:", realResultUser)
	return
}

//数据导出到EXCEL
//@router /datatoexcel [get]
func (c *ChannelController) DataToExcel() {
	condition := ""
	params := []string{}
	//获取渠道
	source := c.GetString("source")
	if source == "所有渠道" {
		source = ""
	}
	var exl [][]string
	exl = [][]string{{"日期", "注册用户", "激活用户", "完成OCR用户", "认证完成用户", "点击立即申请用户数", "第三方导流量", "人均注册平台", "人均创造收益", "总收益"}}
	//拼接条件
	if source != "" {
		condition += " AND u.out_put_source=?"
		params = append(params, source)
	}
	var exportToExcel []models.DailyData
	dailyList, err := models.GetChannelHistoryData(params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取历史数据明细异常！", err.Error(), c.Ctx.Input)
	}
	for k, v := range dailyList {
		if v.RegisterCount != 0 {
			dailyList[k].UserPerCount = utils.SubFloatToFloat(float64(v.ProRegisterCount)/float64(v.RegisterCount), 2)
			dailyList[k].UserPerProfit = utils.SubFloatToFloat(v.TotalProfit/float64(v.RegisterCount), 2)
		}
	}
	//获取今天的统计信息
	//今天微融注册用户
	var dailyData models.DailyData
	todayRegisterCount, err := models.GetWrRegisterCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融注册用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrFirstUserCount1, err := models.GetWrFirstUserCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融登录用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrORCCount1, err := models.GetWrORCCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融ocr用户异常！", err.Error(), c.Ctx.Input)
	}
	todayApplyNowUserCount1, err := models.GetApplyNowUserCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天微融点击立即申请用户异常！", err.Error(), c.Ctx.Input)
	}
	todayWrIdentifyCountToday, err := models.GetWrIdentifyCountToday(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "今天完成认证用户异常！", err.Error(), c.Ctx.Input)
	}
	//获取今日导流量
	todayProRegesiterCount, err := models.GetProRegitsterCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取获取今日导流量异常！", err.Error(), c.Ctx.Input)
	}
	//总收益
	todayProfit, err := models.GetTodayProfit(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取获取今日总收益异常！", err.Error(), c.Ctx.Input)
	}
	if todayRegisterCount != 0 {
		//获取人均注册平台
		dailyData.UserPerCount = float64(todayProRegesiterCount) / float64(todayRegisterCount)
		dailyData.UserPerCount = utils.SubFloatToFloat(dailyData.UserPerCount, 2)
		//获取人均收益
		dailyData.UserPerProfit = todayProfit / float64(todayRegisterCount)
		dailyData.UserPerProfit = utils.SubFloatToFloat(dailyData.UserPerProfit, 2)
	}
	dailyData.TotalProfit = todayProfit
	dailyData.ProRegisterCount = todayProRegesiterCount //导流量
	dailyData.Date = time.Now()
	dailyData.RegisterCount = todayRegisterCount
	dailyData.LoginCount = todayWrFirstUserCount1
	dailyData.OcrCount = todayWrORCCount1
	dailyData.ApplynowCount = todayApplyNowUserCount1
	dailyData.IdentifyCount = todayWrIdentifyCountToday
	exportToExcel = append(exportToExcel, dailyData)
	exportToExcel = append(exportToExcel, dailyList[:]...)
	for _, v := range exportToExcel {
		exp := []string{v.Date.Format("2006-01-02"), strconv.Itoa(v.RegisterCount), strconv.Itoa(v.LoginCount), strconv.Itoa(v.OcrCount), strconv.Itoa(v.IdentifyCount), strconv.Itoa(v.ApplynowCount), strconv.Itoa(v.ProRegisterCount), strconv.FormatFloat(v.UserPerCount, 'f', -1, 64), strconv.FormatFloat(v.UserPerProfit, 'f', -1, 64), strconv.FormatFloat(v.TotalProfit, 'f', -1, 64)}
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

func (c *ChannelController) WlHistogramData() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	histogramData := [4]models.HistogramData{}
	//var histogramData [4]models.HistogramData
	condition := ""
	params := []string{}
	//获取渠道
	source := c.GetString("source")
	if source == "所有渠道" {
		source = ""
	}
	//拼接条件
	if source != "" {
		condition += " AND u.out_put_source=?"
		params = append(params, source)
	}
	registerCount, err := models.GetWrRegisterCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "累计微融注册用户异常！", err.Error(), c.Ctx.Input)
	}
	//累计激活用户
	firstuserCount, err := models.GetWrFirstUserCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "累计激活用户异常！", err.Error(), c.Ctx.Input)
	}
	// 累计完成orc用户
	ORCCount, err := models.GetWrORCCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "累计完成orc用户异常！", err.Error(), c.Ctx.Input)
	}
	//完成认证用户
	identifycount, err := models.GetWrIdentifyCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取累完成认证用户异常！", err.Error(), c.Ctx.Input)
	}
	// 累计点击立即申请用户
	ApplyNowCount, err := models.GetApplyNowUserCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "累计点击立即申请用户异常", err.Error(), c.Ctx.Input)
	}
	histogramData[0].Name = "注册到激活"
	if registerCount != 0 {
		histogramData[0].Count = utils.SubFloatToFloat(100*float64(firstuserCount)/float64(registerCount), 2)
	}
	histogramData[1].Name = "激活到OCR"
	if firstuserCount != 0 {
		histogramData[1].Count = utils.SubFloatToFloat(100*float64(ORCCount)/float64(firstuserCount), 2)
	}
	histogramData[2].Name = "OCR到认证"
	if ORCCount != 0 {
		histogramData[2].Count = utils.SubFloatToFloat(100*float64(identifycount)/float64(ORCCount), 2)
	}
	histogramData[3].Name = "认证到申请"
	if identifycount != 0 {
		histogramData[3].Count = utils.SubFloatToFloat(100*float64(ApplyNowCount)/float64(identifycount), 2)
	}
	resultMap["histogramData"] = histogramData
	resultMap["ret"] = 200
	return
}
