package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/file"

	"github.com/astaxie/beego/toolbox"
	"github.com/tealeg/xlsx"
)

type MsgMarketingController struct {
	BaseController
}

/*
	短信营销接口
*/

//短信营销页面
//@router /getmsgmarketing [get]
func (c *MsgMarketingController) GetMsgMarketing() {
	c.IsNeedTemplate()
	count, err := models.GetCountOfUsers()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取用户总数异常", err.Error(), c.Ctx.Input)
	}
	c.Data["count"] = count
	c.TplName = "sms-management/sms_marketing.html"
}

//发送短信接口
//@router /sendmarketingmsg [post]
func (c *MsgMarketingController) SendMarketingMsg() {
	defer c.ServeJSON()
	sysUserId := c.User.Id
	if !utils.Rc.SetNX(utils.CACHE_KEY_ChECKWITHSENDMSG, 1, time.Second*30) {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "短信已经提交，请不要重复发送哦！"}
		return
	}
	var sendmsg models.SendMsg
	err := c.ParseForm(&sendmsg)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "msgMarketing参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "msgMarketing参数解析异常!"}
		return
	}
	var pushCount int
	var mobilePhones string
	var mobilePhone []string

	//上传的文件名
	fileName := c.GetString("fileName")
	if fileName != "" {
		fileNames := strings.Split(fileName, "\\")
		sendmsg.AccountSource += fileNames[len(fileNames)-1]
	}

	//上传的文件
	filesPath := c.GetString("filesPath")
	if filesPath != "" {
		count1 := 0
		if strings.Contains(filesPath, "text") {
			phoneTxt, err := ReceiveTxt(filesPath)
			if err != nil {
				c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "导入文件失败！"}
				return
			}
			if phoneTxt != "" {
				count1, phoneTxt = DecreaseInMass(strings.Split(phoneTxt, "\n"))
			}
			if errPhone, ok := CheckPhone(phoneTxt); !ok {
				c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "文件中有手机号不正确！phone:" + errPhone}
				return
			}
			pushCount += count1
			mobilePhones += phoneTxt
		} else {
			phones, err := ReceiveExcel(filesPath)
			if err != nil {
				c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "导入文件失败！"}
				return
			}
			count1, phoneExcel := DecreaseInMass(phones)
			if errPhone, ok := CheckPhone(phoneExcel); !ok {
				c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "文件中有手机号不正确！phone:" + errPhone}
				return
			}
			pushCount += count1
			mobilePhones += phoneExcel
		}

	}

	//微融用户
	if sendmsg.Begin != 0 && sendmsg.End != 0 && sendmsg.Begin <= sendmsg.End {
		sendmsg.AccountSource += "   微融：" + strconv.Itoa(sendmsg.Begin) + "到" + strconv.Itoa(sendmsg.End) + "位的用户"
		phones, err := models.GetPhonesByUid(sendmsg.Begin-1, sendmsg.End-sendmsg.Begin+1)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据用户编号获取手机号异常", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "根据用户编号获取手机号异常!"}
			return
		}
		count2, weirongPhone := DecreaseInMass(phones)
		if errPhone, ok := CheckPhone(weirongPhone); !ok {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "微融用户手机号不正确！phone:" + errPhone}
			return
		}
		pushCount += count2
		if mobilePhones != "" {
			mobilePhones = mobilePhones + "\n" + weirongPhone
		} else {
			mobilePhones += weirongPhone
		}
	}

	if sendmsg.ChannelName != 1 && sendmsg.ChannelName != 2 && sendmsg.ChannelName != 3 {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请选择要发送的短信通道"}
		return
	}
	var source, newSource []string
	content := sendmsg.Body
	var url string
	reg := regexp.MustCompile(`[\P{Han}]+`)
	a := reg.FindAllString(content, -1)
	for _, v := range a {
		if strings.Contains(v, "http://") {
			v = strings.Replace(v, ",", "", -1)
			v = strings.Replace(v, "，", "", -1)
			v = strings.TrimSpace(v)
			url = v
		}
	}
	//解析短链接
	newUlr, err := utils.SortToLong(utils.ACCESSTOKEN, url)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析短链接异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "解析短链接异常!"}
		return
	}
	if strings.Contains(newUlr, "source") {
		source = strings.SplitAfter(newUlr, "source=")
		newSource = strings.Split(source[1], "")
		for _, v := range newSource {
			if v == " " || len(v) > 1 {
				break
			}
			sendmsg.Source += v
		}
	} else if strings.Contains(newUlr, "parms") {
		source = strings.SplitAfter(newUlr, "parms=")
		newSource = strings.Split(source[1], "")
		for _, v := range newSource {
			if v == " " || len(v) > 1 {
				break
			}
			sendmsg.Source += v
		}
	}
	//用于区分多条短信循环发送唯一标识
	sendmsg.Flag = time.Now().Format("20060102150405")
	if sendmsg.MobilePhones != "" {
		p := strings.Split(sendmsg.MobilePhones, "\n")
		count3, mobilePhones1 := DecreaseInMass(p)
		if errPhone, ok := CheckPhone(mobilePhones1); !ok {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "输入的手机号不正确！phone:" + errPhone}
			return
		}
		if mobilePhones != "" {
			mobilePhones = mobilePhones + "\n" + mobilePhones1
		} else {
			mobilePhones += mobilePhones1
		}
		pushCount += count3
		//pushCount = strings.Count(mobilePhones, "\n") + 1
	}
	/*mobilePhone = strings.Split(mobilePhones, "\n")
	pushCount, mobilePhones = DecreaseInMass(mobilePhone)*/
	mobilePhone = strings.Split(mobilePhones, "\n")
	var sendCount int
	if i := pushCount % utils.KONGJIANSENDSMSCOUNT; i != 0 {
		sendCount = pushCount/utils.KONGJIANSENDSMSCOUNT + 1
	} else {
		sendCount = pushCount / utils.KONGJIANSENDSMSCOUNT
	}
	if pushCount == 0 {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "手机号为空！"}
		return
	}
	var m []string
	go func() {
		var count int
		for j := 0; j < sendCount; j++ {
			if j == sendCount-1 {
				m = mobilePhone[j*utils.KONGJIANSENDSMSCOUNT:]
				count = len(mobilePhone) - j*utils.KONGJIANSENDSMSCOUNT
			} else {
				m = mobilePhone[j*utils.KONGJIANSENDSMSCOUNT : (j+1)*utils.KONGJIANSENDSMSCOUNT]
				count = utils.KONGJIANSENDSMSCOUNT
			}
			mobilePhones = strings.Join(m, ",")
			if sendmsg.ChannelName == 1 { //当渠道码为0的时候发送微融短信
				// msgResult := utils.SendLocalSMSPostWeiRongMany(mobilePhones, sendmsg.Body, "weirong", c.Ctx.Input.IP(), "0")
				// if msgResult {
				// 	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功"}
				// 	return
				// } else {
				// 	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "发送微融短信异常！", err.Error(), c.Ctx.Input)
				// 	c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "发送微融短信异常!"}
				// 	return
				// }
				c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功"}
				return
			} else if sendmsg.ChannelName == 2 { //当渠道码为1的时候发送空间畅想短信
				name := utils.KJCXNAME                       //账号
				seed := time.Now().Format("20060102150405")  //当前时间
				password := utils.PASSWORD                   //密码
				key := utils.MD5(utils.MD5(password) + seed) //md5(md5(password)+seed
				if sendmsg.PushTime != "" {
					TimeToPush, err := time.Parse(utils.FormatDateTime, sendmsg.PushTime+":00")
					if err != nil {
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "时间解析失败！", err.Error(), c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "时间解析失败！"}
						return
					}
					if !TimeToPush.Before(time.Now()) {
						var spec = utils.TimeToTaskSpec(TimeToPush) //格式化时间
						if spec == "" {
							c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "请输入正确的推送时间！"}
							return
						}
						//定时任务
						taskMission := toolbox.NewTask("DelayPush", spec, func() error {
							statemsg, err := utils.SMSKJCXSend(name, seed, key, mobilePhones, sendmsg.Body, "", "", "")
							if err != nil {
								cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "空间畅想发送短信异常", err.Error(), c.Ctx.Input)
							}
							err = models.AddContentBySysUserId(statemsg.Msg, sysUserId, count, sendmsg)
							if err != nil {
								cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "存储短信信息失败", err.Error(), c.Ctx.Input)
							}
							return nil
						})
						toolbox.StopTask()
						toolbox.AddTask("DelayPush", taskMission)
						toolbox.StartTask()
						c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "设置定时成功"}
						return
					} else {
						c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请输入正确的推送时间！"}
						return
					}
				} else {
					statemsg, err := utils.SMSKJCXSend(name, seed, key, mobilePhones, sendmsg.Body, "", "", "")
					if err != nil {
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "发送空间畅想短信异常！", err.Error(), c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "发送空间畅想短信异常!"}
						return
					}
					err = models.AddContentBySysUserId(statemsg.Msg, sysUserId, count, sendmsg)
					if err != nil {
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "存储短信信息失败", err.Error(), c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "存储短信信息失败!"}
						return
					}
					c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功"}
				}
			} else if sendmsg.ChannelName == 3 { //当渠道码为2的时候发送云融正通的短信
				cmd := "sendBatchMessage"
				userName := utils.YRZTNAME
				passWord := utils.YRZTPASSWORD
				var messageId string
				if sendmsg.PushTime != "" {
					TimeToPush, err := time.Parse(utils.FormatDateTime, sendmsg.PushTime+":00")
					if err != nil {
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "时间解析失败！", err.Error(), c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "时间解析失败！"}
						return
					}
					if !TimeToPush.Before(time.Now()) {
						messageId = strconv.Itoa(sysUserId) + "MSG" + TimeToPush.Format("20060102150405")
						msgResult, err := utils.SendLocalSMSPostYRZT(cmd, userName, passWord, mobilePhones, sendmsg.Body, TimeToPush.Format("20060102150405"), messageId)
						if err != nil {
							cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "发送云融正通短信异常！错误码："+msgResult.ErrorCode, err.Error(), c.Ctx.Input)
							c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "发送云融正通短信异常!错误码：" + msgResult.ErrorCode}
							return
						} else if msgResult.ResultCode != "0" {
							cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "发送云融正通短信异常！错误码："+msgResult.ErrorCode, "", c.Ctx.Input)
							c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "发送云融正通短信异常!错误码：" + msgResult.ErrorCode}
							return
						}
						err = models.AddContentBySysUserId(messageId, sysUserId, count, sendmsg)
						if err != nil {
							cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "存储短信信息失败", err.Error(), c.Ctx.Input)
							c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "存储短信信息失败!"}
							return
						}
						c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功"}
					} else {
						c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请输入正确的推送时间！"}
						return
					}
				} else {
					messageId = strconv.Itoa(sysUserId) + "MSG" + time.Now().Format("20060102150405")
					msgResult, err := utils.SendLocalSMSPostYRZT(cmd, userName, passWord, mobilePhones, sendmsg.Body, "", messageId)
					if err != nil {
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "发送云融正通短信异常！错误码："+msgResult.ErrorCode, err.Error(), c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "发送云融正通短信异常!错误码：" + msgResult.ErrorCode}
						return
					} else if msgResult.ResultCode != "0" {
						if msgResult.ResultCode == "-10" {
							cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "云融正通余额不足", "", c.Ctx.Input)
							c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "云融正通余额不足，请充值！" + msgResult.ErrorCode}
							return
						}
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "发送云融正通短信异常！错误码："+msgResult.ErrorCode, "", c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "发送云融正通短信异常!错误码：" + msgResult.ErrorCode}
						return
					}
					err = models.AddContentBySysUserId(messageId, sysUserId, count, sendmsg)
					if err != nil {
						cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "存储短信信息失败", err.Error(), c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "存储短信信息失败!"}
						return
					}
					c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功"}
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "发送成功。"}
	return
}

