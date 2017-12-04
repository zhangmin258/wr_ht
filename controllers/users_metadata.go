package controllers

import (
	"encoding/json"
	"strconv"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/crypt"
	"zcm_tools"
	"zcm_tools/http"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

/*
认证用户接口
*/
type UserMetadataController struct {
	BaseController
}

//用户详细信息
func (c *UserMetadataController) PersonalInformation() {
	//设置整体加载
	c.IsNeedTemplate()
	id := c.GetString("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.Ctx.WriteString("id错误")
	}
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户详细信息失败", err.Error(), c.Ctx.Input)
	}
	apiMoney, err := models.GetUserMoney(uid, 0)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户收益信息失败", err.Error(), c.Ctx.Input)
	}
	h5Money, err := models.GetUserMoney(uid, 1)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户收益信息失败", err.Error(), c.Ctx.Input)
	}
	userMoney := models.UserMoney{h5Money, apiMoney}
	c.Data["id"] = id
	c.Data["user"] = user
	c.Data["userMoney"] = userMoney
	c.TplName = "user-management/personal_information.html"
}

//@router /proOrderDetails [get]
func (c *UserMetadataController) ProOrderDetails() {
	c.IsNeedTemplate()
	id := c.GetString("id")
	list := models.GetProOrderlist(id)
	c.Data["userorder"] = list
	c.TplName = "user-management/proOrderDetails.html"

}

//获取用户手机运营商数据
//@router /usersoperatedata [get]
func (c *UserMetadataController) UsersOperateData() {
	c.IsNeedTemplate()
	uid, _ := c.GetInt("id")
	//if uid == 1 {
	//	return
	//}
	//fmt.Println(uid)
	var userReportResult models.UserReportResult
	err := utils.MgoSession.DB(utils.MgoDbName).C("ydreportdata").Find(&models.MonGoQuery{Uid: uid}).One(&userReportResult) //100000257
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "mgodb数据获取异常", err.Error(), c.Ctx.Input)
		c.Data["msg"] = "没有手机运营商记录"
		c.Data["name"] = "手机运营商"
		c.TplName = "user-management/user_info_show.html"
		return
	}
	if len(userReportResult.ReportResult.Data.Friend_circle.Peer_num_top_list) != 0 {
		c.Data["calls_3"] = userReportResult.ReportResult.Data.Friend_circle.Peer_num_top_list[0].Top_item
	} else {
		c.Data["msg"] = "没有手机运营商记录"
		c.Data["name"] = "手机运营商"
		c.TplName = "user-management/user_info_show.html"
		return
	}

	if len(userReportResult.ReportResult.Data.Friend_circle.Peer_num_top_list) >= 2 {
		c.Data["calls_6"] = userReportResult.ReportResult.Data.Friend_circle.Peer_num_top_list[1].Top_item
	}

	if len(userReportResult.ReportResult.Data.Friend_circle.Location_top_list) != 0 {
		c.Data["active_location"] = userReportResult.ReportResult.Data.Friend_circle.Location_top_list[0].Top_item
	}

	if len(userReportResult.ReportResult.Data.Cell_behavior) != 0 {
		c.Data["active_location_mon"] = userReportResult.ReportResult.Data.Cell_behavior[0].Behavior
	}

	if len(userReportResult.ReportResult.Data.Contact_region) != 0 {
		c.Data["contact_active_location"] = userReportResult.ReportResult.Data.Contact_region[0].Region_list
	}
	if len(userReportResult.ReportResult.Data.Call_duration_detail) != 0 {
		c.Data["time_collection"] = userReportResult.ReportResult.Data.Call_duration_detail[0].Duration_list
	}

	c.TplName = "user-management/operator_data.html"
}

