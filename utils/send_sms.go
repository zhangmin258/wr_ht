package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"zcm_tools/http"
)

type SmsResult struct {
	Ret int
	Msg string
}

//云融正通获取剩余短信条数
type SmsBalance struct {
	ResultCode        string //返回值状态：0成功；非0错误
	ErrorCode         string //错误码 -1用户名或口令错误；-2为IP验证错误
	Balance           string //余额，元为单位
	UnitPrice         string //发送单价，元为单位
	MessageQtyBalance string //剩余条数
}

//云融正通发送短信
type YRZTSend struct {
	MessageId   string //提交消息
	MobilePhone string //手机号码
	ResultCode  string //处理结果
	ErrorCode   string //错误码
}

//发送多条微融短信
func SendLocalSMSPostWeiRongMany(account, content, merchant, ip, category string) bool {
	params := url.Values{}
	params.Set("account", account)
	params.Set("content", content)
	params.Set("merchant", merchant)
	params.Set("ip", ip)
	params.Set("category", category)
	b, err := http.Post(SMSURLWEIRONGMANY, params.Encode())
	if err != nil {
		fmt.Println("2", err.Error())
		return false
	} else {
		var m SmsResult
		if err := json.Unmarshal(b, &m); err == nil && m.Ret == 200 {
			fmt.Println(true)
			return true
		} else {
			return false
		}
	}
	return true
}

//获取剩余短信条数
func SMSCount(merchant, channel, category string) (count SmsResult, err error) {
	countByte, err := http.Get(SMSURLWEIRONGCOUNT + "?merchant=" + merchant + "&channel=" + channel + "&category=" + category)
	if err != nil {
		return
	}
	err = json.Unmarshal(countByte, &count)
	return
}

//	空间畅想发送短信
func SMSKJCXSend(name, seed, key, dest, content, ext, reference, delay string) (statemsg SmsResult, err error) {
	params := url.Values{}
	params.Set("name", name)           //帐号
	params.Set("seed", seed)           //当前时间
	params.Set("key", key)             //md5( md5(password)  +  seed) )
	params.Set("dest", dest)           //手机号码
	params.Set("content", content)     //短信内容
	params.Set("ext", ext)             //扩展号码:nullable
	params.Set("reference", reference) //参考信息:nullable
	params.Set("delay", delay)         //定时参数 :nullable
	s, err := http.Post(SMSURLKJCXSEND, params.Encode())
	if err != nil {
		return
	}
	statemsg.Msg = string(s)
	smsByteArr := strings.Split(statemsg.Msg, ":")
	statemsg.Msg = smsByteArr[0]

	if statemsg.Msg == "success" {
		statemsg.Ret = 200
		statemsg.Msg = string(smsByteArr[1])
	} else {
		statemsg.Msg = smsByteArr[0] + ":" + string(smsByteArr[1])
		statemsg.Ret = 0
	}
	return
}

//空间畅想查询余额
func SMSKJCXBalance(name, seed, key string) (balancemsg SmsResult, err error) {
	fmt.Println(SMSURLKJCXBALANCE + "?name=" + name + "&seed=" + seed + "&key=" + key)
	msgByte, err := http.Get(SMSURLKJCXBALANCE + "?name=" + name + "&seed=" + seed + "&key=" + key)
	fmt.Println(string(msgByte))
	if err != nil {
		return
	}
	smsByteArr := strings.Split(string(msgByte), ":")
	balancemsg.Msg = smsByteArr[0]
	ret, err := strconv.Atoi(smsByteArr[1])
	if balancemsg.Msg == "success" {
		balancemsg.Ret = ret
	} else {
		balancemsg.Ret = -1
	}
	return
}

//云融正通查询剩余短信条数
func YRZTSMSCount(cmd, userName, passWord string) (result SmsBalance, err error) {
	url:=SMSURLYRZTMESSAGE + "?cmd=" + cmd + "&userName=" + userName + "&passWord=" + passWord
	fmt.Println(url)
	countByte, err := http.Get(url)
	fmt.Println(string(countByte))
	if err != nil {
		return
	}
	//fmt.Println(string(countByte))
	err = json.Unmarshal(countByte, &result)
	return
}

//云融正通批量发送相同内容不同手机号码短信
func SendLocalSMSPostYRZT(cmd, userName, passWord, mobilePhone, body, scheduleDateStr, messageId string) (result YRZTSend, err error) {
	params := url.Values{}
	params.Set("cmd", cmd)
	params.Set("userName", userName)
	params.Set("passWord", passWord)
	params.Set("mobilePhone", mobilePhone)
	params.Set("body", body)
	if scheduleDateStr != "" {
		params.Set("scheduleDateStr", scheduleDateStr)
	} else {
		params.Set("scheduleDateStr", "")
	}
	params.Set("messageId", messageId)
	str, err := http.Post(SMSURLYRZTMESSAGE, params.Encode())
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(str, &result)
	if err == nil && result.ResultCode == "0" {
		return result, nil
	} else {
		return result, err
	}
}
