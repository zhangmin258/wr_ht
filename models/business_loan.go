package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//订单列表展示
type BusinessLoan struct {
	Id            int       `orm:"pk"`                       //商业贷款id
	Uid           int       `orm:"column(uid)"`              //用户ID
	UserName      string    `orm:"column(verify_real_name)"` //用户名字verify_real_name
	Account       string    `orm:"column(account)"`          //用户手机号
	LoanTime      time.Time `orm:"column(loan_time)"`        //申请时间
	Money         float64   `orm:"column(money)"`            //借款金额
	LoanTermCount int       `orm:"column(loan_term_count)"`  //借款期限
	State         string    //订单状态
	FinishTime    time.Time `orm:"column(finish_time)"`  //还款时间
	ReturnMoney   float64   `orm:"column(return_money)"` //应还金额
	OrgId         int       `orm:"column(business_id)"`  //机构Id
	OrgName       string    `orm:"column(orgname)"`      //机构名称
}

//条件分页查询贷款人api
func GetBusinessLoanPageApi(productId int, date string, condition string, params []string, begin, size int) (businessLoan []BusinessLoan, err error) {
	o := orm.NewOrm()
	sql := `select 
 			a.id, a.uid, b.verify_real_name, b.account,a.loan_time,a.money,a.return_money,a.loan_term_count,a.state,a.business_id,d.name as orgname
            from
            business_loan as a
            inner join users_metadata as b
            on a.uid=b.uid
            left join business as d
            on a.business_id=d.id
            where  a.product_id = ? AND
            date_format(a.loan_time , "%Y-%m-%d") = ?  `
	if condition != "" {
		sql += condition
	}
	sql += " order by  id limit ?,? "
	_, err = o.Raw(sql, productId, date, params, begin, size).QueryRows(&businessLoan)
	return
}

//条件分页查询贷款人h5
func GetBusinessLoanPageH5(productId int, date string, condition string, params []string, begin, size int) (businessLoan []BusinessLoan, err error) {
	o := orm.NewOrm()
	sql := `select
 			a.id, a.uid, b.verify_real_name, b.account,a.create_time AS loan_time,a.loan_money AS money,a.return_money,a.loan_term_count
            from
            business_loan_h5 as a
            inner join users_metadata as b
            on a.uid=b.uid
            where  a.product_id = ? AND
            date_format(a.create_time , "%Y-%m-%d") = ?  `
	if condition != "" {
		sql += condition
	}
	sql += " order by  id limit ?,? "
	_, err = o.Raw(sql, productId, date, params, begin, size).QueryRows(&businessLoan)
	return
}

// 条件查询所有贷款人的数量api
func GetBusinessLoanListApi(productId int, date string, condition string, params []string) (count int, err error) {
	o := orm.NewOrm()
	sql := `select count(1)
			from 
			business_loan as a
            inner join users_metadata as b
            on a.uid=b.uid
            left join business as d
            on a.business_id=d.id
            where  a.product_id = ? AND
            date_format(a.loan_time , "%Y-%m-%d") = ? `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, productId, date, params).QueryRow(&count)
	return
}

// 条件查询所有贷款人的数量H5
func GetBusinessLoanListH5(productId int, date string, condition string, params []string) (count int, err error) {
	o := orm.NewOrm()
	sql := `select count(1)
			from
			business_loan_h5 as a
            inner join users_metadata as b
            on a.uid=b.uid
            where  a.product_id = ? AND
            date_format(a.create_time , "%Y-%m-%d") = ? `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, productId, date, params).QueryRow(&count)
	return
}

//根据id查询
func GetBusinessLoanById(id int) (businessloan BusinessLoan, err error) {
	o := orm.NewOrm()
	sql := `select 
 			a.id, a.uid, b.verify_real_name, b.account,a.loan_time,a.money,a.return_money,a.loan_term_count,a.state,a.business_id,d.name as orgname
            from
            business_loan as a
            inner join users_metadata as b
            on a.uid=b.uid
            left join business as d
            on a.business_id=d.id
			where a.id = ? `
	err = o.Raw(sql, id).QueryRow(&businessloan)
	return
}

//根据id查询
func GetBusinessLoanH5ById(id int) (businessloan BusinessLoan, err error) {
	o := orm.NewOrm()
	sql := `SELECT h.id,um.verify_real_name,um.account,h.create_time AS loan_time,p.name AS orgname,h.loan_money AS money,h.loan_term_count
			FROM business_loan_h5 h 
			LEFT JOIN users_metadata um 
			ON h.uid = um.uid 
			LEFT JOIN product p 
			ON h.product_id = p.id
			WHERE h.id = ? `
	err = o.Raw(sql, id).QueryRow(&businessloan)
	return
}

// 条件查询某一个贷款人所有订单
func GetBusinessLoanListById(condition string, params []string, id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `select count(1) from business_loan where id =?`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, params, id).QueryRow(&count)
	return
}