//获取用户有盾实名数据
//@router /usersmetadata [get]
func (c *UserMetadataController) UsersMetaData() {
	c.IsNeedTemplate()
	uid, _ := c.GetInt("id")
	var ydudcredit models.YdudcreditResult
	var ydudcredit1 models.MdbYDYYSUdcreditData
	credit := utils.MgoSession.DB(utils.MgoDbName).C("ydudcredit")
	uidMap := make(map[string]int)
	uidMap["uid"] = uid
	err := credit.Find(uidMap).One(&ydudcredit)
	if err != nil && err.Error() != "not found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取有盾信息异常", err.Error(), c.Ctx.Input)
		c.Data["msg"] = "没有有盾实名数据"
		c.Data["name"] = "有盾实名数据"
		c.TplName = "user-management/user_info_show.html"
		return
	}
	if ydudcredit.Ydudcredit.Header.Ret_code == "" {
		uidMap := make(map[string]string)
		uidMap["uid"] = strconv.Itoa(uid)
		err1 := credit.Find(uidMap).One(&ydudcredit1)
		if err1 != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取有盾信息异常", err1.Error(), c.Ctx.Input)
			c.Data["msg"] = "没有有盾实名数据"
			c.Data["name"] = "有盾实名数据"
			c.TplName = "user-management/user_info_show.html"
			return
		}
		ydudcredit.Ydudcredit = ydudcredit1.YDYYSUdcreditYdudcredit
	}
	pageNum, _ := c.GetInt("pageNum")
	pageCount, _ := c.GetInt("pageCount")
	if pageNum < 1 {
		pageNum = 1
	}
	if len(ydudcredit.Ydudcredit.Body.Devices)%10 > 0 {
		pageCount = len(ydudcredit.Ydudcredit.Body.Devices)/10 + 1
	} else {
		pageCount = len(ydudcredit.Ydudcredit.Body.Devices) / 10
	}
	startpage := (pageNum - 1) * 10
	endpage := pageNum*10 - 1

	if len(ydudcredit.Ydudcredit.Body.Devices) != 0 {
		c.Data["device"] = ydudcredit.Ydudcredit.Body.Devices
		c.Data["device1"] = ydudcredit.Ydudcredit.Body.Devices[0]
	} else {
		device := make([]models.Devices, 0)
		c.Data["device"] = device
	}
	//获取实名信息
	idenuser := models.GetIdentificationUser(uid)
	if idenuser == nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取有盾信息异常", "", c.Ctx.Input)
		}
		c.Data["msg"] = "没有有盾实名数据"
		c.Data["name"] = "有盾实名数据"
		c.TplName = "user-management/user_info_show.html"
		return
	}
	c.Data["idenuser"] = idenuser
	c.Data["user_detail"] = ydudcredit.Ydudcredit.Body.User_detail
	c.Data["body"] = ydudcredit.Ydudcredit.Body
	c.Data["ydudcredit"] = ydudcredit
	c.Data["count"] = len(ydudcredit.Ydudcredit.Body.Devices)
	c.Data["startpage"] = startpage
	c.Data["endpage"] = endpage
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.TplName = "user-management/youdun_data.html"
}

//用户档案查询
//@router /searcherusersmetadata [post]
func (c *UserMetadataController) GetUserYdUdcreditData() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	uid, err := c.GetInt("userId")
	if err != nil {
		resultMap["err"] = "userId获取错误！"
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "userId获取错误", err.Error(), c.Ctx.Input)
		return
	}
	// 查询mongodb 是否有数据 有的话删除
	var ydudcredit models.MdbYDYYSUdcreditData
	credit := utils.MgoSession.DB(utils.MgoDbName).C("ydudcredit")
	uidMap := make(map[string]string)
	uidMap["uid"] = strconv.Itoa(uid)
	err = credit.Find(uidMap).One(&ydudcredit)
	if err != nil {
		if err.Error() != "not found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取有盾信息异常", err.Error(), c.Ctx.Input)
			resultMap["err"] = "有盾用户档案查询错误！"
			return
		}
	} else {
		credit.RemoveAll(bson.M{"uid": strconv.Itoa(uid)})
	}
	idCard := c.GetString("idCard")
	params := map[string]string{}
	params["id_no"] = idCard
	params["pubkey"] = utils.YD_HWY_API_KEY
	params["product_code"] = "carrier"
	code := utils.NewUUID().Hex()
	params["out_order_id"] = code
	s := zcm_tools.ToString(params)
	md5Str := crypt.MD5(s + "|" + "HT1oh4pDajElf1cOazGi")
	creditUrl := utils.YD_HWY_UDCREDIT_URL + "pubkey/" + utils.YD_HWY_API_KEY + "/product_code/" + "Y1001003" + "/out_order_id/" + code + "/signature/"
	resp, err := http.Post(creditUrl+md5Str, s, "application/json")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "有盾用户档案查询错误", err.Error(), c.Ctx.Input)
		resultMap["err"] = "有盾用户档案查询错误！"
		return
	}
	var ydc models.Ydudcredit
	err = json.Unmarshal(resp, &ydc)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析有盾用户档案信息返回数据错误", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析有盾用户档案信息返回数据错误！"
		return
	}
	ydcr := new(models.MdbYDYYSUdcreditData)
	ydcr.Uid = strconv.Itoa(uid)
	ydcr.YDYYSUdcreditYdudcredit = ydc
	authStatus := 3
	if ydc.Header.Ret_code == "000000" {
		authStatus = 2
		if utils.RunMode == "release" {
			ydcredit := utils.MgoSession.DB(utils.MgoDbName).C("ydudcredit")
			err = ydcredit.Insert(ydcr)
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "有盾用户档案信息写入MgoDB失败", err.Error(), c.Ctx.Input)
				resultMap["err"] = "解析有盾用户档案信息返回数据错误！"
				return
			}
		}
	}
	err = models.UpdateYdUdCredit(uid, authStatus)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改用户档案认证信息失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析有盾用户档案信息返回数据错误！"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "同步成功"
}

