package controllers

import (
	"encoding/json"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/services"
	"wr_v1/utils"
	"zcm_tools/email"
	"zcm_tools/log"

	"github.com/astaxie/beego"
)

type NotifyController struct {
	beego.Controller
}

var callbackLog *log.Log

func init() {
	callbackLog = log.Init("20060102.callback")
}

// 提现回调
//@router /lltrade [post]
func (c NotifyController) LlTrade() {
	defer c.ServeJSON()
	var m models.RechargeFeedback
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	if err != nil {
		cache.RecordLogs(0, 0, "", "", "services/LoanResultUpdate", "提现回调解析参数异常", err.Error(), nil)
		email.Send("提现回调解析参数异常!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
		return
	}
	// 修改订单状态
	go services.LoanResultUpdate(m.No_order, "0000", m.Result_pay, 0)
	if m.Result_pay == "FAILURE" || m.Result_pay == "CANCEL" || m.Result_pay == "SUCCESS" {
		c.Data["json"] = map[string]string{"ret_code": "0000", "ret_msg": "交易成功"}
	}
}
