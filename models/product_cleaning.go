package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CleaningData struct {
	Id             int
	Date           time.Time
	RegisterCount  int     //注册人数
	ApplyCount     int     //认证完成人数
	CreditCount    int     //授信人数
	ApplyLoanCount int     //申请借款人数
	MakeLoanCount  int     //放款人数
	MakeLoanAmount float64 //放款金额
	CpaMoney       float64 //收益
	State          int     //结算状态
	JointMode      int     //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CpaDefine      string  //CPA结算的有效事件
	CpaPrice       float64 //CPA的价格
	CpsPrice       float64 //CPS的价格
	AgentProductId int
}

type CleanNot struct {
	Id             int
	Date           time.Time
	RegisterCount  int     //注册人数
	ApplyCount     int     //认证完成人数
	CreditCount    int     //授信人数
	ApplyLoanCount int     //申请借款人数
	MakeLoanCount  int     //放款人数
	MakeLoanAmount float64 //放款金额
	CpaMoney       float64 //收益
	CleaningCycle  string  //结算周期
	QuantityCount  int     //有效数量
	Name           string  //产品名称
}

type CleanNotQuantity struct {
	Id             int
	Date           time.Time
	RegisterCount  int     //注册人数
	ApplyCount     int     //认证人数
	CreditCount    int     //授信人数
	ApplyLoanCount int     //申请借款人数
	MakeLoanCount  int     //放款人数
	MakeLoanAmount float64 //放款金额
	JointMode      int64   //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CpaDefine      string  //CPA结算的有效事件
	CleaningCycle  string
	QuantityCount  int     //有效数量
	Name           string  //产品名称
	CpaMoney       float64 //收益
}

//所有产品的结算状态和当日收益
type CleaningAll struct {
	State    int
	ProPrice float64
}

//开票信息
type ClHisInfo struct {
	Remark    string
	IsBilling int
}

//根据产品按时间查询数据明细
func GetProCleaningDay(id, begin, size int) (cleaningData []*CleaningData, err error) {
	sql := `
		SELECT dd.id,dd.agent_product_id,dd.register_count,dd.apply_count,dd.credit_count,dd.apply_loan_count,dd.make_loan_count,dd.make_loan_amount,dd.date,dd.state,
		dd.pro_price AS cpa_money,dd.cpa_price,dd.joint_mode,dd.cps_first_per AS cps_price
		FROM daily_data dd
		INNER JOIN agent_product ap
		ON dd.agent_product_id=ap.id
		INNER JOIN product_cleaning pc
	 	ON pc.product_id=ap.pro_id
	 	WHERE ap.pro_id = ? GROUP BY date ORDER BY date DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, id, begin, size).QueryRows(&cleaningData)
	if err != nil {
		beego.Debug(err.Error())
	}
	return
}

//获取所有H5产品数据
func GetAllProCleaning() (cleaningData []*CleaningAll, err error) {
	sql := `SELECT dd.state,dd.pro_price
		FROM daily_data dd
		INNER JOIN agent_product ap
		ON dd.agent_product_id=ap.id`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&cleaningData)
	return
}

//产品结算信息总数
func GetProCleaningCount(id int) (count int, err error) {
	sql := `SELECT count(1)
		FROM daily_data dd
		INNER JOIN agent_product ap
		ON dd.agent_product_id=ap.id
	 	INNER JOIN product_cleaning pc
	 	ON pc.product_id=ap.pro_id
	 	WHERE ap.pro_id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

//未结算产品总信息
func GetNotSettleProCleaning(order, condition string, params []string, begin, size int) (cleaningData []*CleanNotQuantity, err error) {
	/*sql := `SELECT p.id,SUM(dd.register_count) AS register_count, SUM(dd.apply_count) AS apply_count,
	SUM(dd.credit_count) AS credit_count, SUM(dd.apply_loan_count) AS apply_loan_count,
	SUM(dd.make_loan_count) AS make_loan_count, SUM(dd.make_loan_amount) AS make_loan_amount,
	SUM(dd.pro_price) AS cpa_money,MIN(dd.date) AS date,p.name,
	pc.cleaning_cycle
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	INNER JOIN product p ON ap.pro_id = p.id
	INNER JOIN product_cleaning pc ON ap.pro_id=pc.product_id
	WHERE dd.state = 0 `*/

	sql := `SELECT p.id,SUM(dd.register_count) AS register_count,SUM(dd.apply_count) AS apply_count,
	SUM(dd.credit_count)AS credit_count,SUM(dd.apply_loan_count) AS apply_loan_count,
	SUM(dd.make_loan_count) AS make_loan_count,SUM(dd.make_loan_amount) AS make_loan_amount,
	SUM(dd.pro_price) AS cpa_money,MIN(dd.date) AS date,p.name,
	pc.joint_mode,pc.cpa_define,pc.cleaning_cycle
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	INNER JOIN product p ON ap.pro_id = p.id
	INNER JOIN product_cleaning pc 	ON ap.pro_id= pc.product_id
	WHERE dd.state = 0 `
	if condition != "" {
		sql += condition
	}
	sql += " GROUP BY name "
	if order != "" {
		sql += order
	}
	sql += ` LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&cleaningData)
	return
}

//获取所以未结算产品的daily_data的id
func GetDailyDataId(pid int, condition string, params []string) (ids []int, err error) {
	sql := `SELECT dd.id
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	WHERE dd.state = 0
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY date DESC"
	_, err = orm.NewOrm().Raw(sql, pid, params).QueryRows(&ids)
	return
}

//单个产品未结算信息
func GetProNotSettle(pid int, condition string, params []string) (cleaningData []*CleanNot, err error) {
	sql := `SELECT dd.register_count,dd.apply_count,
	dd.credit_count,dd.apply_loan_count,
	dd.make_loan_count,dd.make_loan_amount,
	dd.pro_price AS cpa_money,dd.date
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	WHERE dd.state = 0
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY date DESC"
	_, err = orm.NewOrm().Raw(sql, pid, params).QueryRows(&cleaningData)
	return
}

//单个产品已结算信息
func GetProSettle(pid int, condition string, params []string) (cleaningData []*CleanNot, err error) {
	sql := `SELECT dd.register_count,dd.apply_count,
	dd.credit_count,dd.apply_loan_count,
	dd.make_loan_count,dd.make_loan_amount,
	dd.pro_price AS cpa_money,dd.date
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	WHERE dd.state = 1
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY date DESC"
	_, err = orm.NewOrm().Raw(sql, pid, params).QueryRows(&cleaningData)
	return
}

//获取指定时间内最小时间
func GetMinTime(pid int, condition string, params []string) (startDate string, err error) {
	sql := `SELECT MIN(date) FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	WHERE dd.state = 0
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, pid, params).QueryRow(&startDate)
	return
}

//单个产品未结算总数信息
func GetCleanCount(pid int, condition string, params []string) (cleaningData *CleanNot, err error) {
	sql := `SELECT SUM(dd.register_count) AS register_count,SUM(dd.apply_count) AS apply_count,
	SUM(dd.credit_count)AS credit_count,SUM(dd.apply_loan_count) AS apply_loan_count,
	SUM(dd.make_loan_count) AS make_loan_count,SUM(dd.make_loan_amount) AS make_loan_amount,
	SUM(dd.pro_price) AS cpa_money
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	WHERE dd.state = 0
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, pid, params).QueryRow(&cleaningData)
	return
}

