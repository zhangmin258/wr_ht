package controllers

/*
*数据报表接口
*/

import (
	"wr_v1/models"
	"wr_v1/cache"
	"time"
	"wr_v1/utils"
	"fmt"
	"strconv"
	"net/http"
	"os"
	"github.com/tealeg/xlsx"
)

type DataReportController struct {
	BaseController
}

//平台数据明细
func (c *DataReportController) GetDataReport() {
	c.IsNeedTemplate()
	startTime := c.GetString("startTime2")
	endTime := c.GetString("endTime2")
	proId, err := c.GetInt("proId", 0)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品id失败！", err.Error(), c.Ctx.Input)
	}
	var pro models.OnlineProduct
	if proId == 0 {
		pro, err = models.GetOnlineProOne()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取默认产品id失败！", err.Error(), c.Ctx.Input)
		}
		proId = pro.Id
	} else {
		name, err := models.GetProductNameById(proId)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品失败！", err.Error(), c.Ctx.Input)
		}
		pro.Id = proId
		pro.Name = name
	}
	condition := ""
	resCon := ""
	params := []string{}
	if startTime != "" {
		resCon += " AND dd.date >=?"
		condition += " AND data_time >= ?"
		params = append(params, startTime)
	} else {
		startTime = time.Now().AddDate(0, 0, -30).Format(utils.FormatDate)
		condition += " AND data_time >= ?"
		resCon += " AND dd.date >=?"
		params = append(params, startTime)
	}
	if endTime != "" {
		endTime += " 23:59:59"
		resCon += " AND dd.date <=?"
		condition += " AND data_time <= ? "
		params = append(params, endTime)
	} else {
		endTime = time.Now().Format("2006-01-02")
	}
	dataReportShow, err := models.GetDataReport(proId, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品的数据报表失败！", err.Error(), c.Ctx.Input)
	}
	//获取平台返回注册
	restCount, err := models.GetRestCountByProId(proId, resCon, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品平台返回注册数据失败！", err.Error(), c.Ctx.Input)
	}
	for k, v := range dataReportShow {
		for _, v1 := range restCount {
			if v.DataTime == v1.Date {
				dataReportShow[k].RegisterCount = v1.Count
			}
		}
	}
	dataReportCount, dataReportAverage, dataReports := calculationData(dataReportShow)
	c.Data["dataReportCount"] = dataReportCount
	c.Data["dataReportAverage"] = dataReportAverage
	c.Data["startTime"] = startTime
	c.Data["endTime"] = endTime
	c.Data["product"] = pro
	c.Data["dataReport"] = dataReports
	c.Data["excelDate"] = time.Now().Format("2006-01-02")
	c.TplName = "data-operation/data_report.html"
}

//获取上线产品
func (c *DataReportController) GetOnlineProduct() {
	resultMap := make(map[string]interface{})
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	resultMap["ret"] = 403
	condition := ""
	params := []string{}
	name := c.GetString("name")
	if name != "" {
		condition = "AND p.name LIKE ? "
		params = append(params, "%"+name+"%")
	}
	product, err := models.GetOnlineProduct(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取上线产品失败！", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取上线产品失败！"
		return
	}
	resultMap["ret"] = 200
	resultMap["products"] = product
}

//导出全部数据
func (c *DataReportController) ExcelToAll() {
	startTime := c.GetString("startDate")
	endTime := c.GetString("endDate")
	if startTime == "" || endTime == "" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "开始或结束时间为空！", "开始："+startTime+"结束："+endTime, c.Ctx.Input)
	}
	condition := ""
	params := []string{}
	ddCon := " AND dd.date>=? AND dd.date<=? "
	condition += " AND cr.data_time>=? AND cr.data_time<=? "
	params = append(params, startTime, endTime)
	dataReportShow, err := models.GetDataReportAll(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品的数据报表失败！", err.Error(), c.Ctx.Input)
	}
	resgiterCountAll, err := models.GetRestCountAll(ddCon, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品平台返回注册失败！", err.Error(), c.Ctx.Input)
	}
	//获取所有产品
	names, err := models.GetOnlineProductName(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有上线产品失败！", err.Error(), c.Ctx.Input)
	}
	var dataReports []models.DataReport
	for k, v := range dataReportShow {
		for _, v1 := range resgiterCountAll {
			if v.ProId == v1.ProId && v.DataTime == v1.Date {
				dataReportShow[k].RegisterCount = v1.Count
			}
		}
	}

	for _, v1 := range dataReportShow {
		var dataReport models.DataReport
		dataReport.ProId = v1.ProId
		dataReport.DataTime = v1.DataTime
		dataReport.RegisterCount = v1.RegisterCount
		dataReport.ActivateUser = v1.ActivateUser
		dataReport.PlatformRegisterCount = v1.PlatformRegisterCount
		dataReport.AccessCount = v1.AccessCount
		dataReport.PageLoadCount = v1.PageLoadCount
		dataReport.ProName = v1.ProName
		dataReport.Sort = v1.Sort
		if v1.CpsFirstPer != 0 && v1.CpaPrice == 0 {
			dataReport.Price = "CPS:" + fmt.Sprintf("%.2f", v1.CpsFirstPer) + "%"
		} else if v1.CpsFirstPer == 0 && v1.CpaPrice != 0 {
			dataReport.Price = "CPA:" + fmt.Sprintf("%.2f", v1.CpaPrice)
		} else {
			dataReport.Price = "CPA+S:" + fmt.Sprintf("%.2f", v1.CpaPrice) + "+" + fmt.Sprintf("%.2f", v1.CpsFirstPer) + "%"
		}
		if v1.AccessCount != 0 {
			dataReport.UvRegister = float64(v1.RegisterCount * 100 / v1.AccessCount)
		} else {
			dataReport.UvRegister = 0
		}
		dataReport.Income = v1.CpaPrice*float64(v1.RegisterCount) + v1.CpsFirstPer*float64(v1.MakeLoanAmount)
		if dataReport.AccessCount != 0 {
			dataReport.PerUvIncome = dataReport.Income / float64(dataReport.AccessCount)
		} else {
			dataReport.PerUvIncome = 0
		}
		dataReports = append(dataReports, dataReport)
	}
	filename, err := exportDataReport(names, dataReports)
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

