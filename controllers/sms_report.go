package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
	"wr_v1/models"
)

type SmsReportController struct {
	beego.Controller
}

/*
	短信回调接口
*/

//空间畅想短信回调
//@router /kjcxsmsreport [get]
func (c *SmsReportController) KJCXSmsReport() {
	defer c.ServeJSON()
	//resuleMap := make(map[string]interface{})
	name := c.GetString("name")
	report := c.GetString("report")
	reportInfo := strings.Split(report, ";")
	var msgReport []models.MsgReport
	for _, v := range reportInfo {
		var oncereport models.MsgReport
		reportOne := strings.Split(v, ",")
		if len(reportOne) == 4 {
			oncereport.SmsId = reportOne[0]
			oncereport.PhoneNumber = reportOne[1]
			oncereport.StatusReport = reportOne[2]
			oncereport.PushTime = reportOne[3]
			msgReport = append(msgReport, oncereport)
		} else {
			reportOne := strings.Split(v, ",")
			oncereport.SmsId = reportOne[0]
			oncereport.PhoneNumber = reportOne[1]
			oncereport.StatusReport = reportOne[2]
			oncereport.Reference = reportOne[3]
			oncereport.PushTime = reportOne[4]
			msgReport = append(msgReport, oncereport)
		}
	}
	err := models.AddReportData(name, msgReport)
	if err != nil {
		c.Data["json"] = "error:" + err.Error()
		return
	}
	c.Data["json"] = "success"
	return
}

//云融正通短信回调
//@router /yrztsmsreport [post]
func (c *SmsReportController) YRZTSmsReport() {
	defer c.ServeJSON()
	var getYRZTReport []models.GetYRZTReport
	var yrztReport models.GetYRZTReport
	userName := c.GetString("userName")
	passWord := c.GetString("passWord")
	messageQty := c.GetString("messageQty")
	count, err := strconv.Atoi(messageQty)
	if err != nil {
		c.Data["json"] = "-1#解析状态报告条数异常"
		return
	}
	for i := 1; i <= count; i++ {
		yrztReport.UserName = userName
		yrztReport.PassWord = passWord
		yrztReport.MessageQty = messageQty
		j := strconv.Itoa(i)
		yrztReport.SubmitMessageId = c.GetString("submitMessageId" + j)
		yrztReport.ClientMessageBatchId = c.GetString("clientMessageBatchId" + j)
		yrztReport.MobilePhone = c.GetString("MobilePhone" + j)
		timeStr, _ := time.Parse("20060102150405", c.GetString("dateTimeStr"+j))
		yrztReport.DateTimeStr = timeStr.Format("2006-01-02 15:04:05")
		yrztReport.DeliveryStatus = c.GetString("deliveryStatus" + j)
		yrztReport.DeliveryStatusCode = c.GetString("deliveryStatusCode" + j)
		getYRZTReport = append(getYRZTReport, yrztReport)
	}
	err = models.AddYRZTReportData(getYRZTReport)
	if err != nil {
		c.Data["json"] = "-1#保存状态报告失败"
		return
	}
	c.Data["json"] = "0#成功"
}
