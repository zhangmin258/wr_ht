package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Wallet struct {
	Uid            int
	AccountBalance int
	IsFreeze       int
}

type Score struct {
	Uid      int
	Score    int
	IsFreeze int
}

type PayPwd struct {
	PayPwd string
}

type UpdatePayPwd struct {
	PayPwd  string
	Vcode   string
	Account string
	Uid     int
}

//数据库取出钱包收支记录
type WalletRecord struct {
	CreateTime       time.Time //创建时间
	PayOrGet         int       //收入或支出 0:支出 1:收入
	DealType         int       //交易类型(项目)1.话费1元券 2获取新口子 3.网贷征信查询 4.话费10元券 5.平台征信查询 6.1个月会员 7.2个月会员 8.购买抽奖 9.任务奖励 10.签到奖励 11 充值 12 提现 13 抽奖奖励 14 贷款稳下
	MoneyAmount      float64   //交易金额
	AfterMoneyAmount float64   //钱包余额
}

//返回钱包收支记录
type UserWalletRecord struct {
	CreateTime       string //创建时间
	PayOrGet         string //收入或支出
	DealType         string //交易类型(项目)
	MoneyAmount      string //交易金额
	AfterMoneyAmount string //钱包余额
}

//积分消费记录
type ScorePayRecord struct {
	CreateTime       time.Time //创建时间
	DealType         int       //交易类型
	ScoreAmount      int       //融豆价格
	AfterScoreAmount int       //积分余额
}

//积分兑换产品
type ScoreExchangeProduct struct {
	Content    string //产品内容
	ScorePrice int    //兑换所需积分
}

// //返回积分兑换抽奖记录
// type ScoreExchangeLotteryRecord struct {
// 	CreateTime       string //时间
// 	ProductType      string //消耗方式
// 	Prize            string //奖品名称
// 	ScoreAmount      string //积分消耗
// 	AfterScoreAmount string //积分余额
// }

//返回积分兑换记录
type ScoreExchangeRecord struct {
	CreateTime       string //时间
	Prize            string //奖品名称
	ScoreAmount      string //积分消耗
	AfterScoreAmount string //积分余额
}

//抽奖记录
type LotteryRecord struct {
	CreateTime time.Time
	Content    string
}

//返回积分抽奖记录
type ScoreLotteryRecord struct {
	CreateTime string
	Content    string
}

// 获取钱包账户余额
func GetWalletBalance(uid int) (balance float64, err error) {
	sql := `SELECT account_balance FROM wallet WHERE uid=?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&balance)
	return
}

// 初始化用户钱包数据
func InitWallet(uid int) (err error) {
	sql := `INSERT INTO wallet (uid,account_balance) VALUES (?,0.00)`
	_, err = orm.NewOrm().Raw(sql, uid).Exec()
	return
}

// 用户支付，扣除余额
func WalletPay(uid int, money float64) (err error) {
	sql := `UPDATE wallet SET account_balance=account_balance-? WHERE uid=?`
	_, err = orm.NewOrm().Raw(sql, money, uid).Exec()
	return
}

// 充值
func WalletRecharge(uid int, money float64) (err error) {
	sql := `UPDATE wallet SET account_balance=account_balance+? WHERE uid=?`
	_, err = orm.NewOrm().Raw(sql, money, uid).Exec()
	return
}

// 充值-事务
func WalletRechargeTransaction(uid int, money float64, o orm.Ormer) (err error) {
	sql := `UPDATE wallet SET account_balance=account_balance+? WHERE uid=?`
	_, err = o.Raw(sql, money, uid).Exec()
	return
}

// 查询订单状态
func GetUserRechargeStatus(order_code string) (status string, err error) {
	sql := `SELECT status FROM users_recharge_record WHERE order_code=?`
	err = orm.NewOrm().Raw(sql, order_code).QueryRow(&status)
	return
}

// 通过订单号查询用户id和金额
func GetUidByOrderNumber(orderNo string) (uid int, amount float64, err error) {
	sql := `SELECT uid,amount FROM users_withdraw_deposit_records WHERE order_code=?`
	err = orm.NewOrm().Raw(sql, orderNo).QueryRow(&uid, &amount)
	return
}

// 修改提现订单状态
func ModifyWithdrawDepositType(result_pay, ret_code, order_code string, order_type int, o orm.Ormer) (err error) {
	sql := `UPDATE users_withdraw_deposit_records SET type=?,result_pay=?,ret_code=? WHERE order_code=? `
	_, err = o.Raw(sql, order_type, result_pay, ret_code, order_code).Exec()
	return
}

// 得到结果时，修改收支记录
func UpdateRechargeFinance(money float64, status int, orderNo string, o orm.Ormer) (err error) {
	if status == 1 {
		sql := `UPDATE users_finance_record SET finish_time=NOW(), after_money_amount=?,service_states=? WHERE order_code=?`
		_, err = o.Raw(sql, money, status, orderNo).Exec()
	} else {
		sql := `UPDATE users_finance_record SET finish_time=NOW(),service_states=? WHERE order_code=?`
		_, err = o.Raw(sql, status, orderNo).Exec()
	}
	return

}

//获取用户钱包记录
func GetUserWalletRecords(uid, begin, count int) (records []WalletRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT create_time,pay_or_get, deal_type,money_amount,after_money_amount 
			FROM users_finance_record 
			WHERE pay_type = 2 AND uid = ? 
			ORDER BY create_time DESC LIMIT ?, ?`
	_, err = o.Raw(sql, uid, begin, count).QueryRows(&records)
	return
}

//获取用户钱包支出记录数量
func GetUserWalletRecordsCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) 
			FROM users_finance_record 
			WHERE pay_type = 2 AND uid = ?`
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//获取用户积分兑换记录
func GetUserScoreExchangeRecords(uid int) (records []ScorePayRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT create_time,deal_type,score_amount,after_score_amount FROM users_finance_record
			WHERE deal_type < 9 AND uid = ? ORDER BY create_time DESC`
	_, err = o.Raw(sql, uid).QueryRows(&records)
	return
}

//获取积分兑换产品
func GetScoreExchangeProducts() (products []ScoreExchangeProduct, err error) {
	o := orm.NewOrm()
	sql := `SELECT content,score_price FROM score_exchange_product`
	_, err = o.Raw(sql).QueryRows(&products)
	return
}

//获取用户积分抽奖记录
func GetUserLotteryRecords(uid int) (records []LotteryRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT slr.create_time,slp.content FROM score_lottery_record slr 
			LEFT JOIN score_lottery_product slp ON slr.sl_product_id = slp.id 
			WHERE slr.uid = ? ORDER BY slr.create_time DESC`
	_, err = o.Raw(sql, uid).QueryRows(&records)
	return
}