//短信通道剩余条数
//@router /showmsgbalance [post]
func (c *MsgMarketingController) ShowMsgBalance() {
	defer c.ServeJSON()
	channelName, err := c.GetInt("ChannelName")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取通道名称失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取通道名称失败!"}
		return
	}
	var SMSCount int
	var Balance float64
	if channelName == 1 {

		SMSBalance, err := utils.SMSCount("weirong", "weirong", "0")
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取剩余短信条数失败！", err.Error(), c.Ctx.Input)
		}
		SMSCount = SMSBalance.Ret
	} else if channelName == 2 {
		name := utils.KJCXNAME                       //账号
		seed := time.Now().Format("20060102150405")  //当前时间
		password := utils.PASSWORD                   //密码
		key := utils.MD5(utils.MD5(password) + seed) //md5(md5(password)+seed
		SMSBalance, err := utils.SMSKJCXBalance(name, seed, key)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取空间畅想剩余短信条数失败！", err.Error(), c.Ctx.Input)
		}
		SMSCount = SMSBalance.Ret
		Balance = float64(SMSBalance.Ret) * 0.05
	} else if channelName == 3 {
		SMSBalance, err := utils.YRZTSMSCount("getBalance", utils.YRZTNAME, utils.YRZTPASSWORD)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取云融正通剩余短信条数失败！", err.Error(), c.Ctx.Input)
		}
		SMSCount, _ = strconv.Atoi(SMSBalance.MessageQtyBalance)
		Balance, _ = strconv.ParseFloat(SMSBalance.Balance, 64)
		//c.Data["Count"] = Count       //短信发送条数
		//c.Data["SMSCount"] = SMSCount //短信剩余条数
		if SMSBalance.ResultCode != "0" || SMSBalance.ErrorCode != "" {
			return
		}

	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "SMSCount": 0}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "SMSCount": SMSCount, "Balance": fmt.Sprintf("%.2f", Balance)}
	return
}

