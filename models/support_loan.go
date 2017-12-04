package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//付费用户列表
type PayUsersList struct {
	Uid        int       //id
	Name       string    //姓名
	Account    string    //手机号码
	IsValid    int       //是否有效
	IsFinished int       //是否完成
	CreateTime time.Time //下单时间
}

//退款
type Refund struct {
	Uid          int    //id
	RefundAmount string //退款金额
	CreateTime   string //用户订单时间
}

//付费用户导出Excel
type SupportLoanExcel struct {
	Account    string //手机号
	Name       string //姓名
	CreateTime string //下单时间
}

// 用户收支记录
type FinanceInfo struct {
	Uid               int       // 用户uid
	PayType           int       // 支付方式：1融豆 2现金余额 3会员免费
	DealType          int       // 交易类型：
	MoneyAmount       float64   // 人民币
	PayOrGet          int       // 1 付款 2 收款
	CreateTime        time.Time // 创建时间
	ServiceStates     int       // 状态 0 成功 1失败  2执行中
	BeforeScoreAmount int       `orm:"column(befor_score_amount)"` //交易前融豆余额
	AfterScoreAmount  int       //交易后融豆余额
	BeforeMoneyAmount float64   `orm:"column(befor_money_amount)"` //交易前的钱包余额
	AfterMoneyAmount  float64   //交易后的钱包余额
}

//获取贷款稳下价格
func GetSupportLoanInfo(id int) (moneyPrice float64, err error) {
	sql := `SELECT money_price FROM score_exchange_product WHERE id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&moneyPrice)
	return
}

//统计贷款稳下付费人数
func GetPayPersonCount() (count int, err error) {
	sql := `SELECT COUNT(1) AS count FROM support_loan`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//统计贷款稳下付费金额
func GetRefundAmount() (amount float64, err error) {
	sql := `SELECT SUM(refund_amount) FROM refund_record`
	err = orm.NewOrm().Raw(sql).QueryRow(&amount)
	return
}

//付费用户列表查询
func GetPayUsers(condition string, params []string, begin, size int) (payUsersList []PayUsersList, err error) {
	sql := `SELECT  
		s.uid,
		um.verify_real_name AS name,
		u.account,
		s.is_valid,
		s.is_finished,
		s.create_time 
		FROM support_loan AS s 
		LEFT JOIN users_metadata AS um 
		ON s.uid=um.uid 
		LEFT JOIN users AS u
		ON s.uid=u.id
		WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY s.create_time DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&payUsersList)
	return
}

//根据条件获取付费用户列表记录数量
func GetPayUsersCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) count
			FROM support_loan AS s 
			LEFT JOIN users_metadata AS um 
			ON s.uid=um.uid 
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//查询付款凭证
func GetPayToken(uid int, createTime string) (payOrder string, err error) {
	sql := `SELECT pay_order FROM support_loan WHERE uid=? AND create_time=?`
	err = orm.NewOrm().Raw(sql, uid, createTime).QueryRow(&payOrder)
	return
}

type RefundState struct {
	IsValid    int //是否失效
	IsFinished int //是否完成
}

//查询用户是否退满足退款要求
func GetUsersRefundState(uid int, createTime string) (refundState RefundState, err error) {
	sql := `SELECT is_valid,is_finished FROM support_loan WHERE uid=? AND create_time=?`
	err = orm.NewOrm().Raw(sql, uid, createTime).QueryRow(&refundState)
	return
}

//添加退款记录
func AddRefundRecord(uid, scoreProductId int, refundAmount float64, payOrder string) error {
	sql := `INSERT INTO refund_record (uid,score_product_id,refund_amount,submit_time,pay_order)VALUES(?,?,?,NOW(),?)`
	o := orm.NewOrm()
	_, err := o.Raw(sql, uid, scoreProductId, refundAmount, payOrder).Exec()
	return err
}

//退款后更改订单状态为失效、已完成
func ChangeSupportLoanState(uid int, createTime string) error {
	sql := `UPDATE support_loan SET is_valid=2,is_finished=2 WHERE uid=? AND create_time=?`
	o := orm.NewOrm()
	_, err := o.Raw(sql, uid, createTime).Exec()
	return err
}

//付费用户列表导出Excel
func ApplyExportExcel(condition string, params []string) (list []SupportLoanExcel, err error) {
	sql := `SELECT 
		s.create_time,
		um.account,
		um.verify_real_name AS name 
		FROM support_loan AS s LEFT JOIN users_metadata AS um 
		ON s.uid=um.uid 
		WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY s.create_time DESC `
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&list)
	return
}

// 查询用户融豆
func GetScore(uid int) (score int, err error) {
	sql := `SELECT score FROM score WHERE uid = ?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&score)
	return
}

// 添加收支记录
func AddRefundFinance(financeInfo FinanceInfo) (err error) {
	sql := `INSERT INTO users_finance_record (uid,pay_type,deal_type,money_amount,pay_or_get,create_time,service_states,befor_score_amount,after_score_amount,befor_money_amount,after_money_amount) VALUES(?,?,?,?,?,NOW(),?,?,?,?,?)`
	_, err = orm.NewOrm().Raw(sql, financeInfo.Uid, financeInfo.PayType, financeInfo.DealType, financeInfo.MoneyAmount, financeInfo.PayOrGet, financeInfo.ServiceStates, financeInfo.BeforeScoreAmount, financeInfo.AfterScoreAmount, financeInfo.BeforeMoneyAmount, financeInfo.AfterMoneyAmount).Exec()
	return
}
