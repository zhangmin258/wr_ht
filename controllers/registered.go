package controllers

/*
*点击收益接口
*/
import (
	"wr_v1/models"
	"wr_v1/cache"
	"time"
	"sort"
	"strings"
	"fmt"
)

type RegisteredController struct {
	BaseController
}

func (c *RegisteredController) GetRegisterCount() {
	c.IsNeedTemplate()
	var params []string
	var condition = ""
	var restCon = ""
	endTime := c.GetString("endTime", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	startTime := c.GetString("startTime", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	fmt.Println(endTime, startTime)
	if startTime != "" {
		restCon += " AND dd.date >=?"
		condition += " AND c.data_time>=?"
	}
	params = append(params, startTime)
	if endTime != "" {
		endTime += " 23:59:59"
		restCon += " AND dd.date <=?"
		condition += " AND c.data_time<=?"
	}
	params = append(params, endTime)
	clickRegister, err := models.GetClickRegister(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取指定时间内上线的所有产品失败！", err.Error(), c.Ctx.Input)
	}
	//获取平台注册数
	registerCount, err := models.QueryRegisterCount(restCon, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品注册数量失败！", err.Error(), c.Ctx.Input)
	}
	clickRegisters := make(map[int]models.ClickRegister)
	for _, v := range clickRegister {
		if _, ok := clickRegisters[v.ProId]; ok {
			continue
		}
		var cr models.ClickRegister
		cr.ProId = v.ProId
		cr.CpaPrice = v.CpaPrice
		cr.Sort = v.Sort
		cr.Name = v.Name
		cr.Count = v.Count
		cr.CooperationType = v.CooperationType
		clickRegisters[v.ProId] = cr
	}
	var newClickRegister []models.ClickRegister
	for _, v := range clickRegisters {
		newClickRegister = append(newClickRegister, v)
	}
	for k, v := range newClickRegister {
		for _, v1 := range clickRegister {
			if v.ProId == v1.ProId {
				newClickRegister[k].AccessCount += v1.AccessCount
				newClickRegister[k].ActivateUser += v1.ActivateUser
				newClickRegister[k].PlatformRegisterCount += v1.PlatformRegisterCount
			}
		}
		for _, v2 := range registerCount {
			if v.ProId == v2.ProId {
				newClickRegister[k].RegisterCount += v2.Count
			}
		}
	}
	sort.Sort(models.ProductInfos(newClickRegister))
	var sumActivateUser, sumAccessCount, sumRegisterCount, sumPlatformRegisterCount int
	var sumPriceCount, sumPrice float64
	var list []*models.ProIncome
	for k, v := range newClickRegister {
		sumAccessCount += v.AccessCount
		sumActivateUser += v.ActivateUser
		sumRegisterCount += v.RegisterCount
		sumPlatformRegisterCount += v.PlatformRegisterCount
		/*if v.PlatformRegisterCount == 0 || v.RegisterCount == 0 {
			newClickRegister[k].DataRisk = 0
		} else {
			newClickRegister[k].DataRisk = float64(v.RegisterCount) / float64(v.PlatformRegisterCount)
		}*/
		newClickRegister[k].AllEarnings = v.CpaPrice * float64(v.PlatformRegisterCount) // 总收益
		sumPrice += newClickRegister[k].AllEarnings
		if v.AccessCount != 0 && v.CooperationType == 1 {
			newClickRegister[k].ClickEarnings = v.CpaPrice * float64(v.PlatformRegisterCount) / float64(v.AccessCount)
		} else {
			newClickRegister[k].ClickEarnings = 0
		}
		sumPriceCount += newClickRegister[k].ClickEarnings
		if !strings.Contains(newClickRegister[k].Sort, "大额贷款") && newClickRegister[k].CooperationType != 0 && newClickRegister[k].ClickEarnings != 0 {
			list = append(list, &models.ProIncome{newClickRegister[k].ProId, newClickRegister[k].ClickEarnings})
		}
	}
	sort.Sort(models.Pros(list))
	for i := 0; i < len(list); i++ {
		for j := 0; j < len(newClickRegister); j++ {
			if newClickRegister[j].ProId == list[i].ProductId {
				newClickRegister[j].Sorting = i + 1
			}
		}
	}
	c.Data["sumAccessCount"] = sumAccessCount
	c.Data["sumActivateUser"] = sumActivateUser
	c.Data["sumRegisterCount"] = sumRegisterCount
	c.Data["sumPlatformRegisterCount"] = sumPlatformRegisterCount
	c.Data["sumPriceCount"] = sumPriceCount
	c.Data["sumPrice"] = sumPrice
	c.Data["proDates"] = newClickRegister
	c.Data["startTime"] = startTime
	c.Data["endTime"] = strings.Split(endTime, " ")[0]
	c.TplName = "data-operation/reg_statistics.html"
}