//接收文件并回填至表单
//@router /receivefile [post]
func (c *MsgMarketingController) ReceiveFile() {
	defer c.ServeJSON()
	var (
		count int
		err   error
	)
	//接收txt文件
	var phones []string
	filesPath := c.GetString("filesPath")
	if strings.Contains(filesPath, "text") {
		phone, err := ReceiveTxt(filesPath)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "导入文件失败！"}
			return
		}
		if phone != "" {
			count = strings.Count(phone, "\n") + 1
		}
	} else {
		phones, err = ReceiveExcel(filesPath)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "导入文件失败！"}
			return
		}
		count = len(phones)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "文件导入成功！", "Count": count}
	return
}

//接收txt文件
func ReceiveTxt(txtUrl string) (mobilePhones string, err error) {
	if txtUrl != "" {
		txtUrl = strings.Replace(txtUrl, "data:text/plain;base64,", "", -1)
		txtUrl = strings.Replace(txtUrl, "data:;base64,", "", -1)
		arr := strings.Split(txtUrl, "||")
		length := len(arr[0])
		filePath := ""
		deleteFile := []string{}
		if length > 0 {
			timeStr := strconv.FormatInt(time.Now().Unix(), 10)
			txtFileUrl := "txtData" + timeStr + ".txt"
			filePath = "./static/" + txtFileUrl
			err := file.SaveBase64ToFile(arr[0], filePath)
			if err != nil {
				return "", err
			}
			deleteFile = append(deleteFile, txtFileUrl)
		}
		txtFile, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		txtStr, err := ioutil.ReadAll(txtFile)
		if err != nil {
			return "", err
		}
		txtFile.Close()
		if string(txtStr) != "" {
			mobilePhones += string(txtStr)
		}
		if len(deleteFile) > 0 {
			os.Remove("./static/" + deleteFile[0])
		}
	}
	return
}

