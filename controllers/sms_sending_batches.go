package controllers

import (
	"strconv"
	"strings"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type SMSSendingBatchesController struct {
	BaseController
}

//加载发送报表页面
//@router /getsendingbatches [get]
func (c *SMSSendingBatchesController) GetSendingBatches() {
	defer func() {
		c.IsNeedTemplate()
		c.TplName = "sms-management/sms_sending_batches.html"
	}()
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := make([]interface{}, 0, 2)

	if plateform := c.GetString("Source"); plateform != "" {
		p, err := strconv.Atoi(plateform)
		if err != nil {
			c.Data["err"] = "获取渠道信息失败"
			return
		}
		condition += " AND plateform = ? "
		params = append(params, p)
	}
	startTime := c.GetString("startTime")
	if startTime != "" {
		startTime = strings.Replace(startTime, "+", " ", -1)
		condition += " AND push_time >= ? "
		params = append(params, startTime)
	}
	endTime := c.GetString("endTime")
	if endTime != "" {
		endTime = strings.Replace(endTime, "+", " ", -1)
		condition += " AND push_time <= ?"
		params = append(params, endTime)
	}
	sendBatches, err := models.GetSmsMarketingData(condition, params, utils.StartIndex(pageNum, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取短信发送数据失败", err.Error(), c.Ctx.Input)
		c.Data["err"] = "获取短信发送数据失败"
		return
	}
	count, err := models.GetCountForSms(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取短信批次总条数失败", err.Error(), c.Ctx.Input)
		c.Data["err"] = "获取短信批次总条数失败"
		return
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
		c.Data["err"] = "获取需要的页数失败"
		return
	}
	c.Data["SendBatches"] = sendBatches
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count

}

//链接统计接口
//@router /userlinkonproduct [get]
func (c *SMSSendingBatchesController) UserLinkOnProduct() {
	defer c.ServeJSON()
	condition := ""
	params := []string{}
	sms_id := c.GetString("SmsId")
	source := c.GetString("Source")
	plateform := c.GetString("Plateform")
	pushTime := c.GetString("PushTime")
	startTime := c.GetString("startDate")
	endTime := c.GetString("endDate")
	if startTime != "" {
		condition += " AND pr.creat_time >= ? "
		params = append(params, startTime)
	}
	if endTime != "" {
		condition += " AND DATE_FORMAT(pr.creat_time,'%y-%m-%d') <= ? "
		params = append(params, endTime)
	}
	//拿到该条短信的标识
	smsIds, err := models.GetSMSFlag(sms_id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取短信的标识失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取短信的标识失败!"}
		return
	}
	// 获取到达成功数量
	msgCount, err := models.GetSuccessCount(smsIds)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取发送成功条数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取发送成功条数失败!"}
		return
	}
	// 判断source 为空
	if source == "" {
		c.Data["json"] = map[string]interface{}{"ret": 200, "LoadCount": 0, "ClickCount": 0, "ClickReach": 0, "LoadReach": 0, "ReachCount": msgCount, "Plateform": plateform, "PushTime": pushTime}
		return
	}
	clickCount, err := models.GetClickCount(source)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取点击次数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取点击次数失败!"}
		return
	}
	loadCount, err := models.GetLoadCount(source)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取加载次数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取加载次数失败!"}
		return
	}
	// 链接详情查看
	linkLoad, err := models.GetLinkStaLoadCount(source, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取链接加载次数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取链接加载次数失败!"}
		return
	}
	proClickList, err := models.GetLinkStaClickCount(source, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取链接点击次数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200, "err": err.Error(), "msg": "获取链接点击次数失败!"}
		return
	}
	//  分母不能为0
	if msgCount != 0 {
		clickReach := strconv.Itoa(clickCount*100/msgCount) + "%"
		loadReach := strconv.Itoa(loadCount*100/msgCount) + "%"
		c.Data["json"] = map[string]interface{}{"ret": 200, "LoadCount": loadCount, "ClickCount": clickCount, "ClickReach": clickReach, "LoadReach": loadReach, "ReachCount": msgCount,
			"Plateform": plateform, "PushTime": pushTime, "linkLoad": linkLoad, "proClickList": proClickList}
		return
	}
	clickReach := 0
	loadReach := 0
	c.Data["json"] = map[string]interface{}{"ret": 200, "LoadCount": loadCount, "ClickCount": clickCount, "ClickReach": clickReach, "LoadReach": loadReach, "ReachCount": msgCount,
		"Plateform": plateform, "PushTime": pushTime, "linkLoad": linkLoad, "proClickList": proClickList}
	return
}
