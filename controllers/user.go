package controllers

import (
	"regexp"
	"sort"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

/*
用户相关接口
*/
type UserController struct {
	BaseController
}

//分页条件查询用户信息列表
//@router /getUserList [get]
func (c *UserController) GetUserList() {
	//设置整体加载
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//手机号/账号
	if account := c.GetString("account"); account != "" {
		condition += " and u.account =?"
		params = append(params, account)
	}
	//姓名
	if VerifyRealName := c.GetString("idName"); VerifyRealName != "" {
		condition += " and um.verify_real_name = ?"
		params = append(params, VerifyRealName)
	}
	//身份证
	if IdCard := c.GetString("idNo"); IdCard != "" {
		condition += " and um.id_card  =?"
		params = append(params, IdCard)
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
	userList, err := models.GetUserList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户数据失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetUserCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有用户数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["userList"] = userList
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "user-management/user_data.html"
}

//获取用户数据页面数据
//@router /getUserData [get]
func (c *UserController) GetUserData() {
	//设置整体加载
	c.IsNeedTemplate()
	c.TplName = "user-management/dataSummarize.html"
}

//获取用户详细数据
//@router /getUserDataDetail [get]
func (c *UserController) GetUserDataDetail() {
	c.IsNeedTemplate()
	/*
	   每日用户数据明细
	*/
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}

	count, err := models.QueryUsersDataDetailCount()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取今日数据明细异常", err.Error(), c.Ctx.Input)
	}
	count += 1
	//转换页数
	pageCount, err := utils.GetPageCount(count, utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
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
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取今日数据明细异常", err.Error(), c.Ctx.Input)
	}

	var data []models.UsersDataDetail
	//获取今日数据明细
	if pageNum == 1 {
		var todayDetailData models.UsersDataDetail
		var todayActiveUsersCount int
		todayDetailData, err = models.QueryTodayDailyDatas()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取今日数据明细异常", err.Error(), c.Ctx.Input)
		}
		todayActiveUsersCount, err = models.QueryTodayActiveData()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取今日数据明细异常", err.Error(), c.Ctx.Input)
		}
		todayDetailData.ActiveCount = todayActiveUsersCount
		data = append(data, todayDetailData)
	}
	data = append(data, usersDataDetail[:]...)
	c.Data["pageCount"] = pageCount                         //总共多少页
	c.Data["pageNum"] = pageNum                             //第几页
	c.Data["startIndex"] = (pageNum - 1) * utils.PageSize20 //每页起始下标
	c.Data["endIndex"] = pageNum*utils.PageSize20 - 1       //每页终止下标
	c.Data["detailData"] = data                             //所有的数据
	c.Data["count"] = count                                 //数据总共多少条
	c.TplName = "user-management/user_data_details.html"
}

//查询用户信息并展示到页面
//@router /getUserInfo [get]
func (c *UserController) GetUserInfo() {
	//设置整体加载
	c.IsNeedTemplate()
	c.TplName = "user-management/dataSummarize.html"
}

/**
用户信息趋势分析
*/
//@router /getUserStatisticsData [get]
func (this *UserController) GetUserStatisticsData() {
	defer this.ServeJSON()
	startDate := this.GetString("startDate")
	endDate := this.GetString("endDate")
	code, _ := this.GetInt("identifyCode")
	userType, _ := this.GetInt("userType")
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

	var us []models.RegisterUser
	switch userType {
	case 1:
		us, err = cache.GetBeforeRegisterUsersCache() //获取今日之前注册用户
	case 2:
		us, err = cache.GetBeforeIdentifyUsersCache() //获取今日之前认证用户
	case 3:
		us, err = cache.GetBeforeLoanUsersCache() //获取今日之前申请次数
	case 4:
		us, err = cache.GetBeforeCreditUsersCache() //获取今日之前放款用户
	case 5:
		us, err = cache.GetBeforeActiveUsersCache() //今日之前活跃用户
	default:
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取状态失败。"}
		return
	}
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取今日之前数据失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取今日之前数据失败;err:" + err.Error()}
		return
	}
	var tu models.RegisterUser
	switch userType {
	case 1:
		tu, err = models.GetTodayRegisterUsers() //获取今日注册用户
	case 2:
		tu, err = models.GetTodayIdentifyUsers() //获取今日认证用户
	case 3:
		tu, err = models.GetTodayLoanUsers() //获取今日申请次数
	case 4:
		tu, err = models.GetTodayCreditUsers() //获取今日放款用户
	case 5:
		tu, err = models.GetTodayActiveUsers() //今日活跃用户
	default:
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取状态失败。"}
		return
	}
	if err != nil {
		if err.Error() == utils.ErrNoRow() {
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": "今日数据不存在"}
			return
		}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取今日数据失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "获取今日数据失败;err:" + err.Error()}
		return
	}
	us = append(us, tu)
	umap := make(map[string]int)
	for _, u := range us {
		umap[u.CreateDate] = u.Count
	}
	sm := make(map[string]int, 0)
	if code == 1 {
		sm = utils.GetSeriesMonths(startTime, endTime)
	} else if code == 2 {
		sm = utils.GetSeriesWeek(startDate, endDate)
	}
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
	for _, k := range sortedKeys {
		user := new(models.RegisterUser)
		user.CreateDate = k
		user.Count = sm[k]
		resultUser = append(resultUser, *user)
	}
	if code == 1 {
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": resultUser}
		return
	} else if code == 2 {
		realResultUser := this.GetUserPerWeek(resultUser)
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": realResultUser}
		return
	}
}

