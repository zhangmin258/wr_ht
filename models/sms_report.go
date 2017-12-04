package models

import (
	"github.com/astaxie/beego/orm"
)

//接收状态报告
type GetReport struct {
	Name   string //账号
	Report string //状态报告
}

//解析状态报告
type MsgReport struct {
	SmsId        string //消息编号
	PhoneNumber  string //手机号码
	StatusReport string //状态报告
	Reference    string //参考信息
	PushTime     string //推送时间
}

//云融正通接收状态报告
type GetYRZTReport struct {
	UserName             string //用户名
	PassWord             string //密码
	MessageQty           string //状态报告条数
	SubmitMessageId      string //发送的消息ID
	ClientMessageBatchId string //客户端批次ID
	MobilePhone          string //手机号码
	DateTimeStr          string //接收时间
	DeliveryStatus       string //状态报告标识
	DeliveryStatusCode   string //状态报告
}

//空间畅想状态报告存储
func AddReportData(name string, msgReport []MsgReport) error {
	sql := `INSERT INTO sms_report(name,sms_id,phone_number,status_report,reference,push_time,create_time) VALUES(?,?,?,?,?,?,NOW())`
	o := orm.NewOrm()
	reportsql, err := o.Raw(sql).Prepare()
	for _, v := range msgReport {
		_, err = reportsql.Exec(name, v.SmsId, v.PhoneNumber, v.StatusReport, v.Reference, v.PushTime)
	}
	reportsql.Close()
	return err
}

//云融正通状态报告存储
func AddYRZTReportData(getYRZTReport []GetYRZTReport) error {
	sql := `INSERT INTO sms_report(name,password,sms_id,client_message_batch_id,phone_number,delivery_status_code,status_report,push_time,create_time)VALUES(?,?,?,?,?,?,?,?,NOW())`
	o := orm.NewOrm()
	reportsql, err := o.Raw(sql).Prepare()
	for _, v := range getYRZTReport {
		_, err = reportsql.Exec(v.UserName, v.PassWord, v.SubmitMessageId, v.ClientMessageBatchId, v.MobilePhone, v.DeliveryStatusCode, v.DeliveryStatus, v.DateTimeStr)
	}
	reportsql.Close()
	return err
}
