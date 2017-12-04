package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type UsersFinanceRecord struct {
	Id               int
	Uid              int
	PayToken         string //付款凭证，用于调用业务时校验用户是否付款。
	PayType          int    //支付方式：1融豆 2现金余额 3会员免费
	DealType         int    //交易类型：1.话费1元券 2获取新口子 3.网贷征信查询 4.话费10元券 5.平台征信查询 6.1个月会员 7.2个月会员 8.购买抽奖 9.任务奖励 10.签到奖励 11 充值 12 提现 13 抽奖奖励 14 贷款稳下
	ScoreAmount      int
	MoneyAmount      float64
	PayOrGet         int //1 付款 2 收款
	CreateTime       time.Time
	Account          int       //手机号
	Content          string    //内容
	ServiceStates    int       //业务状态：0 成功 1 失败 2执行中
	TaskId           string    //话费充值时返回的taskId
	OrderCode        string    //充值或者提现时的连连订单号
	BeforScoreAmount int       //交易前融豆余额
	AfterScoreAmount int       //交易后融豆余额
	BeforMoneyAmount float64   //交易前的钱包余额
	AfterMoneyAmount float64   //交易后的钱包余额
	FinishTime       time.Time //完成时间（充值，提现等操作获得结果的时间）
}

//累计收入
func GetCounts(condition string, params []interface{}, deal_type int) (count float64, err error) {
	sql := `SELECT SUM(money_amount) FROM users_finance_record WHERE deal_type = ?`
	sql += condition
	err = orm.NewOrm().Raw(sql, params, deal_type).QueryRow(&count)
	return
}

//网贷征信查询累计收入
func GetCreditCounts(condition string, params []interface{}, deal_type int) (count float64, err error) {
	sql := `SELECT SUM(money_amount) FROM users_finance_record WHERE deal_type =? AND service_states=0`
	sql += condition
	err = orm.NewOrm().Raw(sql, params, deal_type).QueryRow(&count)
	return
}

//服务收益分页查询
func GetIncomeList(condition string, params []interface{}, begin, size int) (users []UsersFinanceRecord, err error) {
	sql := `SELECT
			u.create_time,
			us.account,
			s.content,
			u.money_amount,
	        u.service_states
		FROM
			users_finance_record u
		LEFT JOIN score_exchange_product s ON s.id = u.deal_type
		LEFT JOIN users us ON us.id = u.uid
		WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += ` AND u.pay_type=2 ORDER BY u.create_time DESC limit ?, ?`
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&users)
	if err != nil {
		return nil, err
	}
	return
}

//总计
func GetCount(condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT u.create_time FROM users_finance_record u WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += ` AND u.pay_type=2) AS a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//总计人数
func GetUsersCount(condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT us.account FROM users_finance_record u LEFT JOIN users us ON us.id = u.uid WHERE 1 = 1`
	if condition != "" {
		sql += condition
	}
	sql += ` AND u.pay_type=2 GROUP BY account) AS a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//总计收入金额