//单个产品以结算总数信息
func GetCleanedCount(pid int, condition string, params []string) (cleaningData *CleanNot, err error) {
	sql := `SELECT SUM(dd.register_count) AS register_count,SUM(dd.apply_count) AS apply_count,
	SUM(dd.credit_count)AS credit_count,SUM(dd.apply_loan_count) AS apply_loan_count,
	SUM(dd.make_loan_count) AS make_loan_count,SUM(dd.make_loan_amount) AS make_loan_amount,
	SUM(dd.pro_price) AS cpa_money
	FROM daily_data dd
	INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
	WHERE dd.state = 1
	AND ap.pro_id = ?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, pid, params).QueryRow(&cleaningData)
	return
}

//未结算产品总数
func OneNotSettleCleaning(pid int, condition string, params []string) (count int, err error) {
	sql := `SELECT count(1)
		FROM daily_data dd
		INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
		WHERE dd.state =0
		 AND ap.pro_id=?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, pid, params).QueryRow(&count)
	return
}

//已结算产品总数
func OneSettleCleaning(pid int, condition string, params []string) (count int, err error) {
	sql := `SELECT count(1)
		FROM daily_data dd
		INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
		WHERE dd.state =1
		 AND ap.pro_id=?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, pid, params).QueryRow(&count)
	return
}

//未结算产品总数
func NotSettleCleaning(condition string, params []string) (count int, err error) {
	sql := `SELECT count( distinct(p.name))
		FROM daily_data dd
		INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
		INNER JOIN product p ON ap.pro_id = p.id
		INNER JOIN product_cleaning pc 	ON ap.pro_id= pc.product_id
		WHERE dd.state =0 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//撤销结算操作
func CancelSettlement(id int, ids []int) (err error) {
	o := orm.NewOrm()
	o.Begin()
	//更新daily_data数据
	sql := `UPDATE daily_data SET state =0 WHERE id=? `
	updateSql, err := orm.NewOrm().Raw(sql).Prepare()
	defer func() {
		updateSql.Close()
		if err != nil {
			o.Rollback()
			return
		} else {
			o.Commit()
		}
	}()
	if err != nil {
		return
	}
	for _, v := range ids {
		_, err = updateSql.Exec(v)
		if err != nil {
			return
		}
	}
	//删除cleaning_history数据
	sql = `DELETE FROM cleaning_history WHERE id =? `
	_, err = o.Raw(sql, id).Exec()
	if err != nil {
		return
	}
	return
}

//获取开票信息
func GetCleaningHistoryInfo(clId int) (clHisInfo ClHisInfo, err error) {
	sql := `SELECT remark,is_billing FROM cleaning_history WHERE id = ? `
	err = orm.NewOrm().Raw(sql, clId).QueryRow(&clHisInfo)
	return
}