//获取用户支付宝数据
//@router /usersalipaydata [get]
func (c *UserMetadataController) UsersAliPayData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有支付宝记录"
	c.Data["name"] = "支付宝"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户淘宝数据
//@router /userstaobaodata [get]
func (c *UserMetadataController) UsersTaoBaoData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有淘宝记录"
	c.Data["name"] = "淘宝"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户京东数据
//@router /usersjingdongdata [get]
func (c *UserMetadataController) UsersJingDongData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有京东记录"
	c.Data["name"] = "京东"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户学信网数据
//@router /usersletterdata [get]
func (c *UserMetadataController) UsersLetterData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有学信网记录"
	c.Data["name"] = "学信网"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户人行征信数据
//@router /userspersondata [get]
func (c *UserMetadataController) UsersPersonData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有人行征信记录"
	c.Data["name"] = "人行征信"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户银行卡数据
//@router /userscarddata [get]
func (c *UserMetadataController) UsersCardData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有银行卡记录"
	c.Data["name"] = "银行卡"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户银行流水数据
//@router /userswaterdata [get]
func (c *UserMetadataController) UsersWaterData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有银行流水记录"
	c.Data["name"] = "银行流水"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户信用卡账单数据
//@router /usersbilldata [get]
func (c *UserMetadataController) UsersBillData() {
	c.IsNeedTemplate()
	c.Data["msg"] = "没有信用卡账单记录"
	c.Data["name"] = "信用卡账单"
	c.TplName = "user-management/user_info_show.html"
}

//获取用户紧急联系人
//@router /userscontact [get]
func (c *UserMetadataController) UsersContact() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/emergency_contact.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	list, err := models.GetEmergencyContactList(uid, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取联系人失败", err.Error(), c.Ctx.Input)
	}
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetContactCount(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取数量失败", err.Error(), c.Ctx.Input)
	}
	contactcount := utils.PageCount(count, utils.PageSize10)
	c.Data["emergencylist"] = list
	c.Data["user"] = user
	c.Data["page"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["contactcount"] = contactcount
	c.Data["count"] = count
}

//获取用户登录消息数据
//@router /userslogindata [get]
func (c *UserMetadataController) UsersLoginData() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/login_information.html"
	}()
	uid, _ := c.GetInt("id")
	//登录记录
	condition := ""
	pars := []string{}
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	loginRcds, err := models.GetLoginRecordList(uid, condition, pars, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取登陆历史失败", err.Error(), c.Ctx.Input)
	}
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	logincount := models.GetLoginRecordCount(uid, condition, pars)
	loginpagecount := utils.PageCount(logincount, utils.PageSize10)
	c.Data["loginpagecount"] = loginpagecount
	c.Data["loginpage"] = page
	c.Data["logincount"] = logincount
	c.Data["loginRcds"] = loginRcds
	c.Data["pagesize"] = utils.PageSize10
	c.Data["user"] = user
}

