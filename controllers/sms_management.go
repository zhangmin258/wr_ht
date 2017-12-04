package controllers

import (
	"strconv"
	"strings"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"

	//"net/url"
	//"net/http"
	"time"
)

/*
推送信息接口
*/
type SmsManagementController struct {
	BaseController
}

//获取推送历史信息列表
//@router /getsmsmanagelist [get]
func (c *SmsManagementController) GetSMSManageList() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//短信内容
	if content := c.GetString("content"); content != "" {
		condition += " AND content LIKE ?"
		params = append(params, "%"+content+"%")
	}
	//当前用户id
	uid := c.User.Id
	//查询
	smsMange, err := models.GetSMSMange(condition, params, uid, utils.StartIndex(pageNum, utils.PAGE_SIZE), utils.PAGE_SIZE)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取历史短信信息失败！", err.Error(), c.Ctx.Input)
	}
	for i := 0; i < len(smsMange); i++ {
		if smsMange[i].Operator != "" {
			smsMange[i].Nature += "运营商：" + smsMange[i].Operator + " & "
		}
		if smsMange[i].Address != "" {
			smsMange[i].Nature += "地域要求：" + smsMange[i].Address + " & "
		}
		if smsMange[i].App == 1 {
			smsMange[i].Nature += "操作系统：ios & "
		}
		if smsMange[i].App == 2 {
			smsMange[i].Nature += "操作系统：Android & "
		}
		if smsMange[i].App == 3 {
			smsMange[i].Nature += "操作系统：wp & "
		}
		if smsMange[i].MinZmscore != 0 && smsMange[i].MaxZmscore != 0 {
			smsMange[i].Nature += "芝麻分：" + strconv.Itoa(smsMange[i].MinZmscore) + "~" + strconv.Itoa(smsMange[i].MaxZmscore) + " & "
		}
		if smsMange[i].LoanDemand != "" {
			smsMange[i].Nature += "贷款额度：" + smsMange[i].LoanDemand + "(元) & "
		}
		if smsMange[i].LoginTime == 1 {
			smsMange[i].Nature += "3天未登录的用户 & "
		}
		if smsMange[i].LoginTime == 2 {
			smsMange[i].Nature += "7天未登录的用户 & "
		}
		if smsMange[i].LoginTime == 3 {
			smsMange[i].Nature += "3天内登录过的用户 & "
		}
		if smsMange[i].Remark != "" {
			smsMange[i].Nature += "用户标签：" + smsMange[i].Remark + " & "
		}
		if len(smsMange[i].Nature) > 2 {
			smsMange[i].Nature = smsMange[i].Nature[:len(smsMange[i].Nature)-2]
		}
	}
	//总数
	count, err := models.GetSMSManageCount(uid, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取历史短信数量失败！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败！", err.Error(), c.Ctx.Input)
	}
	SMSCount, err := models.GetPushSMSCount(c.User.Id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取短信发送总条数失败！", err.Error(), c.Ctx.Input)
	}
	c.Data["SMSCount"] = SMSCount
	c.Data["smsMange"] = smsMange
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "sms-management/message_history.html"
}

//推送信息页面
//@router /getsmspage [get]
func (c *SmsManagementController) GetSMSPage() {
	c.IsNeedTemplate()
	SMSCount, err := utils.SMSCount("weirong", "weirong", "0")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取剩余短信条数失败！", err.Error(), c.Ctx.Input)
	}
	//空间畅想获取余额
	name := utils.KJCXNAME                           //账号
	seed := time.Now().Format("20060102150405")  //当前时间
	password := utils.PASSWORD                   //密码
	key := utils.MD5(utils.MD5(password) + seed) //md5(md5(password)+seed
	SMSBalance, err := utils.SMSKJCXBalance(name, seed, key)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取空间畅想余额失败！", err.Error(), c.Ctx.Input)
	}
	c.Data["SMSBalance"] = SMSBalance
	c.Data["SMSCount"] = SMSCount
	c.TplName = "sms-management/sms_management.html"
}

