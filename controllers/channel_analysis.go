package controllers

import (
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type ChannelAnalysisController struct {
	BaseController
}

//@router /getchannelanalysis [get]
func (c *ChannelAnalysisController) GetChannelAnalysis() {
	c.IsNeedTemplate()
	name, err := models.GetFristOutPutSourceName()
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取第一个产品异常", err.Error(), c.Ctx.Input)
	}
	c.Data["name"] = name
	c.TplName = "channel-management/channel_data_analysis.html"
}

//@router /getOutChannelProducts [get]
func (this *ChannelAnalysisController) GetOutChannelProducts() {
	defer this.ServeJSON()
	name := this.GetString("name")
	condition := " AND out_put_source like ?"
	params := "%" + name + "%"
	outPutSource, err := models.GetOutPutSourceName(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取外链渠道异常", err.Error(), this.Ctx.Input)
	}
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 200
	resultMap["outPutSource"] = outPutSource
	this.Data["json"] = resultMap
}

//各渠道数据对比
//@router /getchanneldetaildata [get]
func (c *ChannelAnalysisController) GetChannelDetailData() {
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
		condition += " AND u.create_date>=? "
		params = append(params, startDate1)
	}
	if endDate1 := c.GetString("endDate"); endDate1 != "" {
		condition += " AND u.create_date<=? "
		params = append(params, endDate1)
	}

	//渠道列表
	channellist, err := models.GetAnalysisChannel("") //渠道名称列表
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品列表失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取渠道列表失败"
		return
	}
	//1.注册用户，2：激活用户，3：注册到激活转化，4：人均注册平台数
	var AnalysisChannel []models.AnalysisChannel
	switch code {

	case "1": //注册用户数
		h5RegisteredUsersCount, err := models.GetH5RegisteredUsers_Count(condition, params) //注册用户
		if err != nil {
			resultMap["err"] = "获取注册用户数失败"
			return
		}
		h5RUC := make(map[string]int)
		for _, v := range h5RegisteredUsersCount {
			h5RUC[v.Name] = v.Count
		}
		for _, v := range channellist {
			var r models.AnalysisChannel
			r.Name = v.Name
			if _, ok := h5RUC[v.Name]; ok {
				r.RegisteredUsersCount = h5RUC[v.Name]
			}
			AnalysisChannel = append(AnalysisChannel, r)
		}
		// for _, v := range channellist {
		// 	var r models.AnalysisChannel
		// 	r.Name = v.Name
		// 	for _, v1 := range h5RegisteredUsersCount {
		// 		if v.Name == v1.Name {
		// 			r.RegisteredUsersCount = v1.Count
		// 			break
		// 		}
		// 	}
		// 	AnalysisChannel = append(AnalysisChannel, r)
		// }

	case "2": //激活用户数
		h5ActivatedUsersCount, err := models.GetH5ActivatedUsers_Count(condition, params) //激活用户
		if err != nil {
			resultMap["err"] = "获取激活用户数失败"
			return
		}
		h5AUC := make(map[string]int)
		for _, v := range h5ActivatedUsersCount {
			h5AUC[v.Name] = v.Count
		}
		for _, v := range channellist {
			var r models.AnalysisChannel
			r.Name = v.Name
			if _, ok := h5AUC[v.Name]; ok {
				r.ActivatedUsersCount = h5AUC[v.Name]
			}
			AnalysisChannel = append(AnalysisChannel, r)
		}
		// for _, v := range channellist {
		// 	var r models.AnalysisChannel
		// 	r.Name = v.Name
		// 	for _, v1 := range h5ActivatedUsersCount {
		// 		if v.Name == v1.Name {
		// 			r.ActivatedUsersCount = v1.Count
		// 			break
		// 		}
		// 	}
		// 	AnalysisChannel = append(AnalysisChannel, r)
		// }

	case "3": //注册到激活转化率
		h5RegisteredUsersCount, err := models.GetH5RegisteredUsers_Count(condition, params) //注册用户
		if err != nil {
			resultMap["err"] = "获取注册用户数失败"
			return
		}
		h5ActivatedUsersCount, err := models.GetH5ActivatedUsers_Count(condition, params) //激活用户
		if err != nil {
			resultMap["err"] = "获取激活用户数失败"
			return
		}

		h5RUC := make(map[string]int)
		for _, v := range h5RegisteredUsersCount {
			h5RUC[v.Name] = v.Count
		}
		h5AUC := make(map[string]int)
		for _, v := range h5ActivatedUsersCount {
			h5AUC[v.Name] = v.Count
		}
		for _, v := range channellist {
			var r models.AnalysisChannel
			x := 0
			y := 0
			r.Name = v.Name
			if _, ok := h5RUC[v.Name]; ok {
				x = h5RUC[v.Name]
			}
			if _, ok := h5AUC[v.Name]; ok {
				y = h5AUC[v.Name]
			}
			if x == 0 {
				r.Rac = 0
			} else {
				r.Rac = float64(y) / float64(x)
			}
			AnalysisChannel = append(AnalysisChannel, r)
		}

		// for _, v := range channellist {
		// 	var r models.AnalysisChannel
		// 	x := 0
		// 	y := 0
		// 	r.Name = v.Name
		// 	for _, v1 := range h5RegisteredUsersCount {
		// 		if v.Name == v1.Name {
		// 			x = v1.Count
		// 			break
		// 		}
		// 	}
		// 	for _, v1 := range h5ActivatedUsersCount {
		// 		if v.Name == v1.Name {
		// 			y = v1.Count
		// 			break
		// 		}
		// 	}
		// 	if x == 0 {
		// 		r.Rac = 0
		// 	} else {
		// 		r.Rac = float64(y) / float64(x)
		// 	}
		// 	AnalysisChannel = append(AnalysisChannel, r)
		// }
	case "4": //人均注册平台数
		params1 := []string{}
		condition1 := ""
		//开始时间1
		if startDate1 := c.GetString("startDate"); startDate1 != "" {
			condition1 += " AND urp.create_time>=? "
			params1 = append(params1, startDate1)
		}
		if endDate1 := c.GetString("endDate"); endDate1 != "" {
			endDate1 += " 23:59:59"
			condition1 += " AND urp.create_time<=? "
			params1 = append(params1, endDate1)
		}
		h5RegisteredPlatformsCount, err := models.GetH5RegisteredPlatforms_Count(condition1, params1)
		if err != nil {
			resultMap["err"] = "获取用户注册平台数失败"
			return
		}
		h5ChannelUsersCount, err := models.GetH5RegisteredUsers_Count(condition, params)
		if err != nil {
			resultMap["err"] = "获取渠道用户数失败"
			return
		}

		h5RPC := make(map[string]int)
		for _, v := range h5RegisteredPlatformsCount {
			h5RPC[v.Name] = v.Count
		}
		h5CUC := make(map[string]int)
		for _, v := range h5ChannelUsersCount {
			h5CUC[v.Name] = v.Count
		}
		for _, v := range channellist {
			var r models.AnalysisChannel
			x := 0.0
			y := 0
			r.Name = v.Name
			if _, ok := h5RPC[v.Name]; ok {
				x = float64(h5RPC[v.Name])
			}
			if _, ok := h5CUC[v.Name]; ok {
				y = h5CUC[v.Name]
			}
			if y == 0 {
				r.Crpc = 0
			} else {
				r.Crpc = x / float64(y)
			}
			AnalysisChannel = append(AnalysisChannel, r)
		}

		// for _, v := range channellist {
		// 	var r models.AnalysisChannel
		// 	x := 0.0
		// 	y := 0
		// 	r.Name = v.Name
		// 	for _, v1 := range h5RegisteredPlatformsCount {
		// 		if v.Name == v1.Name {
		// 			x = float64(v1.Count)
		// 			break
		// 		}
		// 	}
		// 	for _, v1 := range h5ChannelUsersCount {
		// 		if v.Name == v1.Name {
		// 			y = v1.Count
		// 			break
		// 		}
		// 	}
		// 	if y == 0 {
		// 		r.Crpc = 0 //crpc:capita registration platform count 人均注册平台数
		// 	} else {
		// 		r.Crpc = x / float64(y)
		// 	}
		// 	AnalysisChannel = append(AnalysisChannel, r)
		// }
	}
	resultMap["ret"] = 200
	resultMap["msg"] = AnalysisChannel
}

