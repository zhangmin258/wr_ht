package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/email"
)

type SupportLoanController struct {
	BaseController
}

//付费用户
func (this *SupportLoanController) LoanPlanData() {
	this.IsNeedTemplate()
	//获取贷款稳下付费信息
	moneyPrice, err := models.GetSupportLoanInfo(14)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取贷款稳下付费价格异常", err.Error(), this.Ctx.Input)
	}
	//付费人次
	payPersonCount, err := models.GetPayPersonCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取付费人次异常", err.Error(), this.Ctx.Input)
	}
	//付费金额
	payMoneyCount := moneyPrice * float64(payPersonCount)
	//退款金额
	refundAmount, err := models.GetRefundAmount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取退款金额异常", err.Error(), this.Ctx.Input)
	}
	//累计总收入
	totalAmount := payMoneyCount - refundAmount

	//付费用户列表
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	condition := ""
	params := []string{}
	phoneNum := this.GetString("phone_number")
	if phoneNum != "" {
		params = append(params, phoneNum)
		condition += " AND u.account=?"
	}
	payUsersList, err := models.GetPayUsers(condition, params, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取付费用户列表异常", err.Error(), this.Ctx.Input)
	}
	listCount, err := models.GetPayUsersCount(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取付费用户数量异常", err.Error(), this.Ctx.Input)
	}
	pageCount := utils.PageCount(listCount, utils.PageSize10)
	this.Data["payUsersList"] = payUsersList
	this.Data["count"] = listCount
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = page
	this.Data["payPersonCount"] = payPersonCount
	this.Data["payMoneyCount"] = payMoneyCount
	this.Data["refundAmount"] = refundAmount
	this.Data["totalAmount"] = totalAmount
	this.TplName = "loan-under-steady/paying_user.html"
}

//退款
func (this *SupportLoanController) Refund() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	var refund models.Refund
	err := this.ParseForm(&refund)
	payOrder, err := models.GetPayToken(refund.Uid, refund.CreateTime)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取付款凭证异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "获取付款凭证异常！"
		return
	}
	if payOrder == "" {
		resultMap["msg"] = "未完成付款或付款时异常，无法退款！"
		return
	}
	money, err := strconv.ParseFloat(refund.RefundAmount, 64)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取退款金额异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "获取获取退款金额异常！"
		return
	}
	refundState, err := models.GetUsersRefundState(refund.Uid, refund.CreateTime)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取获取贷款稳下订单状态异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "获取贷款稳下订单状态异常！"
		return
	}
	if refundState.IsFinished != 1 || refundState.IsValid != 1 {
		resultMap["msg"] = "订单已失效或已完成，无法退款！"
		return
	}
	//获取贷款稳下付费信息
	moneyPrice, err := models.GetSupportLoanInfo(14)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取贷款稳下付费价格异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "获取贷款稳下付费信息失败"
		return
	}
	if money > moneyPrice {
		resultMap["msg"] = "退款金额不能大于付费金额"
		return
	}
	err = models.ChangeSupportLoanState(refund.Uid, refund.CreateTime)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "修改贷款稳下订单状态异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "修改贷款稳下订单状态异常！"
		return
	}
	err = models.WalletRecharge(refund.Uid, money)
	if err != nil {
		email.Send("后台贷款稳下页面退款异常", "err:"+err.Error(), utils.ToUsers, "weirong")
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "退款至用户钱包异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "退款至用户钱包异常！"
		return
	}
	err = models.AddRefundRecord(refund.Uid, 14, money, payOrder)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "添加退款记录异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "添加退款记录异常！"
		return
	}
	//获取用户钱包余额
	balance, err := models.GetWalletBalance(refund.Uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取用户钱包余额异常", err.Error(), this.Ctx.Input)
	}
	//获取用户融豆余额
	score, err := models.GetScore(refund.Uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "获取用户融豆余额异常", err.Error(), this.Ctx.Input)
	}
	//添加用户钱包收支记录
	var financeInfo = models.FinanceInfo{
		Uid:               refund.Uid,
		DealType:          15,
		MoneyAmount:       money,
		PayOrGet:          2,
		PayType:           2,
		CreateTime:        time.Now(),
		ServiceStates:     0,
		BeforeScoreAmount: score,
		AfterScoreAmount:  score,
		BeforeMoneyAmount: balance,
		AfterMoneyAmount:  balance + money,
	}
	err = models.AddRefundFinance(financeInfo)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "贷款稳下退款添加用户钱包收支记录异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "添加用户钱包收支记录异常!"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "退款成功！"
	return
}

//导出到Excel
func (this *SupportLoanController) ApplyExportExcel() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	condition := ""
	params := []string{}
	startTime := this.GetString("startTime")
	endTime := this.GetString("endTime")
	if startTime != "" {
		condition += " AND DATE_FORMAT(s.create_time ,'%Y-%m-%d') >=?"
		params = append(params, startTime)
	}
	if endTime != "" {
		condition += " AND DATE_FORMAT(s.create_time ,'%Y-%m-%d')<=?"
		params = append(params, endTime)
	}
	list, err := models.ApplyExportExcel(condition, params)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "导出Excel时查询数据异常", err.Error(), this.Ctx.Input)
		resultMap["err"] = "导出Excel时查询数据异常！"
		return
	}
	loanData := [][]string{{"下单时间", "手机号码", "姓名"}}
	colWidth := []float64{24.0, 20.0, 16.0}
	for _, v := range list {
		listData := []string{v.CreateTime, v.Account, v.Name}
		loanData = append(loanData, listData)
	}
	fileName, err := utils.SupportLoanExportToExcel(loanData, colWidth, "付费用户列表"+time.Now().Format("2006-01-02"))
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "导出Excel异常", err.Error(), this.Ctx.Input)
		return
	}
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fileName)
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, fileName)
	err = os.Remove(fileName)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}