//接收excel文件
func ReceiveExcel(xlsxUrl string) (mobilePhones []string, err error) {
	xlsxUrl = strings.Replace(xlsxUrl, "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,", "", -1)
	xlsxUrl = strings.Replace(xlsxUrl, "data:;base64,", "", -1)
	arr := strings.Split(xlsxUrl, "||")
	length := len(arr[0])
	xlFileUrl := ""
	deleteFile := []string{}
	if length > 0 {
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		fileName := "xlData" + timeStr + ".xlsx"
		filePath := "./static/" + fileName
		err := file.SaveBase64ToFile(arr[0], filePath)
		if err != nil {
			return nil, err
		}
		deleteFile = append(deleteFile, fileName)
		xlFileUrl = filePath
	}
	defer func() {
		if len(deleteFile) > 0 {
			for i := 0; i < len(deleteFile); i++ {
				os.Remove("./static/" + deleteFile[i])
			}
		}
	}()
	xlFile, err := xlsx.OpenFile(xlFileUrl)
	if err != nil {
		return nil, err
	}

	var phones []string
	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows {
		for k, cell := range row.Cells {
			if k == 0 {
				phones = append(phones, cell.Value)
			}
		}
	}
	return phones, nil
}

//根据用户编号查询手机号并回填至表单
//@router /getphonesbyuid [post]
func (c *MsgMarketingController) GetPhonesByUid() {
	defer c.ServeJSON()
	startUid, err := c.GetInt("StartUid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取起始用户编号失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取起始用户编号失败!"}
		return
	}
	endUid, err := c.GetInt("EndUid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取结束用户编号失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取结束用户编号失败!"}
		return
	}

	//根据用户编号查询手机号
	if startUid <= endUid {
		phones, err := models.GetPhonesByUid(startUid-1, endUid-startUid+1)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据用户编号获取手机号异常", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "根据用户编号获取手机号异常!"}
			return
		}
		count, _ := DecreaseInMass(phones)
		c.Data["json"] = map[string]interface{}{"ret": 200, "Count": count, "msg": "用户导入成功！"}
	}
	return
}

//手机号码去重
func DecreaseInMass(phones []string) (count int, phone string) {
	phoneNumbersMap := make(map[string]string)
	var mobilePhones string
	for _, v := range phones {
		if _, ok := phoneNumbersMap[strings.TrimSpace(v)]; !ok {
			phoneNumbersMap[strings.TrimSpace(v)] = v
			mobilePhones += "\n" + strings.TrimSpace(v)
			count++
		}
	}
	mobilePhones = strings.TrimSpace(mobilePhones)
	return count, mobilePhones
}

//手机号判断
func CheckPhone(phones string) (phones1 string, bool bool) {
	bool = true
	regular := "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\\\d{8}$"
	reg := regexp.MustCompile(regular)
	phones = strings.TrimSpace(phones)
	mphone := strings.Split(phones, "\n")
	for _, v := range mphone {
		if len(v) != 11 || reg.MatchString(v) {
			phones1 = v
			bool = false
			return
		}
	}
	return
}
