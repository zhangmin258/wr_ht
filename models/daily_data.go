package models

import (
	"github.com/astaxie/beego/orm"
)

/**
数据明细数据
*/
type DetailedData struct {
	Id                   int
	AgentProductId       int
	CreateDate           string
	RegisterCount        int     //注册数
	ApplyCount           int     //申请数
	CreditExtensionCount int     //授信数
	LoanApplyCount       int     //申请借款数
	CreditCount          int     //放款数
	CreditMoney          float32 //放款金额
	H5Earn               float32 //H5收益
}

func GetCountForH5(condition string, params []interface{}) (dailyData AgentDailyData, err error) {
	sql := `SELECT dd.id,dd.agent_product_id,SUM(register_count) register_count,SUM(apply_count) apply_count,SUM(credit_count) credit_count,
			SUM(apply_loan_count) apply_loan_count,SUM(make_loan_count) make_loan_count,SUM(make_loan_amount) make_loan_amount,
			DATE_FORMAT(dd.date,'%Y-%m-%d') date FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&dailyData)
	return
}
func GetCountForH5GroupByDate(condition string, params []int, pageStart, pageSize int) (dailyDatas []AgentDailyData, err error) {
	/*sql := `SELECT dd.id,dd.agent_product_id,SUM(register_count) register_count,SUM(apply_count) apply_count,
	SUM(credit_count) credit_count,SUM(apply_loan_count) apply_loan_count,SUM(make_loan_count) make_loan_count,SUM(make_loan_amount) make_loan_amount,
	DATE_FORMAT(dd.date, '%Y-%m-%d') date FROM daily_data dd
	INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
	INNER JOIN product p ON p.id = ap.pro_id
	INNER JOIN agent a ON ap.agent_id = a.id
	WHERE 1=1 `*/
	sql := `SELECT id,agent_product_id,register_count,apply_count ,
			credit_count,apply_loan_count,make_loan_count ,make_loan_amount,pro_price,
			DATE_FORMAT(date, '%Y-%m-%d') date FROM daily_data WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY date
			 ORDER BY date DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, pageStart, pageSize).QueryRows(&dailyDatas)
	return
}

func GetDaiDataCount(condition string, params []int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT DATE_FORMAT(date,"%Y-%m-%d") AS d FROM daily_data  WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY d ) a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

/**
根据条件获取注册用户
*/
func GetRegisterUsersByConditionForH5(condition string, params []string) (rus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(dd.date,'%Y-%m-%d') create_date,SUM(register_count) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1  `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(dd.date,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&rus)
	return
}

/**
根据条件获取申请用户
*/
func GetLoanUsersByConditionForH5(condition string, params []string) (lus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(dd.date,'%Y-%m-%d') create_date,SUM(apply_count) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(dd.date,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&lus)
	return
}

/**
根据条件获取授信用户
*/
func GetCreditExtensionUsersByConditionForH5(condition string, params []string) (ceus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(dd.date,'%Y-%m-%d') create_date,SUM(credit_count) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(dd.date,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&ceus)
	return
}

/**
根据条件获取申请借款用户
*/
func GetLoanTotalUsersByConditionForH5(condition string, params []string) (ltus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(dd.date,'%Y-%m-%d') create_date,SUM(apply_loan_count) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(dd.date,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&ltus)
	return
}

/**
根据条件获取放款用户
*/
func GetCreditUsersByConditionForH5(condition string, params []string) (cus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(dd.date,'%Y-%m-%d') create_date,SUM(make_loan_count) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(dd.date,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&cus)
	return
}

/**
根据条件获取放款金额
*/
func GetCreditMoneyUsersByConditionForH5(condition string, params []string) (cmu []CreditMoneyUser, err error) {
	sql := `SELECT DATE_FORMAT(dd.date,'%Y-%m-%d') create_date,SUM(make_loan_amount) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(dd.date,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&cmu)
	return
}

//根据条件获取h5产品收益
func GetProfitH5(condition string, params []string) (cmu []CreditMoneyUser, err error) {
	sql := `SELECT DATE(dd.date) create_date,SUM(pro_price) count FROM daily_data dd
			INNER JOIN agent_product ap ON ap.id = dd.agent_product_id
			INNER JOIN product p ON p.id = ap.pro_id
			INNER JOIN agent a ON ap.agent_id = a.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE(dd.date)
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&cmu)
	return
}