//发送推送信息
//@router /pushsms [post]
func (c *SmsManagementController) PushSMS() {
	defer c.ServeJSON()
	var smsManagement models.SMSManagement
	err := c.ParseForm(&smsManagement)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "smsManagement参数解析异常", err.Error(), c.Ctx.Input)
		//解析参数异常
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "smsManagement参数解析异常!"}
		return
	}
	if smsManagement.PushTime == "" {
		smsManagement.PushTime = time.Now().Format("2006-01-02 15:04:05")
	}
	//当前用户
	smsManagement.SysUserId = c.User.Id
	condition := ""
	params := []interface{}{}
	//登录时间筛选
	switch smsManagement.LoginTime {
	case 1: //3天未登录
		condition += " AND DateDiff(now(),u.login_time)>=?"
		params = append(params, 3)
	case 2: //7天未登录
		condition += " AND DateDiff(now(),u.login_time)>=?"
		params = append(params, 7)
	case 3: //3天内登录过
		condition += " AND DateDiff(now(),u.login_time)<=?"
		params = append(params, 3)
	}
	//使用手机设备筛选
	if smsManagement.App != 0 {
		condition += " AND u.app = ? "
		params = append(params, smsManagement.App)
	}
	//芝麻信用筛选
	if smsManagement.MinZmscore != 0 && smsManagement.MaxZmscore != 0 {
		condition += " AND um.zm_score>= ? AND um.zm_score<=? "
		params = append(params, smsManagement.MinZmscore)
		params = append(params, smsManagement.MaxZmscore)
	}
	smsManagement.Address = strings.Replace(smsManagement.Address, "，", ",", -1)
	smsManagement.Operator = strings.Replace(smsManagement.Operator, "，", ",", -1)
	smsManagement.LoanDemand = strings.Replace(smsManagement.LoanDemand, "，", ",", -1)
	//addrs := strings.Split(smsManagement.Address, ",")
	operator := strings.Split(smsManagement.Operator, ",")
	loanDemand := strings.Split(smsManagement.LoanDemand, ",")
	//地址筛选<--暂时关闭该条件-->
	/*if smsManagement.Address != "" {
		condition += " AND ( u.address LIKE ? "
		for k, v := range addrs {
			if k != len(addrs)-1 {
				condition += " OR u.address LIKE ? "
			} else {
				condition += ") "
			}
			params = append(params, "%"+v+"%")
		}
	}*/
	//运营商筛选
	if smsManagement.Operator != "" {
		condition += " AND ( 1>1 "
		for k, v := range operator {
			switch v {
			case "中国移动":
				condition += ` OR  TRIM(u.account) REGEXP "^1([3][4-9]|[4][7]|[5][0-27-9]|[8][2-478])[0-9]{8}$" `
			case "中国联通":
				condition += ` OR TRIM(u.account) REGEXP "^1([3][0-2]|[4][5]|[5][5-6]|[7][6]|[8][5-6])[0-9]{8}$" `
			case "中国电信":
				condition += ` OR TRIM(u.account) REGEXP "^1(3[3]|5[3]|7[37]|8[019])[0-9]{8}$" `
			case "其他":
				condition += ` OR( TRIM(u.account) NOT REGEXP "^1(3[3]|5[3]|7[37]|8[019])[0-9]{8}$"
				AND TRIM(u.account) NOT REGEXP "^1([3][4-9]|[4][7]|[5][0-27-9]|[8][2-478])[0-9]{8}$"
				 AND  TRIM(u.account) NOT REGEXP "^1([3][0-2]|[4][5]|[5][5-6]|[7][6]|[8][5-6])[0-9]{8}$" )`
			}
			if k == len(operator)-1 {
				condition += ") "
			}
		}
	}
	//贷款金额筛选
	if smsManagement.LoanDemand != "" {
		condition += " AND ( 1>1 "
		var l int
		for k, v := range loanDemand {
			if strings.Contains(v, "以内") {
				condition += " OR  ub.loan_amount<=? "
				params = append(params, 1000)
				if k == len(loanDemand)-1 {
					condition += ") "
				}
			} else if strings.Contains(v, "以上") {
				condition += " OR  ub.loan_amount>=? "
				params = append(params, 300000)
				if k == len(loanDemand)-1 {
					condition += ") "
				}
			} else {
				condition += " OR (ub.loan_amount>=? AND ub.loan_amount<=?) "
				if k == len(loanDemand)-1 {
					condition += ") "
				}
				loan := strings.Split(v, "~")
				l, _ = strconv.Atoi(loan[0])
				params = append(params, l)
				l, _ = strconv.Atoi(loan[1])
				params = append(params, l)
			}
		}
	}
	//获取满足条件的所有手机号码
	phones, err := models.GetPushSMSUserAccount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取手机号码失败", err.Error(), c.Ctx.Input)
	}

	//发送的条数
	smsManagement.PushCount = len(phones)
	count := len(phones)
	pageSize := 1
	pageCount, _ := utils.GetPageCount(count, pageSize)
	if pageCount == 1 {
		//手机号数组转换字符串
		phone := strings.Join(phones, ",")
		go utils.SendLocalSMSPostWeiRongMany(phone, smsManagement.Content, "weirong", c.Ctx.Input.IP(), "0")
	} else {
		start := 0
		end := 0
		for i := 1; i <= pageCount; i++ {
			start = (i - 1) * pageSize
			end = start + pageSize
			var phoneList = phones[start:end]
			phone := strings.Join(phoneList, ",")
			go utils.SendLocalSMSPostWeiRongMany(phone, smsManagement.Content, "weirong", c.Ctx.Input.IP(), "0")
		}
	}
	err = models.SaveSMSManage(&smsManagement)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存smsManagement异常", err.Error(), c.Ctx.Input)
		//解析参数异常
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "保存smsManagement异常!"}
		return
	}

	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功!"}
}