func (this *UserController) GetUserPerWeek(users []models.RegisterUser) (realResultUser []models.RegisterUser) {
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
	return
}

//用户数据分析展示到页面_1
//@router /sataAnalysis [get]
func (this *UserController) DataAnalysis() {
	this.IsNeedTemplate()
	this.TplName = "user-management/user_data_analysis.html"
}

//用户资信分析
//@router /userDataAnalysis [get]
func (this *UserController) UserDataAnalysis() {
	defer this.ServeJSON()
	//设置整体加载
	params := this.GetString("params")
	smap := map[string]int{}
	if params == "gender" { //性别
		genderlist, err := models.StaGender()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "性别解析异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "性别获取异常"}
			return
		}
		for _, v := range genderlist {
			smap[v.Sex] = v.Num
		}
		genderInit := models.AgeInit()
		for i := 0; i < len(genderInit); i++ {
			genderInit[i].Num = smap[genderInit[i].Sex]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": genderInit}
		return
	}
	if params == "zm_score" { //芝麻分
		zmlist, err := models.StaZMScore()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取芝麻分异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取芝麻分异常"}
			return
		}
		for _, v := range zmlist {
			smap[v.Score] = v.Num
		}
		scoreInit := models.ZMScoreInit()
		for i := 0; i < len(scoreInit); i++ {
			scoreInit[i].Num = smap[scoreInit[i].Score]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": scoreInit}
		return
	}
	if params == "age" { //年龄
		agelist, err := models.StaAge()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取年龄数据异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取年龄数据异常"}
			return
		}
		for _, v := range agelist {
			smap[v.Data] = v.Num
		}
		ageInit := models.StaAgeInit()
		for i := 0; i < len(ageInit); i++ {
			ageInit[i].Num = smap[ageInit[i].Data]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": ageInit}
		return
	}
	if params == "credit" { //信用评分
		creditlist, err := models.StaCredit()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取信用评分异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取信用评分异常"}
			return
		}
		for _, v := range creditlist {
			smap[v.Credit] = v.Num
		}
		creditInit := models.StaCreditInit()
		for i := 0; i < len(creditInit); i++ {
			creditInit[i].Num = smap[creditInit[i].Credit]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": creditInit}
		return
	}
	if params == "job" { //职业身份
		joblist, err := models.StaJob()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取职业身份数据异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取职业身份数据异常"}
			return
		}
		for _, v := range joblist {
			smap[v.Userjob] = v.Num
		}
		jobInit := models.StaJobInit()
		for i := 0; i < len(jobInit); i++ {
			jobInit[i].Num = smap[jobInit[i].Userjob]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": jobInit}
		return
	}
	if params == "os" { //操作系统
		oslist, err := models.StaOS()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取操作系统异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取操作系统异常"}
			return
		}
		for _, v := range oslist {
			smap[v.Os] = v.Num
		}
		osInit := models.StaOSInit()
		for i := 0; i < len(osInit); i++ {
			osInit[i].Num = smap[osInit[i].Os]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": osInit}
		return
	}
	if params == "operators" { //运营商
		accounts, err := models.GetAccounts()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取运营商数据异常", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取运营商数据异常"}
			return
		}
		reg1 := regexp.MustCompile("^1([3][4-9]|[4][7]|[5][0-27-9]|[8][2-478])[0-9]{8}$")
		reg2 := regexp.MustCompile("^1([3][0-2]|[4][5]|[5][5-6]|[7][6]|[8][5-6])[0-9]{8}$")
		reg3 := regexp.MustCompile("^1(3[3]|5[3]|7[37]|8[019])[0-9]{8}$")
		counts := make([]int, 4)
		for _, account := range accounts {
			switch {
			case reg1.MatchString(account):
				counts[0] += 1
			case reg2.MatchString(account):
				counts[1] += 1
			case reg3.MatchString(account):
				counts[2] += 1
			default:
				counts[3] += 1
			}
		}
		opi := models.UserOperatorsInit()
		for i := 0; i < len(opi); i++ {
			opi[i].Num = counts[i]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": opi}
		return
	}
	if params == "fundsecurity" { //社保公积金
		fundCount, securityCount, securityFundCount, totalCount, err := models.GetSecurityFundUsers()
		if err != nil {
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取社保公积金信息失败"}
			return
		}
		securityfundinit := models.SecurityFundInit()
		securityfundinit[0].Num = securityCount
		securityfundinit[1].Num = fundCount
		securityfundinit[2].Num = securityFundCount
		securityfundinit[3].Num = totalCount
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": securityfundinit}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 403, "err": "请求参数解析错误！"}
	return
}