//导出对应产品数据
func (c *DataReportController) ExcelToOne() {
	condition := ""
	resCon := ""
	params := []string{}
	var exl [][]string
	if startDate := c.GetString("startDate"); startDate != "" {
		resCon += " AND dd.date >= ?"
		condition += " AND data_time >= ?"
		params = append(params, startDate)
	}
	if endDate := c.GetString("endDate"); endDate != "" {
		endDate += " 23:59:59"
		resCon += " AND dd.date <= ?"
		condition += " AND data_time <= ? "
		params = append(params, endDate)
	}
	proId, err := c.GetInt("proId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品id失败！", err.Error(), c.Ctx.Input)
	}
	dataReportShow, err := models.GetDataReport(proId, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品的数据报表失败！", err.Error(), c.Ctx.Input)
	}
	//获取平台返回注册
	restCount, err := models.GetRestCountByProId(proId, resCon, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品平台返回注册数据失败！", err.Error(), c.Ctx.Input)
	}
	proName, err := models.GetProductNameById(proId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品名称失败！", err.Error(), c.Ctx.Input)
	}
	for k, v := range dataReportShow {
		for _, v1 := range restCount {
			if v.ProId == v1.ProId && v.DataTime == v1.Date {
				dataReportShow[k].RegisterCount = v1.Count
			}
		}
	}
	exl = [][]string{{proName}}
	exp := []string{"日期", "平台位置", "PV", "UV", "我司统计注册", "我司统计激活", "平台返回注册", "UV注册转化", "当日价格", "收入", "每UV收益"}
	exl = append(exl, exp)
	dataReportCount, dataReportAverage, dataReports := calculationData(dataReportShow)
	for _, v := range dataReports {
		exp := []string{v.DataTime.Format("2006-01-02"), v.Sort, strconv.Itoa(v.PageLoadCount),
			strconv.Itoa(v.AccessCount), strconv.Itoa(v.PlatformRegisterCount), strconv.Itoa(v.ActivateUser), strconv.Itoa(v.RegisterCount), fmt.Sprintf("%.2f", v.UvRegister) + " %",
			v.Price, fmt.Sprintf("%.2f", v.Income), fmt.Sprintf("%.2f", v.PerUvIncome)}
		exl = append(exl, exp)
	}
	exp = []string{"合计", "-", strconv.Itoa(dataReportCount.PvAllCount), strconv.Itoa(dataReportCount.UvAllCount), strconv.Itoa(dataReportCount.PlatformRegisterAllCount),
		strconv.Itoa(dataReportCount.ActivateAllCount), strconv.Itoa(dataReportCount.RegisterAllCount), "-", "-", fmt.Sprintf("%.2f", dataReportCount.IncomeAll), "-"}
	exl = append(exl, exp)
	exp = []string{"平均", "-", fmt.Sprintf("%.2f", dataReportAverage.PvAverage), fmt.Sprintf("%.2f", dataReportAverage.UvAverage),
		fmt.Sprintf("%.2f", dataReportAverage.PlatformRegisterAverage), fmt.Sprintf("%.2f", dataReportAverage.ActivateAverage), fmt.Sprintf("%.2f", dataReportAverage.RegisterAverage),
		fmt.Sprintf("%.2f", dataReportAverage.UvRegisterAverage) + " %", "-", fmt.Sprintf("%.2f", dataReportAverage.IncomeAverage), fmt.Sprintf("%.2f", dataReportAverage.PerUvIncomeAverage)}
	exl = append(exl, exp)
	filename, err := utils.ExportDataReportOne(exl)
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