//空间畅想发送信息
//@router /kjcxpushsms [post]
func (c *SmsManagementController) KJCXSendSMS() {
	defer c.ServeJSON()
	var smsManagement models.SMSManagement
	err := c.ParseForm(&smsManagement)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "smsManagement参数解析异常", err.Error(), c.Ctx.Input)
		//解析参数异常
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "smsManagement参数解析异常!"}
		return
	}
	if smsManagement.PushTime == "" {
		smsManagement.PushTime = time.Now().Format("2006-01-02 15:04:05")
	}
	//当前用户
	smsManagement.SysUserId = c.User.Id
	condition := ""
	params := []interface{}{}
	//登录时间筛选
	switch smsManagement.LoginTime {
	case 1: //3天未登录
		condition += " AND DateDiff(now(),u.login_time)>=?"
		params = append(params, 3)
	case 2: //7天未登录
		condition += " AND DateDiff(now(),u.login_time)>=?"
		params = append(params, 7)
	case 3: //3天内登录过
		condition += " AND DateDiff(now(),u.login_time)<=?"
		params = append(params, 3)
	}
	//使用手机设备筛选
	if smsManagement.App != 0 {
		condition += " AND u.app = ? "
		params = append(params, smsManagement.App)
	}
	//芝麻信用筛选
	if smsManagement.MinZmscore != 0 && smsManagement.MaxZmscore != 0 {
		condition += " AND um.zm_score>= ? AND um.zm_score<=? "
		params = append(params, smsManagement.MinZmscore)
		params = append(params, smsManagement.MaxZmscore)
	}
	smsManagement.Address = strings.Replace(smsManagement.Address, "，", ",", -1)
	smsManagement.Operator = strings.Replace(smsManagement.Operator, "，", ",", -1)
	smsManagement.LoanDemand = strings.Replace(smsManagement.LoanDemand, "，", ",", -1)
	//addrs := strings.Split(smsManagement.Address, ",")
	operator := strings.Split(smsManagement.Operator, ",")
	loanDemand := strings.Split(smsManagement.LoanDemand, ",")
	//地址筛选<--暂时关闭该条件-->
	/*if smsManagement.Address != "" {
		condition += " AND ( u.address LIKE ? "
		for k, v := range addrs {
			if k != len(addrs)-1 {
				condition += " OR u.address LIKE ? "
			} else {
				condition += ") "
			}
			params = append(params, "%"+v+"%")
		}
	}*/
	//运营商筛选
	if smsManagement.Operator != "" {
		condition += " AND ( 1>1 "
		for k, v := range operator {
			switch v {
			case "中国移动":
				condition += ` OR  TRIM(u.account) REGEXP "^1([3][4-9]|[4][7]|[5][0-27-9]|[8][2-478])[0-9]{8}$" `
			case "中国联通":
				condition += ` OR TRIM(u.account) REGEXP "^1([3][0-2]|[4][5]|[5][5-6]|[7][6]|[8][5-6])[0-9]{8}$" `
			case "中国电信":
				condition += ` OR TRIM(u.account) REGEXP "^1(3[3]|5[3]|7[37]|8[019])[0-9]{8}$" `
			case "其他":
				condition += ` OR( TRIM(u.account) NOT REGEXP "^1(3[3]|5[3]|7[37]|8[019])[0-9]{8}$"
				AND TRIM(u.account) NOT REGEXP "^1([3][4-9]|[4][7]|[5][0-27-9]|[8][2-478])[0-9]{8}$"
				 AND  TRIM(u.account) NOT REGEXP "^1([3][0-2]|[4][5]|[5][5-6]|[7][6]|[8][5-6])[0-9]{8}$" )`
			}
			if k == len(operator)-1 {
				condition += ") "
			}
		}
	}
	//贷款金额筛选
	if smsManagement.LoanDemand != "" {
		condition += " AND ( 1>1 "
		var l int
		for k, v := range loanDemand {
			if strings.Contains(v, "以内") {
				condition += " OR  ub.loan_amount<=? "
				params = append(params, 1000)
				if k == len(loanDemand)-1 {
					condition += ") "
				}
			} else if strings.Contains(v, "以上") {
				condition += " OR  ub.loan_amount>=? "
				params = append(params, 300000)
				if k == len(loanDemand)-1 {
					condition += ") "
				}
			} else {
				condition += " OR (ub.loan_amount>=? AND ub.loan_amount<=?) "
				if k == len(loanDemand)-1 {
					condition += ") "
				}
				loan := strings.Split(v, "~")
				l, _ = strconv.Atoi(loan[0])
				params = append(params, l)
				l, _ = strconv.Atoi(loan[1])
				params = append(params, l)
			}
		}
	}
	//获取满足条件的所有手机号码
	dests, err := models.GetPushSMSUserAccount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取手机号码失败", err.Error(), c.Ctx.Input)
	}
	//发送的条数
	smsManagement.PushCount = len(dests)
	name := utils.KJCXNAME                           //账号
	seed := time.Now().Format("20060102150405")  //当前时间
	password := utils.PASSWORD                   //密码
	key := utils.MD5(utils.MD5(password) + seed) //md5(md5(password)+seed
	count := len(dests)
	pageSize := 1
	pageCount, _ := utils.GetPageCount(count, pageSize)
	var returnmsg utils.SmsResult
	if pageCount == 1 {
		//手机号数组转换字符串
		dest := strings.Join(dests, ",")
		returnmsg, _ = utils.SMSKJCXSend(name, seed, key, dest, smsManagement.Content, "", "", "")

	} else {
		start := 0
		end := 0
		for i := 1; i <= pageCount; i++ {
			start = (i - 1) * pageSize
			end = start + pageSize
			var destList = dests[start:end]
			dest := strings.Join(destList, ",")
			returnmsg, _ = utils.SMSKJCXSend(name, seed, key, dest, smsManagement.Content, "", "", "")

		}
	}
	err = models.SaveSMSManage(&smsManagement)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存smsManagement异常", err.Error(), c.Ctx.Input)
		//解析参数异常
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "保存smsManagement异常!"}
		return
	}

	c.Data["json"] = map[string]interface{}{"ret": returnmsg.Ret, "msg": returnmsg.Msg}

}
