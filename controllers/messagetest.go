package controllers

import (
	"strings"
	"time"
	"wr_v1/models"
	"wr_v1/utils"

	"wr_v1/cache"
)

type MsgtestController struct {
	BaseController
}

//APP运营-短信管理-短信测试-测试历史
//@router  /testhistory [get]
func (this *MsgtestController) TestHistory() {
	this.IsNeedTemplate()
	//查询分页信息
	pageNum, _ := this.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	if msgContent := this.GetString("send_text"); msgContent != "" {
		condition += " and content LIKE ?"
		params = append(params, "%"+msgContent+"%")
	}
	msgHistoryList, err := models.FindMsgHistory(condition, params, utils.StartIndex(pageNum, utils.PageSize10), utils.PageSize10)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询短信测试历史异常！", err.Error(), this.Ctx.Input)
	}
	count, err := models.GetHistoryCount(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询测试历史总条数异常！", err.Error(), this.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PageSize10)

	this.Data["msgHistoryList"] = msgHistoryList
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.TplName = "sms-management/sms_test.html"
}

//APP运营-短信管理-短信测试-发送测试
//@router /sendmessage [post]
func (this *MsgtestController) SendMessage() {
	defer this.ServeJSON()
	var messagetest models.MessageTest
	err := this.ParseForm(&messagetest)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "短信参数解析异常", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "短信参数解析异常!"}
		return
	}
	messagetest.UserId = this.User.Id
	messagetest.Phone = strings.Replace(messagetest.Phone, "，", ",", -1)
	var count int
	if messagetest.Phone != "" {
		count = strings.Count(messagetest.Phone, ",") + 1
	} else {
		count = strings.Count(messagetest.Phone, ",")
	}
	if utils.SendLocalSMSPostWeiRongMany(messagetest.Phone, messagetest.Message, "weirong", this.Ctx.Input.IP(), "0") {
		err = models.AddHistoryMessage(messagetest.Message, messagetest.UserId, messagetest.Phone, count)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "插入短信测试历史异常！", err.Error(), this.Ctx.Input)
			//插入申请条件异常
			this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "插入数据异常"}
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "插入数据成功"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送短信成功"}
	return
}

//APP运营-短信管理-短信测试-发送测试空间畅想
//@router /kjcxsendmessage [post]
func (this *MsgtestController) KjcxSendMessage() {
	defer this.ServeJSON()
	var messagetest models.MessageTest
	err := this.ParseForm(&messagetest)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "短信参数解析异常", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "短信参数解析异常!"}
		return
	}
	messagetest.UserId = this.User.Id
	messagetest.Phone = strings.Replace(messagetest.Phone, "，", ",", -1)
	name := utils.KJCXNAME                       //账号
	seed := time.Now().Format("20060102150405")  //当前时间
	password := utils.PASSWORD                   //密码
	key := utils.MD5(utils.MD5(password) + seed) //md5(md5(password)+seed
	var count int
	if messagetest.Phone != "" {
		count = strings.Count(messagetest.Phone, ",") + 1
	} else {
		count = strings.Count(messagetest.Phone, ",")
	}
	var returnmsg utils.SmsResult
	returnmsg, _ = utils.SMSKJCXSend(name, seed, key, messagetest.Phone, messagetest.Message, "", "", "")
	err = models.AddHistoryMessage(messagetest.Message, messagetest.UserId, messagetest.Phone, count)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "插入短信测试历史异常！", err.Error(), this.Ctx.Input)
		//插入申请条件异常
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "插入数据异常"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "插入数据成功"}
	this.Data["json"] = map[string]interface{}{"ret": returnmsg.Ret, "msg": returnmsg.Msg}
	return
}
