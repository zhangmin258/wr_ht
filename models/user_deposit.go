package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//用户佣金
type UserDeposit struct {
	Id          int
	Uid         int       //用户id
	Account     string    //用户手机号
	Amount      float64   //用户提现金额
	Real_amount float64   //实际提现金额
	CheckResult int64     //1:等待连连回调2：需要人工审核3：正常放款4：放款失败5:拒绝
	AmountTime  time.Time //提现时间
	CheckTime   time.Time //审批时间
	CheckPeople int       //审批人
	OrderCode   string    //商户编号
	State       int64     //用户被拒绝状态：大于0为拒绝过
}

//用户邀请的好友信息
type UserInvitation struct {
	Id             int
	Account        string    //好友手机号
	MobileType     string    //好友设备型号
	State          int       `orm:"column(type)"` //1:注册2:激活3:完成
	Money          float64   //获得金额
	InvitationTime time.Time //邀请时间
}

//用户银行卡
type UsersBankcard struct {
	Id         int       `orm:"column(id);pk"`
	Uid        int       `orm:"column(uid)"`
	BankName   string    `orm:"column(bank_name)"`
	CardNumber string    `orm:"column(card_number);size(30)"`
	State      string    `orm:"column(state)"`
	CreateTime time.Time `orm:"column(create_time)"`
	Account    string    `orm:"column(account)"`
	BankMobile string    `orm:"column(bank_mobile)"` // 银行的预留手机号
}

type UserDepositInfo struct {
	Uid    int
	Amount float64
}

//获取所有佣金列表
func GetUserDeposit(condition string, params []string, begin, size int) (userDeposit []*UserDeposit, err error) {
	sql := `SELECT ur.id,ur.uid,ur.amount_time,ur.type AS check_result,u.account,ur.amount,a.state FROM users_withdraw_deposit_records ur
	LEFT JOIN (SELECT uid,count(uid) as state from users_withdraw_deposit_records where type=5 group by uid ) a
	ON a.uid= ur.uid
	INNER JOIN users u ON ur.uid = u.id
	 LEFT JOIN wallet w
	 ON ur.uid=w.uid
	 WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY ur.amount_time DESC LIMIT ?,? "
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&userDeposit)
	return
}

//获取所有佣金列表
func GetUserDepositCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_withdraw_deposit_records ur
	INNER JOIN users u ON ur.uid = u.id
	 LEFT JOIN wallet w
	 ON ur.uid=w.uid
		WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//获取用户邀请好友列表
func UserInvitationList(uid int, condition string, params []string, begin, size int) (userInvitation []*UserInvitation, err error) {
	sql := `SELECT ir.id,u.account,u.mobile_type,ir.type,ir.money,ir.invitation_time FROM invitation_record ir
	LEFT JOIN users u ON ir.new_uid = u.id WHERE ir.old_uid=? `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY ir.invitation_time LIMIT ?,? "
	_, err = orm.NewOrm().Raw(sql, uid, params, begin, size).QueryRows(&userInvitation)
	return
}

//获取用户邀请好友总数
func UserInvitationCount(uid int, condition string, params []string) (count int, err error) {
	sql := `SELECT count(1) FROM invitation_record ir
	LEFT JOIN users u ON ir.new_uid = u.id WHERE ir.old_uid=? `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, uid, params).QueryRow(&count)
	return
}

//获取用户详细佣金信息
func GetUserDepositDetail(uid int, begin, size int) (userDeposit []*UserDeposit, err error) {
	sql := `SELECT ur.amount_time,ur.check_time,ur.amount,ur.type AS check_result FROM users_withdraw_deposit_records ur
			INNER JOIN users u ON ur.uid = u.id
			WHERE ur.uid = ?
			ORDER BY ur.amount_time DESC LIMIT ?,? `
	_, err = orm.NewOrm().Raw(sql, uid, begin, size).QueryRows(&userDeposit)
	return
}

//获取用户详细佣金信息总数
func GetUserDepositDetailCount(uid int) (count int, err error) {
	sql := `SELECT count(1) FROM users_withdraw_deposit_records ur
			INNER JOIN users u ON ur.uid = u.id
			WHERE ur.uid = ? `
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

// 根据id获取提现信息
func GetDepositById(id int) (userDeposit *UserDeposit, err error) {
	sql := `SELECT id,uid,order_code,amount,real_amount FROM users_withdraw_deposit_records WHERE id = ? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&userDeposit)
	return
}

// 获取单张银行卡信息
func GerUsersBankcardById(Uid int) (bkcard *UsersBankcard, err error) {
	o := orm.NewOrm()
	err = o.Raw(`SELECT id,uid,bank_name,card_number,state,create_time,account,bank_mobile FROM users_bankcards WHERE uid = ?`, Uid).QueryRow(&bkcard)
	return
}

// 查询订单状态
func GetUserWithdrawDepositStatus(order_code string) (resultPay string, err error) {
	sql := `SELECT result_pay FROM users_withdraw_deposit_records WHERE order_code=?`
	err = orm.NewOrm().Raw(sql, order_code).QueryRow(&resultPay)
	return
}

// 将订单设置为异常订单
func SetIllegalOrder(order_code string) (err error) {
	sql := `UPDATE users_withdraw_deposit_records SET is_illegal_order=1 WHERE order_code=?`
	_, err = orm.NewOrm().Raw(sql, order_code).Exec()
	return
}

func RefuseWithdrawDeposit(userDeposit *UserDeposit, sysId int) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err != nil {
			o.Rollback()
			return
		}
		o.Commit()
	}()
	//修改提现记录表的状态
	sql := `UPDATE users_withdraw_deposit_records SET type=5,check_people=?,check_time=NOW() WHERE uid= ? AND order_code=? `
	_, err = o.Raw(sql, sysId, userDeposit.Uid, userDeposit.OrderCode).Exec()
	if err != nil {
		return
	}
	/*	//修改提现表的用户记录状态
		sql = `UPDATE users_withdraw_deposit SET state=1 WHERE uid = ?`
		_, err = o.Raw(sql, userDeposit.Uid).Exec()
		if err != nil {
			return
		}*/
	//修改钱包余额
	sql = `UPDATE wallet SET account_balance=account_balance+? WHERE uid=?`
	_, err = o.Raw(sql, userDeposit.Amount, userDeposit.Uid).Exec()
	if err != nil {
		return
	}
	//查出用户的余额
	money := 0.00
	sql = `SELECT account_balance FROM wallet WHERE uid = ? `
	err = o.Raw(sql, userDeposit.Uid).QueryRow(&money)
	if err != nil {
		return
	}
	//修改用户收支记录表
	sql = `UPDATE users_finance_record SET service_states=1,after_money_amount=? WHERE uid=? AND order_code=?`
	_, err = o.Raw(sql, money, userDeposit.Uid, userDeposit.OrderCode).Exec()
	return
}