func GetTotalMoney(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT SUM(money_amount) FROM (SELECT u.money_amount FROM users_finance_record u WHERE 1 = 1`
	if condition != "" {
		sql += condition
	}
	sql += `) AS a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//失败或执行中金额
func GetMoney(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT SUM(money_amount) FROM (SELECT u.money_amount,u.service_states FROM users_finance_record u WHERE u.service_states IN (1,2) AND 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += `) AS a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//一元话费消耗
func GetOneBillCount(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT COUNT(1) FROM users_finance_record WHERE deal_type =1 AND service_states = 0`
	sql += condition
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//十元话费消耗
func GetTenBillCount(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT COUNT(1) FROM users_finance_record WHERE deal_type =4 AND service_states = 0`
	sql += condition
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//现金消耗
func GetMoneyCount(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT SUM(money_amount) FROM users_finance_record WHERE deal_type=13`
	sql += condition
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//活动消耗分页查询
func GetActivityList(source string, condition string, params []interface{}, begin, size int) (users []UsersFinanceRecord, err error) {
	sql := ` `
	if source == "1,4" || source == "13"{
		sql += `SELECT u.create_time,us.account,
	         CASE deal_type
				WHEN 1 THEN '一元话费券'
				WHEN 4 THEN '十元话费券'
				WHEN 13 THEN '抽奖奖励'
				END AS content,
				u.money_amount
				FROM users_finance_record u
				LEFT JOIN score_exchange_product s ON s.id = u.deal_type
				LEFT JOIN users us ON us.id = u.uid
				WHERE 1 = 1`
	}else{
		sql += ` SELECT * FROM (
					SELECT u.create_time,us.account,
					CASE deal_type
					WHEN 1 THEN '一元话费券'
					WHEN 4 THEN '十元话费券'
					END AS content, u.money_amount
					FROM users_finance_record u
					LEFT JOIN score_exchange_product s ON s.id = u.deal_type
					LEFT JOIN users us ON us.id = u.uid
					WHERE 1 = 1
					AND u.service_states = 0
					AND u.deal_type IN (1, 4)
					UNION ALL
					SELECT u.create_time,us.account,
					CASE deal_type
					WHEN 13 THEN '抽奖奖励'
					END AS content,u.money_amount
					FROM users_finance_record u
					LEFT JOIN score_exchange_product s ON s.id = u.deal_type
					LEFT JOIN users us ON us.id = u.uid
					WHERE 1 = 1
					AND u.deal_type IN (13)
					AND u.money_amount > 0) AS u WHERE 1=1`
	}

	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY create_time DESC limit ?, ?`
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&users)
	if err != nil {
		return nil, err
	}
	return
}

//总计
func GetActivityCount(source string,condition string, params []interface{}) (count int, err error) {
	sql := ` `
	if source == "1,4" || source == "13" {
		sql += ` SELECT COUNT(1) FROM (SELECT u.create_time,us.account FROM users_finance_record u LEFT JOIN users us ON us.id = u.uid WHERE 1=1 `
		if condition != "" {
			sql += condition
		}
		sql += ` ) AS a`
	}else{
		sql += ` SELECT COUNT(1) FROM(
		SELECT u.create_time,us.account
		FROM users_finance_record u
		LEFT JOIN users us ON us.id = u.uid
		WHERE u.deal_type IN (13)
        AND u.money_amount > 0
		UNION ALL
		SELECT u.create_time,us.account
		FROM users_finance_record u
		LEFT JOIN users us ON us.id = u.uid
		WHERE u.service_states = 0
		AND u.deal_type IN (1, 4)) AS u WHERE 1=1`
		if condition != "" {
			sql += condition
		}
	}

	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//总计人数
func GetActivityUsersCount(source,condition string, params []interface{}) (count int, err error) {
	sql := ` `
	if source == "1,4" || source == "13" {
		sql += `SELECT COUNT(1) FROM (SELECT us.account FROM users_finance_record u LEFT JOIN users us ON us.id = u.uid WHERE 1 = 1`
		if condition != "" {
			sql += condition
		}
		sql += ` GROUP BY account) AS a`
		err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	}else{
		sql += ` SELECT COUNT(1) FROM (
		    SELECT account
		    FROM(SELECT us.account
				FROM users_finance_record u
				LEFT JOIN users us ON us.id = u.uid
				WHERE 1=1 `
		if condition != "" {
			sql += condition
		}
		sql += `
		        AND u.deal_type IN (13)
				AND u.money_amount > 0
				UNION ALL
				SELECT us.account
				FROM users_finance_record u
				LEFT JOIN users us ON us.id = u.uid
				WHERE 1=1 `
		if condition != "" {
			sql += condition
		}
		sql += `  AND u.service_states = 0
				AND u.deal_type IN (1, 4)) AS a
                GROUP BY account) AS b`
		err = orm.NewOrm().Raw(sql, params,params).QueryRow(&count)
	}
	return
}

//获取话费支出金额
func GetBillCount(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT SUM(money_amount)
             FROM(SELECT
			 CASE deal_type
		      WHEN 1 THEN 1
		      WHEN 4 THEN 10
		      END AS money_amount
	      	FROM users_finance_record u
		    WHERE 1 = 1`
	if condition != "" {
		sql += condition
	}
	sql += `) AS a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//获取所有项目注册支出金额
func GetAllActivitymoney(condition string, params []interface{}) (count float64, err error) {
	sql := `SELECT SUM(money_amount) FROM(
		     SELECT
		     create_time,
			 CASE deal_type
		     WHEN 1 THEN 1
		     WHEN 4 THEN 10
		     END AS money_amount
		     FROM users_finance_record
		     WHERE deal_type IN (1,4)
	         AND service_states = 0
		     UNION ALL
			 SELECT create_time,money_amount
			 FROM users_finance_record
			 WHERE deal_type = 13) AS u WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}
