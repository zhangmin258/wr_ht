package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/crypt"
	"zcm_tools/email"
	"zcm_tools/http"
	"zcm_tools/pay"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

//银行卡卡bin查询接口
func GetLLBankCardBin(cardnumber string) (*models.LLBandCardBinResponse, error) {
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("sign_type", "RSA")
	param.Add("card_no", cardnumber)
	b, err := pay.DoRequest(param, utils.LL_BANDCARD_BIN, utils.LL_RSA_PRIVATE_KEY, pay.ContentTypeJson, pay.EncryptMD5withRsa)
	if err != nil {
		return nil, err
	}
	var m *models.LLBandCardBinResponse
	err = json.Unmarshal(b, &m)
	if m == nil {
		return nil, errors.New("解析连连返回数据出错")
	}
	return m, nil
}

//连连API签约申请
func LLSignApply(param http.Values) (*models.LLTradeApiResponse, error) {
	b, err := pay.DoRequest(param, utils.LL_SignApply, utils.LL_RSA_PRIVATE_KEY, pay.ContentTypeJson, pay.EncryptMD5withRsa)
	if err != nil {
		return nil, err
	}
	var m *models.LLTradeApiResponse
	err = json.Unmarshal(b, &m)
	if m == nil {
		return nil, errors.New("解析连连签约申请返回数据出错")
	}
	return m, nil
}

// 获取连连API签约申请token
func CheckLLSignToken(key string) string {
	r_token, _ := utils.Rc.RedisString(key)
	if r_token == "" {
		return ""
	} else {
		return r_token
	}
}

//连连API签约验证
func LLSignVerfy(param http.Values) (*models.LLTradeResponse, error) {
	b, err := pay.DoRequest(param, utils.LL_SignVerfy, utils.LL_RSA_PRIVATE_KEY, pay.ContentTypeJson, pay.EncryptMD5withRsa)
	if err != nil {
		return nil, err
	}
	var m *models.LLTradeResponse
	err = json.Unmarshal(b, &m)
	if m == nil {
		return nil, errors.New("解析连连签约申请返回数据出错")
	}
	return m, nil
}

// 发送提现请求
func LLTrade(order_number, card_no, user_name, money_order string) (res *models.LLTradeResponse, err error) {
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("api_version", "1.0")
	param.Add("sign_type", "RSA")
	param.Add("no_order", order_number)
	param.Add("dt_order", time.Now().Format("20060102150405")) //time.Now().Format("20060102150405")
	if utils.RunMode != "release" {
		money_order = "0.01"
	}
	param.Add("money_order", money_order)
	param.Add("card_no", card_no)
	param.Add("acct_name", user_name)
	param.Add("info_order", "微融")
	param.Add("flag_card", "0")
	param.Add("notify_url", utils.WR_API_URL+utils.LL_Trade_CallBack)
	sign, err := pay.Sign(param, utils.LL_RSA_PRIVATE_KEY, pay.EncryptMD5withRsa)
	fmt.Println("sign:", sign)
	if err != nil {
		beego.Debug(err.Error())
	}
	param.Add("sign", sign)
	fmt.Println(string(param.Data()))
	pay_load := crypt.LLEncrypt(param.Data(), []byte(utils.LL_PUBLIC_KEY))
	body, err := http.Post(utils.LL_Trade_API, `{"pay_load":"`+pay_load+`","oid_partner":"`+utils.LL_oid_partner+`"}`, "application/json;charset=UTF-8")
	if err != nil {
		return
	}
	json.Unmarshal(body, &res)
	if res.Ret_code == "4002" || res.Ret_code == "4003" || res.Ret_code == "4004" {
		// 发送确认消息
		res, err = LLTradeConfirm(order_number, res.Confirm_code)
	}
	return
}

//连连放款(确认)
func LLTradeConfirm(order_number, confirm_code string) (res *models.LLTradeResponse, err error) {
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("api_version", "1.0")
	param.Add("sign_type", "RSA")
	param.Add("no_order", order_number)
	param.Add("confirm_code", confirm_code)
	param.Add("notify_url", utils.WR_API_URL+utils.LL_Trade_CallBack)
	sign, err := pay.Sign(param, utils.LL_RSA_PRIVATE_KEY, pay.EncryptMD5withRsa)
	param.Add("sign", sign)
	pay_load := crypt.LLEncrypt(param.Data(), []byte(utils.LL_PUBLIC_KEY))
	body, err := http.Post(utils.LL_Trade_CONFIRM, `{"pay_load":"`+pay_load+`","oid_partner":"`+utils.LL_oid_partner+`"}`, "application/json;charset=UTF-8")
	if err != nil {
		return
	}
	json.Unmarshal(body, &res)
	fmt.Println(res.Ret_code)
	return
}

//连连放款订单查询
func LLTradeQuery(order_number string) (*models.LLTradeQueryResponse, error) {
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("sign_type", "RSA")
	param.Add("api_version", "1.0")
	param.Add("no_order", order_number)
	body, err := pay.DoRequest(param, utils.LL_Trade_Query, utils.LL_RSA_PRIVATE_KEY, pay.ContentTypeJson, pay.EncryptMD5withRsa)
	if err != nil {
		return nil, err
	}
	var m *models.LLTradeQueryResponse
	err = json.Unmarshal(body, &m)
	if m == nil {
		return nil, errors.New("解析连连返回数据出错")
	}
	return m, nil
}

