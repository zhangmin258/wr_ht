package controllers

import (
	"fmt"
	"strconv"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/services"
	"wr_v1/utils"
	"zcm_tools/email"
)

/*
*用户佣金接口
 */

type UserDepositController struct {
	BaseController
}

//待审核提现列表
func (c *UserDepositController) DepositList() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//手机号/账号
	if account := c.GetString("account"); account != "" {
		condition += " AND u.account LIKE ? "
		params = append(params, "%"+account+"%")
	}
	//待审批用户信息
	condition += " AND ur.type = 2 "
	userDeposit, err := models.GetUserDeposit(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取佣金列表失败", err.Error(), c.Ctx.Input)
		}
	}
	count, err := models.GetUserDepositCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有佣金数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userDeposit"] = userDeposit
	c.TplName = "pay-management/commission_withdrawal.html"
}

//邀请记录
func (c *UserDepositController) InvitationInfo() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	userId, err := c.GetInt("userId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户id失败", err.Error(), c.Ctx.Input)
	}

	condition := ""
	params := []string{}
	//手机号/账号
	if account := c.GetString("account"); account != "" {
		condition += " and u.account LIKE ? "
		params = append(params, "%"+account+"%")
	}
	userInvitation, err := models.UserInvitationList(userId, condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户邀请好友信息失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.UserInvitationCount(userId, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户邀请好友数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	//获取用户信息
	user, err := models.GetUserById(userId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["user"] = user
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userInvitation"] = userInvitation
	c.TplName = "show-management/invited_record.html"
}

// 审批记录
func (c *UserDepositController) ApproveInfo() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//手机号/账号
	if account := c.GetString("account"); account != "" {
		condition += " AND u.account LIKE ? "
		params = append(params, "%"+account+"%")
	}
	userDeposit, err := models.GetUserDeposit(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取佣金列表失败", err.Error(), c.Ctx.Input)
		}
	}
	count, err := models.GetUserDepositCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有佣金数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userDeposit"] = userDeposit
	c.TplName = "pay-management/examine_approve_record.html"
}

// 提现明细
func (c *UserDepositController) WithdrawInfo() {
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	userId, err := c.GetInt("userId")
	fmt.Println("userId:", userId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户id失败", err.Error(), c.Ctx.Input)
	}
	userDepositDetail, err := models.GetUserDepositDetail(userId, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户提现明细失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetUserDepositDetailCount(userId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有佣金数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败", err.Error(), c.Ctx.Input)
	}
	//获取用户信息
	user, err := models.GetUserById(userId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户失败", err.Error(), c.Ctx.Input)
	}
	fmt.Println("user数据:", user)
	fmt.Println("审批记录数据明细：", userDepositDetail)
	c.Data["user"] = user
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userDepositDetail"] = userDepositDetail
	c.TplName = "pay-management/withdrawal_subsidiary.html"
}

// 用户的提现申请的操作
func (c *UserDepositController) WithdrawDeposit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	router := "userdeposit/withdrawdeposit"
	sysId := c.User.Id
	jsonstr := string(c.Ctx.Input.RequestBody)
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	depId, err := c.GetInt("depId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "获取用户id失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取用户佣金id失败"
		return
	}
	// 防止重复提交
	if !utils.Rc.SetNX(utils.CACHE_KEY_ChECKWITHDRAWCASH_DEPID+strconv.Itoa(depId), 1, time.Minute) {
		resultMap["err"] = "请不要重复提交哦"
		return
	}
	//获取用户提现信息
	userDeposit, err := models.GetDepositById(depId)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "获取用户提现信息失败", err.Error(), c.Ctx.Input)
		}
		resultMap["err"] = "获取用户提现信息失败"
		return
	}
	//获取用户银行卡信息
	bc, err := models.GerUsersBankcardById(userDeposit.Uid)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "获取用户银行卡信息失败", err.Error(), c.Ctx.Input)
		}
		resultMap["err"] = "获取用户银行卡信息失败"
		return
	}
	//获取用户实名认证信息
	um, err := models.GetUsersMetadataById(userDeposit.Uid)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "获取用户银行卡信息失败", err.Error(), c.Ctx.Input)
		}
		resultMap["err"] = "查询用户实名信息错误"
		return
	}
	// 发起提现
	res, err := services.LLTrade(userDeposit.OrderCode, bc.CardNumber, um.VerifyRealName, utils.SubFloatToString(userDeposit.Real_amount, 2))
	fmt.Println(res)
	// 主动查询
	go services.CostActiveQuery(userDeposit.OrderCode, sysId)

	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "发送请求过程发生异常", err.Error(), c.Ctx.Input)
		email.Send("发送请求过程发生异常!", err.Error(), utils.ToUsers, c.Ctx.Input.IP()+jsonstr+router)
		resultMap["err"] = "发送请求过程发生异常!"
		return
	}
	// 发送成功
	if res.Ret_code == "0000" {
		resultMap["ret"] = 200
		resultMap["msg"] = "发送提现申请成功!"
	}

	// 作为失败处理
	if res.Ret_code != "0000" && res.Ret_code != "9999" && res.Ret_code != "4006" && res.Ret_code != "4007" && res.Ret_code != "4009" && res.Ret_code != "1002" && res.Ret_code != "2005" {
		go services.LoanResultUpdate(userDeposit.OrderCode, res.Ret_code, "", sysId)
		resultMap["ret"] = 403
		resultMap["err"] = "当前订单状态异常!已取消提现操作！"
		if res.Ret_msg != "" {
			resultMap["err"] = res.Ret_msg
		}
	}
	return
}

// 用户的拒绝提现申请的操作
func (c *UserDepositController) RefuseWithdrawDeposit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	router := "userdeposit/withdrawdeposit"
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	depId, err := c.GetInt("depId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "获取用户id失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取用户佣金id失败"
		return
	}
	if !utils.Rc.SetNX(utils.CACHE_KEY_ChECKWITHDRAWCASH_DEPID+strconv.Itoa(depId), 1, time.Minute) {
		resultMap["err"] = "请不要重复提交哦"
		return
	}
	//获取用户提现信息
	userDeposit, err := models.GetDepositById(depId)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "获取用户提现信息失败", err.Error(), c.Ctx.Input)
		}
		resultMap["err"] = "获取用户提现信息失败"
		return
	}
	sysId := c.User.Id
	err = models.RefuseWithdrawDeposit(userDeposit, sysId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, router, "拒绝提现失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "拒绝提现失败"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "拒绝提现成功！"
	return
}
