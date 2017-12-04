package models

import (
	"github.com/astaxie/beego/orm"
)

type ProductInfo struct {
	Id        int
	Name      string  //平台名称
	JointMode int64   //合作方式
	CpaDefine string  //有效性定义
	Price     float64 //价格
	CpaPrice  float64 //cpa价格
	CpsPrice  float64 //cps价格
}

type QuantityCount struct {
	RegisterCount  int     //注册人数
	ApplyCount     int     //认证人数
	CreditCount    int     //授信人数
	ApplyLoanCount int     //申请借款人数
	MakeLoanCount  int     //放款人数
	MakeLoanAmount float64 //放款金额
	JointMode      int64   //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CpaDefine      string  //CPA结算的有效事件
	CleaningCycle  string
}

//根据pid获取产品信息
func GetProInfoByPid(pid int) (productInfo *ProductInfo, err error) {
	sql := `SELECT p.id,p.name,pc.joint_mode,pc.cpa_define,pc.cpa_price,pc.cps_price FROM product p INNER JOIN agent_product ap
		ON p.id=ap.pro_id
		INNER JOIN product_cleaning pc
		ON ap.pro_id= pc.product_id
		WHERE p.id =? limit 1`
	err = orm.NewOrm().Raw(sql, pid).QueryRow(&productInfo)
	return
}

//保存清算历史记录
func SaveCleaningHistory(pid, isbilling int, DailyDataIds, remark, condition string, cleanCount *CleanNot, params []string) (err error) {
	o := orm.NewOrm()
	o.Begin()
	var ids []int
	//查出所有agent_product的id
	sql := `SELECT id	FROM agent_product WHERE pro_id = ?`
	_, err = o.Raw(sql, pid).QueryRows(&ids)
	//更新指定时间段daily_data的状态
	sql = `UPDATE daily_data dd SET dd.state=1
	        WHERE dd.agent_product_id=? `
	if condition != "" {
		sql += condition
	}
	usql, err := o.Raw(sql).Prepare()
	for _, v := range ids {
		_, err = usql.Exec(v, params[0], params[1])
	}
	//结算信息存入清算历史
	sql = `INSERT INTO cleaning_history (daily_data_id,pro_id,settle_time,begin_time,end_time,
	quantity_count,make_loan_amount,settle_money,is_billing,remark) VALUES(?,?,NOW(),?,?,?,?,?,?,?)`
	_, err = o.Raw(sql, DailyDataIds, pid, params, cleanCount.QuantityCount, cleanCount.MakeLoanAmount, cleanCount.CpaMoney, isbilling, remark).Exec()
	defer func() {
		usql.Close()
		if err != nil {
			o.Rollback()
			return
		}
		o.Commit()
	}()
	return
}

//获取有效数据和放款金额总数
func GetQuantityCount(pid int, condition string, params []string) (cleaningData *QuantityCount, err error) {
	sql := `SELECT SUM(dd.register_count) AS register_count,SUM(dd.apply_count) AS apply_count,
	SUM(dd.credit_count)AS credit_count,SUM(dd.apply_loan_count) AS apply_loan_count,
	SUM(dd.make_loan_count) AS make_loan_count,SUM(dd.make_loan_amount) AS make_loan_amount,
	pc.joint_mode,pc.cpa_define
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	INNER JOIN product_cleaning pc 	ON ap.pro_id= pc.product_id
	WHERE dd.state = 0
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, pid, params).QueryRow(&cleaningData)
	return
}