// 连连充值订单主动查询
func CostActiveQuery(orderCode string, sysId int) {
	var count = 1
	for true {

		time.Sleep(time.Second * 15)
		// beego.Info("主动查询次数:", count)
		res, err := LLTradeQuery(orderCode) // 发起查询
		if err != nil {
			cache.RecordLogs(0, 0, "", "", "services/CostActiveQuery", "连连充值订单主动查询", err.Error(), nil)
			email.Send("连连充值订单主动查询!", err.Error(), utils.ToUsers, "services/CostActiveQuery")
			return
		}
		fmt.Println(res)
		// 修改订单状态
		go LoanResultUpdate(orderCode, res.Ret_code, res.Result_pay, sysId)
		if res.Ret_code != "0000" || res.Result_pay == "FAILURE" || res.Result_pay == "SUCCESS" || res.Result_pay == "CANCEL" { // 得到确切结果后结束循环
			break
		}
		count++
		if count > 5 {
			// TODO 异常订单记录到数据库
			// beego.Info("异常订单:", orderCode, res.Ret_code, res.Result_pay)
			err := models.SetIllegalOrder(orderCode)
			if err != nil {
				cache.RecordLogs(0, 0, "", "", "services/CostActiveQuery", "异常订单记录到数据库操作失败 ", err.Error(), nil)
				email.Send("异常订单记录到数据库操作失败!", err.Error(), utils.ToUsers, "services/CostActiveQuery")
			}
			break
		}
	}
}

// 提现结果处理
func LoanResultUpdate(orderNumber, retCode, resultPay string, sysId int) {

	if utils.Rc.Lock(utils.CACHE_KEY_Loan_Accept_State_Deal+orderNumber, time.Minute) {
		defer utils.Rc.Delete(utils.CACHE_KEY_Loan_Accept_State_Deal + orderNumber)
		var state int
		var paylog int
		// 从mysql中获取提现订单状态
		orderStatus, err := models.GetUserWithdrawDepositStatus(orderNumber)
		if err != nil {
			cache.RecordLogs(sysId, 0, "", "", "services/LoanResultUpdate", "mysql中获取提现订单状态发生异常", err.Error(), nil)
			email.Send("mysql中获取提现订单状态发生异常!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
			return
		}
		// 订单状态已处理
		if orderStatus == "SUCCESS" || orderStatus == "FAILURE" || orderStatus == "CANCEL" {
			return
		}
		// 查询订单金额和用户id
		id, money, err := models.GetUidByOrderNumber(orderNumber)
		if err != nil {
			cache.RecordLogs(sysId, 0, "", "", "services/LoanResultUpdate", "查询订单金额和用户id", err.Error(), nil)
			email.Send("查询订单金额和用户id!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
			return
		}
		// 查询用户钱包余额
		balance, err := models.GetWalletBalance(id)
		if err != nil {
			cache.RecordLogs(sysId, 0, "", "", "services/LoanResultUpdate", "查询用户钱包余额", err.Error(), nil)
			email.Send("查询用户钱包余额!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
			return
		}
		o := orm.NewOrm()
		o.Begin()
		// 订单作失败处理
		if retCode != "0000" && retCode != "9999" && retCode != "4006" && retCode != "4007" && retCode != "4009" && retCode != "1002" && retCode != "2005" {
			resultPay = "FAILURE"
		}
		if resultPay == "SUCCESS" {
			state = 3
			paylog = 0
		} else if resultPay == "FAILURE" || resultPay == "CANCEL" || resultPay == "CLOSED" || resultPay == "CHECK" {
			state = 4
			paylog = 1
			// 提现失败给用户钱包补款
			err = models.WalletRechargeTransaction(id, money, o)
			if err != nil {
				cache.RecordLogs(sysId, 0, "", "", "services/LoanResultUpdate", "提现失败给用户钱包补款", err.Error(), nil)
				email.Send("提现失败给用户钱包补款!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
				o.Rollback()
				return
			}

		} else if resultPay == "PROCESSING" || resultPay == "APPLY" {
			state = 1
		}
		// 修改收支记录
		err = models.UpdateRechargeFinance(balance+20, paylog, orderNumber, o)
		if err != nil {
			cache.RecordLogs(sysId, 0, "", "", "services/LoanResultUpdate", "修改收支记录", err.Error(), nil)
			email.Send("修改收支记录!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
			o.Rollback()
			return
		}
		// 修改提现记录
		err = models.ModifyWithdrawDepositType(resultPay, retCode, orderNumber, state, o)
		if err != nil {
			cache.RecordLogs(sysId, 0, "", "", "services/LoanResultUpdate", "修改提现记录", err.Error(), nil)
			email.Send("修改提现记录!", err.Error(), utils.ToUsers, "services/LoanResultUpdate")
			o.Rollback()
			return
		}
		o.Commit()
	}

}