//用户贷款需求
//@router /businessLoanAnalysis [get]
func (this *UserController) BusinessLoanAnalysis() {
	defer this.ServeJSON()
	params := this.GetString("params")
	var smap = map[string]int{}
	if params == "money" {
		money, err := models.StaLoanMoney()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取贷款金额失败", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取贷款金额失败"}
			return
		}
		for _, v := range money {
			smap[v.MoneyAccount] = v.Num
		}
		lmInit := models.StaLoanMoneyInit()
		for i := 0; i < len(lmInit); i++ {
			lmInit[i].Num = smap[lmInit[i].MoneyAccount]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": lmInit}
		return
	}
	if params == "termcount" {
		termcount, err := models.StaLoanTermCount()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取贷款期限失败", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取贷款期限失败"}
			return
		}
		for _, v := range termcount {
			smap[v.Termcount] = v.Num
		}
		qx := models.BusinessLoanTermCountInit()
		for i := 0; i < len(qx); i++ {
			qx[i].Num = smap[qx[i].Termcount]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": qx}
		return
	}
	if params == "termtimes" {
		termtimes, err := models.StaBusinessLoanTimes()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取借款次数失败", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取借款次数失败"}
			return
		}
		for _, v := range termtimes {
			smap[v.Times] = v.Num
		}
		timesInit := models.BusinessLoanTimesInit()
		for i := 0; i < len(timesInit); i++ {
			timesInit[i].Num = smap[timesInit[i].Times]
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "arr": timesInit}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 403, "err": "请求参数解析错误！"}
}

//用户地域分布
//@router /businessLoanSiteAnalysis [get]
func (this *UserController) BusinessLoanSiteAnalysis() {
	defer this.ServeJSON()
	var list []models.UserAddress
	var err error
	resultMap := map[string]interface{}{}
	params := this.GetString("params")
	smap := map[string]int{}
	resultMap["ret"] = 403
	if params == "registerUser" {
		list, err = models.StaRegisterUserAddress()
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取注册成功用户数据失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "获取注册成功用户数据失败"
			return
		}
	}
	if params == "approvedUser" {
		list, err = models.StaApprovedUserAddress() //approvedUser
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取认证成功用户数据失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "获取认证成功用户数据失败"
			return
		}
	}
	if params == "loanedUser" {
		list, err = models.StaLoanedUserAddress() //loanedUser
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取借款成功用户数据失败", err.Error(), this.Ctx.Input)
			resultMap["err"] = "获取借款成功用户数据失败"
			return
		}
	}
	if len(list) > 0 {
		for _, v := range list {
			smap[v.Provience] = v.Num
		}
		initAddress := models.ProvienceInit()
		for i := 0; i < len(initAddress); i++ {
			initAddress[i].Num = smap[initAddress[i].Provience]
		}
		resultMap["ret"] = 200
		resultMap["arr"] = initAddress
		this.Data["json"] = resultMap
	} else {
		resultMap["msg"] = "暂无数据"
	}
}