//获取用户贷款记录数据
//@router /usersloandata [get]
func (c *UserMetadataController) UsersLoanData() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/loan_records.html"
	}()
	uid, _ := c.GetInt("id")
	condition := " and l.uid=? "
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	//贷款详情
	pars := []string{strconv.Itoa(uid)}
	list, err := models.GetLoanList(condition, utils.StartIndex(page, utils.PageSize10), utils.PageSize10, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取贷款详情失败", err.Error(), c.Ctx.Input)
	}
	count := models.GetLoanCount(condition, pars...)
	condition += " and l.state in (?,?,?,?,?) "
	//SUPPLEMENT:代补充资料，REFUSE：审核未通过，WAITING：待审核，PASS：审核通过，PAYING：待放款，PAYFAILED:放款失败，CONFIRM：已放款，FINISH：已完成，OVERDUE：逾期中

	//借款成功的总数
	pars = append(pars, "PASS")
	pars = append(pars, "PAYING")
	pars = append(pars, "CONFIRM")
	pars = append(pars, "FINISH")
	pars = append(pars, "OVERDUE")
	passCount := models.GetLoanCount(condition, pars...)
	//总借款金额
	maxMoney := models.GetLoanMoneyMax(condition, pars...)

	//还款成功的总数
	condition = " and l.uid= ? and l.state in (?) "
	pars = pars[:1]
	pars = append(pars, "FINISH")
	finishCount := models.GetLoanCount(condition, pars...)
	//逾期总数
	pars = pars[:1]
	pars = append(pars, "OVERDUE")
	overdueCount := models.GetLoanCount(condition, pars...)

	loanpagecount := utils.PageCount(count, utils.PageSize10)
	token := c.GetString("token")
	if utils.MD5("demaxiya"+strconv.Itoa(uid)+time.Now().Format(utils.FormatDate)) == token {
		c.Data["repayment"] = true
	}
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["loancount"] = count
	c.Data["loanpage"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["loanpagecount"] = loanpagecount
	c.Data["loanlist"] = list

	//各种总数
	c.Data["passcount"] = passCount
	c.Data["finishcount"] = finishCount
	c.Data["overduecount"] = overdueCount
	c.Data["maxmoney"] = fmt.Sprintf("%1.2f", maxMoney/10000)
	c.Data["user"] = user
}

