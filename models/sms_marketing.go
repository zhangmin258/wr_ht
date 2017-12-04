package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type SendMsg struct {
	MobilePhones  string //手机号码
	Body          string //短信内容
	ChannelName   int    //通道名称
	PushTime      string //定时发送时间
	Begin         int    //起始用户编号
	End           int    //终止用户编号
	Source        string //来源
	Flag          string //用于区分多条短信循环发送唯一标识
	AccountSource string //手机号来源
}

//根据用户编号查询用户手机号码
func GetPhonesByUid(startUid, endUid int) (phones []string, err error) {
	sql := `SELECT account FROM users ORDER BY id LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, startUid, endUid).QueryRows(&phones)
	return
}

//发送短信前存储所发送的短信
func AddContentBySysUserId(smsId string, sysUserId, pushCount int, sendMsg SendMsg) error {
	var pushTime string
	if sendMsg.PushTime != "" {
		pushTime = sendMsg.PushTime
	} else {
		pushTime = time.Now().Format("2006-01-02 15:04:05")
	}
	sql := `INSERT INTO sms_management(sms_id,content,sys_user_id,push_count,push_time,create_time,begin,end,plateform,source,flag,account_source)VALUES(?,?,?,?,?,NOW(),?,?,?,?,?,?)`
	_, err := orm.NewOrm().Raw(sql, smsId, sendMsg.Body, sysUserId, pushCount, pushTime, sendMsg.Begin, sendMsg.End, sendMsg.ChannelName, sendMsg.Source, sendMsg.Flag, sendMsg.AccountSource).Exec()
	return err
}

//查询用户数量
func GetCountOfUsers() (count int, err error) {
	sql := `SELECT COUNT(1) FROM users `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}
