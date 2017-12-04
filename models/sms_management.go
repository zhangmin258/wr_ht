package models

import (
	"github.com/astaxie/beego/orm"

	"time"
)

type SMSManagement struct {
	Id         int    //主键
	Content    string //短信文本
	SysUserId  int    //创建人id
	Address    string //地域要求，多个地名之间以逗号隔开
	MinZmscore int    //最小芝麻分要求
	MaxZmscore int    //最大芝麻分要求
	App        int8   //1 ios,2 android,3 wp
	LoginTime  int    //多久未登录的时间
	Operator   string //运营商要求,多个之间用逗号隔开
	Remark     string //用户标签
	LoanDemand string //贷款需求
	PushCount  int    //发送条数
	PushTime   string //推送时间
	Nature     string //用户属性
}

type UserMessage struct {
	Id         int       //
	Title      string    //标题
	Content    string    //内容
	MsgType    int64     //该消息的类型：1.新口子推荐。2.审核未通过。3.审核通过。4.未放款成功。5.跳转url 。6注册成功信息
	IsRead     int64     //是否已读0未读1已读
	CreateTime time.Time //创建时间
}

//分页查询历史短信信息
func GetSMSMange(condition string, params []string, uid, begin, size int) (smsManage []SMSManagement, err error) {
	sql := `SELECT id,content,push_count,address,
	min_zmscore,max_zmscore,
	app,login_time,push_time,
	remark,operator,
	loan_demand FROM sms_management WHERE sys_user_id = ?`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY push_time DESC LIMIT ?,?"
	_, err = orm.NewOrm().Raw(sql, uid, params, begin, size).QueryRows(&smsManage)
	return
}

// 查询所有历史短信数量
func GetSMSManageCount(uid int, condition string, params []string) (count int, err error) {
	sql := `SELECT count(1)
	FROM sms_management
	WHERE sys_user_id = ? `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, uid, params).QueryRow(&count)
	return
}

//保存短信
func SaveSMSManage(smsManagement *SMSManagement) error {
	sql := `INSERT INTO sms_management (content,push_count,
	sys_user_id,address,
	min_zmscore,max_zmscore,
	app,login_time,push_time,
	remark,operator,
	loan_demand,create_time )
	VALUES
	(?,?,?,?,?,?,?,?,?,?,?,?,now()) `
	_, err := orm.NewOrm().Raw(sql, smsManagement.Content,
		smsManagement.PushCount, smsManagement.SysUserId, smsManagement.Address,
		smsManagement.MinZmscore, smsManagement.MaxZmscore,
		smsManagement.App, smsManagement.LoginTime, smsManagement.PushTime,
		smsManagement.Remark, smsManagement.Operator, smsManagement.LoanDemand).Exec()
	return err
}

//获取全部推送手机号码
func GetPushSMSUserAccount(condition string, params interface{}) (accounts []string, err error) {
	sql := `SELECT u.account FROM users u
	INNER JOIN users_metadata um
	ON um.uid = u.id
	INNER JOIN users_basedata ub
	ON u.id=ub.uid
	WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&accounts)
	return
}

//总发送短信条数
func GetPushSMSCount(uid int) (count int, err error) {
	sql := `SELECT SUM(push_count) FROM sms_management
	WHERE sys_user_id = ?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

//用户消息的总数
func GetUserMessageCount(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_message
	WHERE uid = ?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

//用户消息list
func GetUserMessageList(uid, begin, count int) (userMessage []*UserMessage, err error) {
	sql := `SELECT id,title,content,msg_type,is_read,create_time FROM users_message
	WHERE uid = ?
	ORDER BY create_time DESC
	LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, uid, begin, count).QueryRows(&userMessage)
	return
}