//获取用户消息中心数据
//@router /usersmessagecenter [get]
func (c *UserMetadataController) UsersMessageCenter() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/message_center.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	//用户消息详情
	list, err := models.GetUserMessageList(uid, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户消息列表失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetUserMessageCount(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户消息总数失败", err.Error(), c.Ctx.Input)
	}
	loanpagecount := utils.PageCount(count, utils.PageSize10)
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["messcount"] = count
	c.Data["messpage"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["messpagecount"] = loanpagecount
	c.Data["messlist"] = list
	c.Data["user"] = user
}

//获取邀请好友信息
//@router /invitefriend [get]
func (c *UserMetadataController) InviteFriend() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/invite_friend.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	list, err := models.GetInviteFriendList(uid, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户邀请好友信息失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetInviteFriendListCount(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户邀请好友总数失败", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PageSize10)
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户邀请好友信息失败", err.Error(), c.Ctx.Input)
	}
	c.Data["count"] = count
	c.Data["page"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["pagecount"] = pageCount
	c.Data["userList"] = list
	c.Data["user"] = user
}

//获取网贷记账信息
//@router /loanaccount [get]
func (c *UserMetadataController) LoanAccount() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/loan_account.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}

	list, err := models.GetLoanAccountList(uid, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户网贷记账信息失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetLoanAccountListCount(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户网贷记账信息总数失败", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PageSize10)
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}

	c.Data["count"] = count
	c.Data["page"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["pagecount"] = pageCount
	c.Data["userList"] = list
	c.Data["user"] = user
}

//获取社保信息
//@router /userssecuritydata [get]
func (c *UserMetadataController) UsersSecurityData() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/social_security.html"
	}()
	uid, _ := c.GetInt("id")
	//获取用户社保信息
	var createTime time.Time
	var mxSbData models.MXSBData
	if utils.RunMode == "release" { //从MongoDB获取
		var mdbMxSBData models.MdbMxSBData
		//从MongoDB获取社保信息
		session := utils.GetSession()
		defer session.Close()
		credit := session.DB(utils.MgoDbName).C("mxsecuritydata")
		uidMap := make(map[string]int)
		uidMap["uid"] = uid
		err := credit.Find(uidMap).One(&mdbMxSBData)
		if err != nil && err.Error() != "not found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "从数据库获取用户社保信息出错", err.Error(), c.Ctx.Input)
			c.Data["msg"] = "获取用户社保信息出错"
			c.Data["name"] = "社保"
			c.TplName = ""
			return
		}
		if err != nil && err.Error() == "not found" {
			c.Data["msg"] = "没有该用户社保数据"
			c.Data["name"] = "社保"
			c.TplName = ""
			return
		}
		createTime = mdbMxSBData.CreateTime
		mxSbData = mdbMxSBData.MXSBData
	} else { //从mysql获取
		data, create, err := models.GetUsersMXSBGJJInfo(uid, "security")
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "从数据库获取用户社保信息出错", err.Error(), c.Ctx.Input)
			c.Data["msg"] = "获取用户社保信息出错"
			c.Data["name"] = "社保"
			c.TplName = ""
			return
		}
		if err != nil && err.Error() == utils.ErrNoRow() {
			c.Data["msg"] = "没有该用户社保数据"
			c.Data["name"] = "社保"
			c.TplName = ""
			return
		}
		err = json.Unmarshal([]byte(data), &mxSbData)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析用户社保信息出错", err.Error(), c.Ctx.Input)
			c.Data["msg"] = "解析用户社保信息出错"
			c.Data["name"] = "社保"
			c.TplName = ""
			return
		}
		createTime = create
	}
	var usersecuritydata models.UserSecurityData
	usersecuritydata.City = mxSbData.City                                    //所属城市
	usersecuritydata.AreaCode = mxSbData.AreaCode                            //地区编码
	usersecuritydata.SocialSecurityNo = mxSbData.BaseInfo.SocialSecurityNo   //社保卡号
	usersecuritydata.PersonalNo = mxSbData.BaseInfo.PersonalNo               //个人编号
	usersecuritydata.RealName = mxSbData.BaseInfo.RealName                   //姓名
	usersecuritydata.BaseNumber = strconv.Itoa(mxSbData.BaseInfo.BaseNumber) //缴存基数
	usersecuritydata.LastPayDate = mxSbData.BaseInfo.LastPayDate             //最后缴存日期
	switch mxSbData.BaseInfo.IDType { //证件类型
	case "ID_CARD", "":
		usersecuritydata.IdType = "身份证"
	case "PASSPORT":
		usersecuritydata.IdType = "护照"
	}
	usersecuritydata.IdCard = mxSbData.BaseInfo.IDCard                               //证件号码
	usersecuritydata.Nation = mxSbData.BaseInfo.Nation                               //民族
	usersecuritydata.BirthDay = mxSbData.BaseInfo.BirthDay                           //生日
	usersecuritydata.Phone = mxSbData.BaseInfo.Phone                                 //手机号码
	usersecuritydata.Address = mxSbData.BaseInfo.Address                             //家庭住址
	usersecuritydata.PersonnelStatus = mxSbData.BaseInfo.PersonnelStatus             //人员状态
	usersecuritydata.HouseholdRegistration = mxSbData.BaseInfo.HouseholdRegistration //户口属性
	usersecuritydata.FirstInsuredDate = mxSbData.BaseInfo.FirstInsuredDate           //首次参保时间
	usersecuritydata.WorkTime = mxSbData.BaseInfo.WorkTime                           //参加工作时间
	usersecuritydata.InsuredUnit = mxSbData.BaseInfo.InsuredUnit                     //参保单位
	usersecuritydata.InsuredUnitCode = mxSbData.BaseInfo.InsuredUnitCode             //参保单位编码
	usersecuritydata.UnitType = mxSbData.BaseInfo.UnitType                           //单位类型
	switch mxSbData.BaseInfo.PayStatus { //缴存状态
	case "NONE":
		usersecuritydata.PayStatus = "未缴纳"
	case "NORMAL":
		usersecuritydata.PayStatus = "正常"
	case "SUSPENSE":
		usersecuritydata.PayStatus = "停缴"
	case "CLOSED":
		usersecuritydata.PayStatus = "注销"
	}
	usersecuritydata.BeginDate = mxSbData.BaseInfo.BeginDate //开户日期
	for _, v := range mxSbData.Insurances {
		switch v.InsuranceType {
		case "工伤保险":
			usersecuritydata.IndustrialInsuranceBaseNumber = strconv.Itoa(v.BaseNumber) //工伤保险基数
		case "失业保险":
			usersecuritydata.UnemploymentInsuranceBaseNumber = strconv.Itoa(v.BaseNumber) //失业保险基数
		case "医疗保险":
			usersecuritydata.MedicalInsuranceBaseNumber = strconv.Itoa(v.BaseNumber) //医疗保险基数
		case "养老保险":
			usersecuritydata.EndowmentInsuranceBaseNumber = strconv.Itoa(v.BaseNumber) //养老保险基数
		case "生育保险":
			usersecuritydata.MaternityInsuranceBaseNumber = strconv.Itoa(v.BaseNumber) //生育保险基数
		}
	}
	usersecuritydata.FetchTime = mxSbData.BaseInfo.FetchTime                                           //抓取时间
	usersecuritydata.MedicalInsuranceBalance = strconv.Itoa(mxSbData.BaseInfo.MedicalInsuranceBalance) //医疗保险账户余额
	//获取用户信息
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["user"] = user
	c.Data["securityinfo"] = usersecuritydata
	c.Data["createTime"] = createTime.Format("2006-01")
}

