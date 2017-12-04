package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//测试历史
type MsgHistory struct {
	PushDate  string //发送时间
	Content   string //推送内容
	Phone     string //测试号码
	PushCount int    //发送条数
}

//分页信息和搜索内容
type PageAndSearch struct {
	Page    int    //分页信息
	Message string //发送文本
}

type MessageTest struct {
	UserId  int    //创建人id
	Message string //短信文本
	Phone   string //测试手机号
}

//分页查询短信测试历史
func FindMsgHistory(condition string, params []string, begin int, size int) (msghistory []MsgHistory, err error) {
	sql := `SELECT DATE_FORMAT(push_time,'%Y-%m-%d %H:%m')push_date,content,phone,push_count FROM sms_test WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY push_date DESC limit ?,?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&msghistory)
	return
}

//查询测试历史总条数
func GetHistoryCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM sms_test WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//插入短信测试历史
func AddHistoryMessage(message string, uid int, phone string, count int) error {
	sql := `INSERT INTO sms_test(content,sys_user_id,phone,push_time,push_count,create_time)VALUES(?,?,?,?,?,?)`
	o := orm.NewOrm()
	_, err := o.Raw(sql, message, uid, phone, time.Now(), count, time.Now()).Exec()
	return err
}
