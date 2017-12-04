package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type UserLoan struct {
	Id              int       //buisness_loan id
	Business_id     int       //
	Uid             string    //用户id
	Name            string    //产品名称
	Money           float64   //借款金额
	Loan_term_count int       //借款期限
	Create_time     time.Time //批复时间
	Finish_time     time.Time //还款时间
	State           string    // 状态
}

type UserOrder struct {
	Id               int       //产品id
	Name             string    //产品名称
	Account          string    //用户手机号
	Create_time      time.Time //借款时间
	Verify_real_name string    //用户姓名
	Money            float64   //借款金额
	Loan_term_count  int       //借款期限
}

func GetLoanList(condition string, begin, count int, pars ...string) (list []UserLoan, err error) {
	sql := `SELECT l.id,l.business_id,l.uid,p.name,l.money,l.loan_term_count,l.create_time,l.finish_time,l.state FROM business_loan l LEFT JOIN product p on p.id=l.product_id WHERE 1=1   `
	if condition != "" {
		sql += condition
	}
	sql += " order by l.create_time desc limit ?, ?"
	_, err = orm.NewOrm().Raw(sql, pars, begin, count).QueryRows(&list)
	return
}

func GetLoanCount(condition string, pars ...string) int {
	sql := `SELECT count(l.id) FROM business_loan l LEFT JOIN product p on p.id=l.product_code WHERE 1=1  `
	if condition != "" {
		sql += condition
	}
	var count int
	o := orm.NewOrm()
	o.Raw(sql, pars).QueryRow(&count)
	return count
}

func GetLoanMoneyMax(condition string, pars ...string) float64 {
	sql := `SELECT SUM(l.money) AS count FROM business_loan l LEFT JOIN users_metadata um ON l.uid=um.uid left join product p on p.id=l.product_code
			WHERE 1=1  `
	if condition != "" {
		sql += condition
	}
	o := orm.NewOrm()
	var count float64
	o.Raw(sql, pars).QueryRow(&count)
	return count
}

func GetProOrderlist(id string) (list *UserOrder) {
	sql := `SELECT l.id,p.name,um.account,l.create_time,um.verify_real_name,l.money,l.loan_term_count FROM business_loan l LEFT JOIN users_metadata um ON l.uid=um.uid left join product p on p.id=l.product_code WHERE l.id = ?`
	orm.NewOrm().Raw(sql, id).QueryRow(&list)
	return
}