//渠道数据明细
//@router /getchanneldata [get]
func (c *ChannelAnalysisController) GetChannelData() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	resultMap["ret"] = 403
	code := c.GetString("code")
	params := []string{}
	condition := ""
	//渠道列表
	channellist, err := models.GetAnalysisChannel("") //渠道名称列表
	if err != nil {
		resultMap["err"] = "获取渠道列表失败"
		return
	}
	//渠道
	condition += " AND u.out_put_source=? "
	channel := c.GetString("name")
	flag := 1
	cou, err := models.GetOutPutSourceNameByName(channel)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据name查询渠道异常！", err.Error(), c.Ctx.Input)
	}
	if cou == 0 {
		flag = 0
	}
	if channel == "" {
		channel = channellist[0].Name
	}
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	st, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		resultMap["err"] = "获取时间错误"
		return
	}
	et, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		resultMap["err"] = "获取时间错误"
		return
	}
	//开始时间2
	if startDate != "" {
		condition += " AND u.create_date>=? "
		params = append(params, startDate)
	}
	//结束时间2
	if endDate != "" {
		condition += " AND u.create_date<=? "
		params = append(params, endDate)
	}

	timediff := et.Sub(st) + time.Duration(24*time.Hour)
	res2 := make([]models.AnalysisChannel, int(timediff.Hours()/24))
	x := st
	for k, _ := range res2 {
		res2[k].Date = x
		dd, _ := time.ParseDuration("24h")
		x = x.Add(dd)
	}
	//1.注册用户，2.激活用户，3.注册到激活转化，4.点击-注册转化率，5.人均注册平台数
	switch code {
	case "1": //指定日期渠道注册用户数
		h5RegisteredUsersCount, err := models.GetH5RegisteredUsers_CountByTime(condition, channel, params) //注册用户
		if err != nil {
			resultMap["err"] = "获取注册用户数失败"
			return
		}

		for k, v := range res2 {
			for _, v1 := range h5RegisteredUsersCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					res2[k].RegisteredUsersCount = v1.RegisteredUsersCount
				}
			}
		}

	case "2": //指定日期渠道激活用户数
		h5ActivatedUsersCount, err := models.GetH5ActivatedUsers_CountByTime(condition, channel, params) //激活用户
		if err != nil {
			resultMap["err"] = "获取激活用户数失败"
			return
		}
		for k, v := range res2 {
			for _, v1 := range h5ActivatedUsersCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					res2[k].ActivatedUsersCount = v1.ActivatedUsersCount
				}
			}
		}
	case "3": //指定日期渠道注册激活转化率
		h5RegisteredUsersCount, err := models.GetH5RegisteredUsers_CountByTime(condition, channel, params) //注册用户
		if err != nil {
			resultMap["err"] = "获取注册用户数失败"
			return
		}
		h5ActivatedUsersCount, err := models.GetH5ActivatedUsers_CountByTime(condition, channel, params) //激活用户
		if err != nil {
			resultMap["err"] = "获取激活用户数失败"
			return
		}
		for k, v := range res2 {
			x := 0
			y := 0
			for _, v1 := range h5RegisteredUsersCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					x = v1.RegisteredUsersCount
				}
			}
			for _, v1 := range h5ActivatedUsersCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					y = v1.ActivatedUsersCount
				}
			}
			if x == 0 {
				res2[k].Rac = 0 //rac:Register activation conversion
			} else {
				res2[k].Rac = float64(y) / float64(x)
			}
		}
	case "4": //指定日期渠道点击注册转化率
		params1 := []string{}
		condition1 := ""
		condition1 += " AND u.out_put_source=?"

		//开始时间2
		if startDate != "" {
			condition1 += " AND plpr.create_date>=?"
			params1 = append(params1, startDate)
		}
		//结束时间2
		if endDate != "" {
			endDate += " 23:59:59"
			condition1 += " AND plpr.create_date<=?"
			params1 = append(params1, endDate)
		}
		h5HitsCount, err := models.GetHits_CountByTime(condition1, channel, params1) //点击次数
		if err != nil {
			resultMap["err"] = "获取点击次数失败"
			return
		}
		h5RegisteredUsersCount, err := models.GetH5RegisteredUsers_CountByTime(condition, channel, params) //注册人数
		if err != nil {
			resultMap["err"] = "获取注册人数失败"
			return
		}
		for k, v := range res2 {
			x := 0
			y := 0
			for _, v1 := range h5HitsCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					x = v1.HitsCount
					res2[k].HitsCount = v1.HitsCount
				}
			}
			for _, v1 := range h5RegisteredUsersCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					y = v1.RegisteredUsersCount
					res2[k].RegisteredUsersCount = v1.RegisteredUsersCount
				}
			}
			if x == 0 {
				res2[k].Hrc = 0
			} else {
				res2[k].Hrc = float64(y) / float64(x)
			}
		}
	case "5": //指定日期渠道人均注册平台数
		params1 := []string{}
		condition1 := ""
		condition1 += " AND u.out_put_source=?"

		//开始时间2
		if startDate != "" {
			condition1 += " AND urp.create_time>=?"
			params1 = append(params1, startDate)
		}
		//结束时间2
		if endDate != "" {
			endDate += " 23:59:59"
			condition1 += " AND urp.create_time<=?"
			params1 = append(params1, endDate)
		}
		h5RegisteredPlatformsCount, err := models.GetH5RegisteredPlatforms_CountByTime(condition1, channel, params1)
		if err != nil {
			resultMap["err"] = "获取用户注册平台数失败"
			return
		}
		h5ChannelUsersCount, err := models.GetH5RegisteredUsers_CountByTime(condition, channel, params)
		if err != nil {
			resultMap["err"] = "获取渠道用户数失败"
			return
		}
		for k, v := range res2 {
			x := 0.0
			y := 0
			for _, v1 := range h5RegisteredPlatformsCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					x = float64(v1.Count)
				}
			}
			for _, v1 := range h5ChannelUsersCount {
				if v.Date.Format("2006-01-02") == v1.Date.Format("2006-01-02") {
					y = v1.RegisteredUsersCount
				}
			}
			if y == 0 {
				res2[k].Crpc = 0
			} else {
				res2[k].Crpc = x / float64(y)
			}
		}
	}
	resultMap["flag"] = flag
	resultMap["ret"] = 200
	resultMap["msg"] = res2

}

//获取渠道列表
//@router /getchannellist [get]
func (c *ChannelAnalysisController) GetChannellist() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	resultMap["ret"] = 403
	channellist, err := models.GetAnalysisChannel("") //渠道名称列表
	if err != nil {
		resultMap["err"] = "获取渠道列表失败"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = channellist

	return
}