//获取公积金信息
//@router /usersfunddata [get]
func (c *UserMetadataController) UsersFundData() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/accumulation_fund.html"
	}()
	uid, _ := c.GetInt("id")
	var mxGJJData models.MXGJJData
	var createTime time.Time
	//获取用户公积金信息
	if utils.RunMode == "release" { //从MongoDB获取
		//从MongoDB获取公积金信息
		var mdbMxGJJData models.MdbMxGJJData
		session := utils.GetSession()
		defer session.Close()
		credit := session.DB(utils.MgoDbName).C("mxfunddata")
		uidMap := make(map[string]int)
		uidMap["uid"] = uid
		err := credit.Find(uidMap).One(&mdbMxGJJData)
		if err != nil && err.Error() != "not found" {
			cache.RecordLogs(uid, 0, c.User.Name, c.User.DisplayName, "", "从数据库获取用户公积金信息出错", err.Error(), c.Ctx.Input)
			c.Data["msg"] = "获取用户公积金信息出错"
			c.Data["name"] = "公积金"
			return
		}
		if err != nil && err.Error() == "not found" {
			c.Data["msg"] = "没有该用户公积金数据"
			c.Data["name"] = "公积金"
			return
		}
		mxGJJData = mdbMxGJJData.MXGJJData
		createTime = mdbMxGJJData.CreateTime
	} else { //从mysql获取
		data, create, err := models.GetUsersMXSBGJJInfo(uid, "fund")
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(uid, 0, c.User.Name, c.User.DisplayName, "", "从数据库获取用户公积金信息出错", err.Error(), c.Ctx.Input)
			c.Data["msg"] = "获取用户公积金信息出错"
			c.Data["name"] = "公积金"
			return
		}
		if err != nil && err.Error() == utils.ErrNoRow() {
			c.Data["msg"] = "没有该用户公积金数据"
			c.Data["name"] = "公积金"
			return
		}
		err = json.Unmarshal([]byte(data), &mxGJJData)
		if err != nil {
			cache.RecordLogs(uid, 0, c.User.Name, c.User.DisplayName, "", "解析用户社保信息出错", err.Error(), c.Ctx.Input)
			c.Data["msg"] = "解析用户公积金信息出错"
			c.Data["name"] = "公积金"
			return
		}
		createTime = create
	}
	var userFundData models.YDGJJUserInfo
	userFundData.RealName = mxGJJData.UserInfo.RealName //姓名
	switch mxGJJData.UserInfo.Gender { //性别
	case "0":
		userFundData.Gender = "女"
	case "1":
		userFundData.Gender = "男"
	}
	userFundData.Birthday = mxGJJData.UserInfo.Birthday                           //出生年月
	userFundData.Mobile = mxGJJData.UserInfo.Mobile                               //手机号
	userFundData.Email = mxGJJData.UserInfo.Email                                 //邮箱
	userFundData.CustomerNumber = mxGJJData.UserInfo.CustomerNumber               //客户号
	userFundData.GjjNumber = mxGJJData.UserInfo.GjjNumber                         //公积金账号
	userFundData.Balance = strconv.Itoa(mxGJJData.UserInfo.Balance)               //账号余额
	userFundData.FundBalance = mxGJJData.UserInfo.FundBalance                     //公积金余额
	userFundData.SubsidyBalance = strconv.Itoa(mxGJJData.UserInfo.SubsidyBalance) //补贴公积金账户余额
	userFundData.SubsidyIncome = strconv.Itoa(mxGJJData.UserInfo.SubsidyIncome)   //补贴月缴存
	switch mxGJJData.UserInfo.PayStatus { //缴存状态
	case "NONE":
		userFundData.PayStatus = "未缴纳"
	case "NORMAL":
		userFundData.PayStatus = "正常"
	case "SUSPENSE":
		userFundData.PayStatus = "停缴"
	case "CLOSED":
		userFundData.PayStatus = "注销"
	}
	switch mxGJJData.UserInfo.CardType { //证件类型
	case "ID_CARD", "":
		userFundData.CardType = "身份证"
	case "PASSPORT":
		userFundData.CardType = "护照"
	}
	userFundData.IdCard = mxGJJData.UserInfo.IDCard                                                   //证件号
	userFundData.HomeAddress = mxGJJData.UserInfo.HomeAddress                                         //通讯地址
	userFundData.CorporationNumber = mxGJJData.UserInfo.CorporationNumber                             //企业账户号码
	userFundData.CorporationName = mxGJJData.UserInfo.CorporationName                                 //当前缴存企业名称
	userFundData.MonthlyCorporationIncome = strconv.Itoa(mxGJJData.UserInfo.MonthlyCorporationIncome) //企业月度缴存
	userFundData.MonthlyCustomerIncome = strconv.Itoa(mxGJJData.UserInfo.MonthlyCustomerIncome)       //个人月度缴存
	userFundData.MonthlyTotalIncome = strconv.Itoa(mxGJJData.UserInfo.MonthlyTotalIncome)             //月度总缴存
	userFundData.CorporationRatio = mxGJJData.UserInfo.CorporationRatio                               //企业缴存比例
	userFundData.CustomerRatio = mxGJJData.UserInfo.CustomerRatio                                     //个人缴存比例
	userFundData.SubsidyCorporationRatio = mxGJJData.UserInfo.SubsidyCorporationRatio                 //补贴公积金公司缴存比例
	userFundData.SubsidyCustomerRatio = mxGJJData.UserInfo.SubsidyCustomerRatio                       //补贴公积金个人缴存比例
	userFundData.BaseNumber = strconv.Itoa(mxGJJData.UserInfo.BaseNumber)                             //缴存基数
	userFundData.LastPayDate = mxGJJData.UserInfo.LastPayDate                                         //最新缴存日期
	userFundData.BeginDate = mxGJJData.UserInfo.BeginDate                                             //开户日期
	//获取用户信息
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(uid, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["user"] = user
	c.Data["fundinfo"] = userFundData
	c.Data["createTime"] = createTime.Format("2006-01")
}

