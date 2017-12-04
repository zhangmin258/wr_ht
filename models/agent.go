package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//代理商信息
type Agent struct {
	Id          int       `orm:"pk"` //(11)主键
	OrgName     string    //机构信息：名称
	OrgAddress  string    //机构信息：通讯录地址
	OrgEmail    string    //机构信息：邮箱
	BusName     string    //商务联系人：联系人姓名
	BusPhone    string    //商务联系人：电话
	BusWechat   string    //商务联系人：微信号
	BusQQ       string    `orm:"column(bus_qq)"` //商务联系人：qq
	BusEmail    string    //商务联系人：邮箱
	AccountName string    //开票：账户名称
	AccountBank string    //开票：开户银行
	BankAccount string    //开票：银行账号
	BackUrl     string    //后台：后台网址
	BackAccount string    //后台：后台账号
	BackPwd     string    //后台:密码
	FromNum     string    //后台:渠道号
	CreateTime  time.Time //创建时间
}

//条件分页查询代理商信息
func GetAgentList(condition string, params []string, begin, size int) (agents []Agent, err error) {
	sql := ` SELECT id, org_name,bus_name,bus_phone FROM agent WHERE id >0`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY create_time DESC LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&agents)
	return
}

// 条件查询所有代理商数量
func GetAgentCount(condition string, params []string) (count int, err error) {
	sql := `SELECT count(1) FROM agent p WHERE id >0`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}
func GetAgentIdByName(agentName string) (id int, err error) {
	sql := `SELECT id FROM agent WHERE org_name = ?`
	err = orm.NewOrm().Raw(sql, agentName).QueryRow(&id)
	return
}
func InsertAgentProduct(productId, agentId int) (id int64, err error) {
	sql := `INSERT INTO agent_product(pro_id,agent_id) VALUES(?,?)`
	rs := orm.NewOrm().Raw(sql, productId, agentId)
	result, err := rs.Exec()
	if err == nil {
		id, err = result.LastInsertId()
	}
	return
}

// 保存下级代理数据
func AddLowerLevel(agent *Agent) (pid int, err error) {
	agent.CreateTime = time.Now()
	id, err := orm.NewOrm().Insert(agent)
	return int(id), err
}

// update下级代理数据
func EditLowerLevel(agent *Agent) (err error) {
	sql := `SELECT create_time FROM agent WHERE id=?`
	var time time.Time
	orm.NewOrm().Raw(sql, agent.Id).QueryRow(&time)
	agent.CreateTime = time
	_, err = orm.NewOrm().Update(agent)
	return
}

// 根据id查询代理商
func GetAgentById(id int) (agent Agent, err error) {
	sql := `SELECT id,org_name,bus_name,
	org_address,org_email,
	bus_phone,bus_wechat,
	bus_qq,bus_email,
	account_name,account_bank,back_account,
	bank_account,back_url,back_pwd,
	from_num,create_time
	FROM agent WHERE id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&agent)
	return
}