func calculationData(dataReportShow []models.DataReportShow) (dataReportCount models.DataReportCount, dataReportAverage models.DataReportAverage, dataReports []models.DataReport) {
	var (
		dataReport       models.DataReport
		uvRegisterCount  float64
		perUvIncomeCount float64
	)

	for _, v := range dataReportShow {
		dataReport.PageLoadCount = v.PageLoadCount
		dataReport.DataTime = v.DataTime
		dataReport.AccessCount = v.AccessCount
		dataReport.PlatformRegisterCount = v.PlatformRegisterCount
		dataReport.ActivateUser = v.ActivateUser
		dataReport.RegisterCount = v.RegisterCount
		dataReport.Sort = v.Sort
		if v.CpsFirstPer != 0 && v.CpaPrice == 0 {
			dataReport.Price = "CPS:" + fmt.Sprintf("%.2f", v.CpsFirstPer) + "%"
		} else if v.CpsFirstPer == 0 && v.CpaPrice != 0 {
			dataReport.Price = "CPA:" + fmt.Sprintf("%.2f", v.CpaPrice)
		} else {
			dataReport.Price = "CPA+S:" + fmt.Sprintf("%.2f", v.CpaPrice) + "+" + fmt.Sprintf("%.2f", v.CpsFirstPer) + "%"
		}
		if v.AccessCount != 0 {
			dataReport.UvRegister = float64(v.RegisterCount * 100 / v.AccessCount)
			uvRegisterCount += dataReport.UvRegister
		} else {
			dataReport.UvRegister = 0
		}
		dataReport.Income = v.CpaPrice*float64(v.RegisterCount) + v.CpsFirstPer*float64(v.MakeLoanAmount)
		if dataReport.AccessCount != 0 {
			dataReport.PerUvIncome = dataReport.Income / float64(dataReport.AccessCount)
			perUvIncomeCount += dataReport.PerUvIncome
		} else {
			dataReport.PerUvIncome = 0
		}
		dataReportCount.PvAllCount += dataReport.PageLoadCount
		dataReportCount.UvAllCount += dataReport.AccessCount
		dataReportCount.PlatformRegisterAllCount += dataReport.PlatformRegisterCount
		dataReportCount.RegisterAllCount += dataReport.RegisterCount
		dataReportCount.ActivateAllCount += dataReport.ActivateUser
		dataReportCount.IncomeAll += dataReport.Income
		dataReports = append(dataReports, dataReport)
	}
	count := len(dataReportShow)
	if count != 0 {
		dataReportAverage.PvAverage = float64(dataReportCount.PvAllCount) / float64(count)
		dataReportAverage.PerUvIncomeAverage = perUvIncomeCount / float64(count)
		dataReportAverage.UvRegisterAverage = uvRegisterCount / float64(count)
		dataReportAverage.UvAverage = float64(dataReportCount.UvAllCount) / float64(count)
		dataReportAverage.PlatformRegisterAverage = float64(dataReportCount.PlatformRegisterAllCount) / float64(count)
		dataReportAverage.RegisterAverage = float64(dataReportCount.RegisterAllCount) / float64(count)
		dataReportAverage.ActivateAverage = float64(dataReportCount.ActivateAllCount) / float64(count)
		dataReportAverage.IncomeAverage = dataReportCount.IncomeAll / float64(count)
	}
	return
}

//导出全部数据报表
func exportDataReport(names []string, datareport []models.DataReport) (filename string, err error) {
	//遍历exportToExcel
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	file = xlsx.NewFile()
	if len(names) == 0 {
		names = append(names, "未命名")
	}
	for _, name := range names {
		sheet, err = file.AddSheet(name)
		if err != nil {
			return "", err
		}
		var exl [][]string
		style := xlsx.NewStyle()
		style.Alignment.Horizontal = "center"
		style.Alignment.Vertical = "center"
		if name != "未命名" {
			exl = [][]string{{name}}
			exp := []string{"日期", "平台位置", "PV", "UV", "我司统计注册", "我司统计激活", "平台返回注册", "UV注册转化", "当日价格", "收入", "每UV收益"}
			exl = append(exl, exp)
		} else {
			exl = [][]string{{"暂无数据"}}
		}
		for _, v := range datareport {
			if name == v.ProName {
				exp := []string{v.DataTime.Format("2006-01-02"), v.Sort, strconv.Itoa(v.PageLoadCount),
					strconv.Itoa(v.AccessCount), strconv.Itoa(v.PlatformRegisterCount), strconv.Itoa(v.ActivateUser), strconv.Itoa(v.RegisterCount), fmt.Sprintf("%.2f", v.UvRegister) + " %",
					v.Price, fmt.Sprintf("%.2f", v.Income), fmt.Sprintf("%.2f", v.PerUvIncome)}
				exl = append(exl, exp)
			}
		}
		//遍历添加数据
		for k, v := range exl {
			row = sheet.AddRow()
			for k1, t := range v {
				cell = row.AddCell()
				if k == 0 {
					cell.Merge(10, 0)
					cell.SetStyle(style)
				} else {
					sheet.Cols[k1].Width = 16
				}
				cell.Value = t
			}
		}
	}
	filename = "数据报表-" + time.Now().Format("2006-01-02 15-04-05") + ".xlsx"
	err = file.Save(filename)
	return filename, err
}