//获取用户钱包收支记录
//@router /userswalletdata [get]
func (c *UserMetadataController) UsersWalletData() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/purse_purse.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	//获取用户钱包收支记录
	records, err := models.GetUserWalletRecords(uid, utils.StartIndex(page, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户钱包收支记录出错", err.Error(), c.Ctx.Input)
	}
	//整理返回数据
	var userWalletRecords []models.UserWalletRecord
	for _, v := range records {
		var userWalletRecord models.UserWalletRecord
		userWalletRecord.CreateTime = v.CreateTime.Format("2006-01-02 15:04")
		userWalletRecord.AfterMoneyAmount = strconv.FormatFloat(v.AfterMoneyAmount, 'f', 2, 64)
		switch v.PayOrGet {
		case 1:
			userWalletRecord.PayOrGet = "消耗"
			userWalletRecord.MoneyAmount = strconv.FormatFloat(0-v.MoneyAmount, 'f', 2, 64)
		case 2:
			userWalletRecord.PayOrGet = "收入"
			userWalletRecord.MoneyAmount = "+" + strconv.FormatFloat(v.MoneyAmount, 'f', 2, 64)
		}
		switch v.DealType {
		case 1:
			userWalletRecord.DealType = "话费1元券"
		case 2:
			userWalletRecord.DealType = "获取新口子"
		case 3:
			userWalletRecord.DealType = "网贷征信查询"
		case 4:
			userWalletRecord.DealType = "话费10元券"
		case 5:
			userWalletRecord.DealType = "平台征信查询"
		case 6:
			userWalletRecord.DealType = "1个月会员"
		case 7:
			userWalletRecord.DealType = "2个月会员"
		case 8:
			userWalletRecord.DealType = "购买抽奖"
		case 9:
			userWalletRecord.DealType = "任务奖励"
		case 10:
			userWalletRecord.DealType = "签到奖励"
		case 11:
			userWalletRecord.DealType = "充值"
		case 12:
			userWalletRecord.DealType = "提现"
		case 13:
			userWalletRecord.DealType = "抽奖奖励"
		case 14:
			userWalletRecord.DealType = "贷款稳下"
		}
		userWalletRecords = append(userWalletRecords, userWalletRecord)
	}
	//获取记录条数
	count, err := models.GetUserWalletRecordsCount(uid)
	//获取页数
	pagecount, err := utils.GetPageCount(count, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取页数失败", err.Error(), c.Ctx.Input)
	}
	//获取用户信息
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["messcount"] = count
	c.Data["messpage"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["messpagecount"] = pagecount
	c.Data["walletrecords"] = userWalletRecords
	c.Data["user"] = user
}

//获取用户积分兑换记录
//@router /userexchangerecord [get]
func (c *UserMetadataController) UserExchangeRecord() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/integration_drawingRecords.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	//获取积分兑换记录
	scoreExchangeRecords, err := models.GetUserScoreExchangeRecords(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户积分消费记录失败", err.Error(), c.Ctx.Input)
	}
	var userExchangeRecords []models.ScoreExchangeRecord
	for _, v := range scoreExchangeRecords {
		var userExchangeRecord models.ScoreExchangeRecord
		switch v.DealType {
		case 1:
			userExchangeRecord.Prize = "话费1元券"
		case 2:
			userExchangeRecord.Prize = "获取新口子"
		case 3:
			userExchangeRecord.Prize = "网贷征信查询"
		case 4:
			userExchangeRecord.Prize = "话费10元券"
		case 5:
			userExchangeRecord.Prize = "平台征信查询"
		case 6:
			userExchangeRecord.Prize = "1个月会员"
		case 7:
			userExchangeRecord.Prize = "2个月会员"
		case 8:
			userExchangeRecord.Prize = "购买抽奖"
		}
		userExchangeRecord.CreateTime = v.CreateTime.Format("2006-01-02 15:04")
		userExchangeRecord.ScoreAmount = strconv.Itoa(v.ScoreAmount)
		userExchangeRecord.AfterScoreAmount = strconv.Itoa(v.AfterScoreAmount)
		userExchangeRecords = append(userExchangeRecords, userExchangeRecord)
	}
	count := len(userExchangeRecords)
	//获取页数
	pagecount, err := utils.GetPageCount(count, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取页数失败", err.Error(), c.Ctx.Input)
	}
	if pagecount > 0 {
		start := utils.PageSize10 * (page - 1)
		end := utils.PageSize10 * page
		if count > start && start < end {
			if count >= end {
				userExchangeRecords = userExchangeRecords[start:end]
			} else {
				userExchangeRecords = userExchangeRecords[start:]
			}
		}
	}
	//获取用户信息
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["messcount"] = count
	c.Data["messpage"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["messpagecount"] = pagecount
	c.Data["scoreexchange"] = userExchangeRecords
	c.Data["user"] = user
	return
}

//获取用户积分抽奖记录
//@router /userlotteryrecord [get]
func (c *UserMetadataController) UserLotteryRecord() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "user-management/user_lottery_records.html"
	}()
	uid, _ := c.GetInt("id")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	//获取用户积分抽奖记录
	lotteryRecords, err := models.GetUserLotteryRecords(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户积分抽奖记录失败", err.Error(), c.Ctx.Input)
	}
	var scoreLotteryRecords []models.ScoreLotteryRecord
	for _, v := range lotteryRecords {
		var scoreLotteryRecord models.ScoreLotteryRecord
		scoreLotteryRecord.CreateTime = v.CreateTime.Format("2006-01-02 15:04")
		scoreLotteryRecord.Content = v.Content
		scoreLotteryRecords = append(scoreLotteryRecords, scoreLotteryRecord)
	}
	//降序排列
	count := len(scoreLotteryRecords)
	//获取页数
	pagecount, err := utils.GetPageCount(count, utils.PageSize10)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取页数失败", err.Error(), c.Ctx.Input)
	}
	if pagecount > 0 {
		start := utils.PageSize10 * (page - 1)
		end := utils.PageSize10 * page
		if count > start && start < end {
			if count >= end {
				scoreLotteryRecords = scoreLotteryRecords[start:end]
			} else {
				scoreLotteryRecords = scoreLotteryRecords[start:]
			}
		}
	}
	//获取用户信息
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["messcount"] = count
	c.Data["messpage"] = page
	c.Data["pagesize"] = utils.PageSize10
	c.Data["messpagecount"] = pagecount
	c.Data["scorelottery"] = scoreLotteryRecords
	c.Data["user"] = user
	return
}

//获取用户个人资信
//@router /usercredit [get]
func (c *UserMetadataController) UserCredit() {
	c.IsNeedTemplate()
	uid, _ := c.GetInt("id")
	userCredit, err := models.GetUserCredit(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户个人资信失败", err.Error(), c.Ctx.Input)
	}
	//获取用户信息
	user, err := models.GetUserById(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询用户失败", err.Error(), c.Ctx.Input)
	}
	c.Data["user"] = user
	c.Data["userCredit"] = userCredit
	c.TplName = "user-management/personal_credit.html"
}
