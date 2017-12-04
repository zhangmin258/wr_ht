package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type AgentDailyData struct {
	Id             int
	AgentProductId int       //代理产品Id
	RegisterCount  int       //注册人数
	ApplyCount     int       //申请人数
	CreditCount    int       //授信人数
	ApplyLoanCount int       //申请借款人数
	MakeLoanCount  int       //放款人数
	MakeLoanAmount float32   //放款金额
	Date           time.Time //日期
	ProPrice       float64   //收益
	State          int       //结算状态
	JointMode      int       //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CpaDefine      string    //CPA结算的有效事件
	CpaPrice       float64   //CPA的价格
	CpsFirstPer    float64   //cps的首借百分比
	CpsAgainPer    float64   //cps的复借百分比
}

type AnalysisProduct struct {
	Name                 string    //产品名称
	Loan_money           float64   //放款金额
	HitsCount            int       //点击次数
	RegisterCount        int       //注册人数
	PlatformRegister     int       //平台注册人数
	Credit               float64   //信用
	ApplyRate            float64   //申请通过率
	Cac                  float64   //获客成本
	Amount               float64   //件均金额
	HitsRegister         float64   //我司统计点击-注册转化率
	PlatformHitsRegister float64   //平台点击-注册转化率
	Data                 time.Time //日期
}

//产品列表
type ProductIdList struct {
	Id   int    //产品id
	Name string //产品名字
	Type int    //产品类型
}

//产品人数
type ProductCount struct {
	Id    int
	Count int
}

//产品金额
type ProductMoney struct {
	Id    int
	Money float64
}

//产品点击次数和注册人事
type ProductHitsAndRegister struct {
	Id            int
	HitsCount     int
	RegisterCount int
}

//指定日期产品人数
type ProductCountDate struct {
	Data  time.Time
	Count int
}

//指定日期产品金额
type ProductMoneyData struct {
	Data  time.Time
	Money float64
}

//指定日期产品点击次数和注册人事
type ProductHitsAndRegisterData struct {
	Data          time.Time
	HitsCount     int
	RegisterCount int
}

type AgentForSelect struct {
	Id      int
	OrgName string
}

//Excel数据解析
type ExcelData struct {
	Date           time.Time //日期
	RegisterCount  int       //注册人数
	ApplyCount     int       //认证完成人数
	CreditCount    int       //授信人数
	ApplyLoanCount int       //申请借款人数
	MakeLoanCount  int       //放款人数
	MakeLoanAmount float64   //放款金额
	ProPrice       float64   //收益
	JointMode      int       //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CpaDefine      string    //CPA结算的有效事件
	CpaPrice       float64   //CPA的价格
	CpsFirstPer    float64   //cps的首借百分比
	CpsAgainPer    float64   //cps的复借百分比
}

func (this *AgentDailyData) TableName() string {
	return "daily_data"
}

func QueryAgentsByProductId(productId int) (agents []AgentForSelect, err error) {
	sql := `SELECT a.id,a.org_name FROM agent a INNER JOIN agent_product ap ON a.id = ap.agent_id WHERE pro_id = ?`
	_, err = orm.NewOrm().Raw(sql, productId).QueryRows(&agents)
	return
}
func QueryAgentDailyDataById(id int, date string) (dailyData AgentDailyData, err error) {
	sql := `SELECT * FROM daily_data WHERE id = ? AND date=?  `
	err = orm.NewOrm().Raw(sql, id, date).QueryRow(&dailyData)
	return
}
func InsertAgentDailyData(data AgentDailyData, date string) error {
	sql := `INSERT INTO daily_data (agent_product_id,register_count,apply_count,pro_price,
	 credit_count,apply_loan_count,
	 make_loan_count,make_loan_amount,joint_mode,cpa_price,cps_first_per,cps_again_per,cpa_define,date)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := orm.NewOrm().Raw(sql, data.AgentProductId, data.RegisterCount, data.ApplyCount, data.ProPrice, data.CreditCount, data.ApplyLoanCount, data.MakeLoanCount,
		data.MakeLoanAmount,
		data.JointMode, data.CpaPrice, data.CpsFirstPer, data.CpsAgainPer, data.CpaDefine,
		date).Exec()
	return err
}

//根据产品id和代理商id获取中间表id
func GetAgentProduct(agentId, productId int) (agentProductId int, err error) {
	sql := `SELECT id FROM agent_product WHERE agent_id = ? AND pro_id = ?`
	err = orm.NewOrm().Raw(sql, agentId, productId).QueryRow(&agentProductId)
	return
}

//根据时间获取产品的注册人数、申请人数、授信人数、申请借款人数、放款人数、放款金额信息
func GetDailyDataByDate(agentProductId int, startDate string) (proDailyData *AgentDailyData, err error) {
	sql := `SELECT id,agent_product_id,register_count,apply_count,
	 credit_count,apply_loan_count,
	 make_loan_count,make_loan_amount,
	 joint_mode,cpa_price,cps_first_per,cps_again_per,cpa_define
	 FROM daily_data WHERE agent_product_id = ? AND date = ?`
	err = orm.NewOrm().Raw(sql, agentProductId, startDate).QueryRow(&proDailyData)
	return
}
func UpdateAgentDailyData(data AgentDailyData, date string) error {
	sql := `UPDATE daily_data SET register_count=?,
	apply_count=?,credit_count=?,apply_loan_count=?,make_loan_count=?,make_loan_amount=?,pro_price=?,
	joint_mode=?,cpa_price=?,cps_first_per=?,cps_again_per=?,cpa_define=?
	WHERE agent_product_id =? AND date =?`
	_, err := orm.NewOrm().Raw(sql, data.RegisterCount, data.ApplyCount, data.CreditCount, data.ApplyLoanCount, data.MakeLoanCount, data.MakeLoanAmount, data.ProPrice,
		data.JointMode, data.CpaPrice, data.CpsFirstPer, data.CpsAgainPer, data.CpaDefine,
		data.AgentProductId, date).Exec()
	return err
}

//SUPPLEMENT:代补充资料，REFUSE：审核未通过，WAITING：待审核，PASS：审核通过，PAYING：待放款，PAYFAILED:放款失败，CONFIRM：已放款，FINISH：已完成，OVERDUE：逾期中

//获取数据分析-产品列表
func GetAnalyisProduct(condition string) (product []ProductIdList, err error) {
	sql := `SELECT id,name,cooperation_type type FROM product WHERE id!=0 AND name !="微融" ` //cooperation_type 0:api  1:H5
	sql += condition
	_, err = orm.NewOrm().Raw(sql).QueryRows(&product)
	return
}

//获取产品放款次数
func GetLoan_count(condition string, params []string) (res []ProductCount, err error) {
	sql := `SELECT bl.product_id id,count(1) count FROM business_loan bl WHERE (bl.state='CONFIRM' OR bl.state='FINISH' OR bl.state='OVERDUE')`
	sql += condition + " GROUP BY bl.product_id"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品申请次数
func GetApplycount(condition string, params []string) (res []ProductCount, err error) {
	sql := `SELECT bl.product_id id,count(1) count FROM business_loan bl WHERE 1=1`
	sql += condition + " GROUP BY bl.product_id"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品放款金额
func GetLoanmoney(condition string, params []string) (res []ProductMoney, err error) {
	sql := `SELECT bl.product_id id ,Sum(bl.real_money) money FROM business_loan bl WHERE 1=1 `
	sql += condition + " GROUP BY bl.product_id HAVING bl.product_id IS NOT NULL"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品总支出金额
func GetAllloanmoney(condition string, params []string) (res []ProductMoney, err error) {
	sql := `  SELECT bl.product_id id,(COUNT(pru.product_id)*pc.cpa_price) money
	FROM business_loan bl
	LEFT JOIN product_cleaning pc
	ON bl.product_id=pc.product_id
	LEFT JOIN
	( SELECT product_id, COUNT(uid)
	FROM product_register_user
	GROUP BY product_id) pru
	ON bl.product_id=pru.product_id
	WHERE 1=1  `
	sql += condition + " GROUP BY bl.product_id"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品新老用户总放款人数
func GetAllloancount(condition string, params []string) (res []ProductCount, err error) {
	sql := `SELECT bl.product_id id ,count(1) count FROM business_loan bl WHERE
		(bl.state='CONFIRM' OR bl.state='FINISH' OR bl.state='OVERDUE')`
	sql += condition + " GROUP BY bl.product_id"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品信用
func GetCredit(condition string, params []string) (res []ProductMoney, err error) {
	sql := `SELECT bl.product_id id  ,SUM(um.credit_code) money FROM business_loan bl LEFT JOIN users_metadata um ON um.uid=bl.uid WHERE 1=1 `
	sql += condition + " GROUP BY bl.product_id"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取H5放款人数
func GetH5Loan_count(params []string) (res []ProductCount, err error) {
	sql := `SELECT a.pro_id id, SUM(d.make_loan_count) count
			FROM daily_data d
			LEFT JOIN agent_product a ON d.agent_product_id = a.id
			WHERE d.date>=? AND d.date<= ?
			GROUP BY a.pro_id
			HAVING pro_id IS NOT NULL `
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取H5申请人数
func GetH5Applycount(params []string) (res []ProductCount, err error) {
	sql := `SELECT a.pro_id id, SUM(d.apply_loan_count) count
			FROM daily_data d
			LEFT JOIN agent_product a ON d.agent_product_id = a.id
			WHERE d.date>=? AND d.date<= ?
			GROUP BY a.pro_id
			HAVING pro_id IS NOT NULL  `
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取H5产品放款金额
func GetH5Loanmoney(params []string) (res []ProductMoney, err error) {
	sql := `SELECT agent_product.pro_id id,SUM(daily_data.make_loan_amount) money FROM
	daily_data LEFT JOIN agent_product ON daily_data.agent_product_id=agent_product.id WHERE
	 date>=? AND date<= ? GROUP BY agent_product.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取H5总支出金额
func GetH5Allloanmoney(params []string) (res []ProductMoney, err error) {
	sql := `select  ap.pro_id id,SUM(pc.cpa_price*register_count) money from daily_data dd
	left join agent_product ap
	on dd.agent_product_id = ap.id
 	left join product_cleaning pc
  	on ap.pro_id=pc.product_id WHERE
	date>=? AND date<= ? GROUP BY ap.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取H5新老用户总放款人数
func GetH5Allloancount(params []string) (res []ProductCount, err error) {
	sql := `SELECT agent_product.pro_id id,SUM(daily_data.make_loan_count) count FROM 
	daily_data LEFT JOIN agent_product ON daily_data.agent_product_id=agent_product.id WHERE
	 date>=? AND date<= ? GROUP BY agent_product.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品点击次数
func GetHits(params []string) (res []ProductCount, err error) {
	sql := `SELECT product_id id ,count(1) count FROM product_landing_page_record WHERE create_time>=? AND create_time<=? GROUP BY product_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品注册人数
func GetRegister(params []string) (res []ProductCount, err error) {
	sql := `SELECT product_id id ,count(1) count FROM product_register_user WHERE create_time>=? AND create_time<=? GROUP BY product_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取平台返回注册
func GetPlatformRegister(params []string) (res []ProductCount, err error) {
	sql := `SELECT agent_product.pro_id id,SUM(daily_data.register_count) count 
			FROM agent_product LEFT JOIN daily_data ON daily_data.agent_product_id=agent_product.id 
			WHERE date>=? AND date<= ? GROUP BY agent_product.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取产品按时间取放款人数
func GetLoan_countByTime(condition string, proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT DATE_FORMAT(bl.create_time,'%Y-%m-%d') data ,count(1) count FROM business_loan bl WHERE ( bl.state='CONFIRM' OR bl.state='FINISH' OR bl.state='OVERDUE') `
	sql += condition + " GROUP BY DATE_FORMAT(bl.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品按时间取申请人数
func GetApplycountByTime(condition string, proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT DATE_FORMAT(bl.create_time,'%Y-%m-%d') data,count(1) count FROM business_loan bl WHERE 1=1`
	sql += condition + " GROUP BY DATE_FORMAT(bl.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品按时间取放款金额
func GetLoanmoneyByTime(condition string, proid int, params []string) (res []ProductMoneyData, err error) {
	sql := `SELECT DATE_FORMAT(bl.create_time,"%Y-%m-%d") data, sum(bl.real_money) money
			FROM business_loan bl
			WHERE 1 = 1 `
	sql += condition + " GROUP BY DATE_FORMAT(bl.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品按时间取总支出金额
func GetAllloanmoneyByTime(condition string, proid int, params []string) (res []ProductMoneyData, err error) {
	sql := `SELECT DATE_FORMAT(bl.create_time,'%Y-%m-%d') data,(COUNT(pru.product_id)*pc.cpa_price) money
	FROM business_loan bl
	LEFT JOIN product_cleaning pc
	ON bl.product_id=pc.product_id
	LEFT JOIN
	( SELECT product_id, COUNT(uid)
	FROM product_register_user
	GROUP BY product_id) pru
	ON bl.product_id=pru.product_id
	 WHERE 1=1`
	sql += condition + " GROUP BY DATE_FORMAT(bl.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品按时间取新老用户总放款人数
func GetAllloancountByTime(condition string, proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT DATE_FORMAT(bl.create_time,'%Y-%m-%d') data ,count(1) count FROM business_loan bl WHERE
		(bl.state='CONFIRM' OR bl.state='FINISH' OR bl.state='OVERDUE')`
	sql += condition + " GROUP BY DATE_FORMAT(bl.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品按时间取信用
func GetCreditByTime(condition string, proid int, params []string) (res []ProductMoneyData, err error) {
	sql := `SELECT DATE_FORMAT(bl.create_time,'%Y-%m-%d') data ,Sum(um.credit_code) money FROM business_loan bl LEFT JOIN users_metadata um ON um.uid=bl.uid WHERE 1=1 `
	sql += condition + " GROUP BY DATE_FORMAT(bl.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品按时间取点击次数
func GetHitsByTime(proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT DATE_FORMAT(plpr.create_time,"%Y-%m-%d") data, SUM(plpr.count) count
			FROM product_landing_page_record plpr
			WHERE plpr.product_id = ? AND plpr.create_time>=? AND DATE_FORMAT(plpr.create_time,"%Y-%m-%d")<= ?
			GROUP BY DATE_FORMAT(plpr.create_time,"%Y-%m-%d") `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)

	return
}

//获取产品按时间取注册人数
func GetRegisterByTime(proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT DATE_FORMAT(pr.create_time,"%Y-%m-%d") data, COUNT(1) count
			FROM product_register_user pr
			WHERE pr.product_id = ? AND pr.create_time >= ?  AND DATE_FORMAT(pr.create_time,"%Y-%m-%d") <= ?
			GROUP BY DATE_FORMAT(pr.create_time,"%Y-%m-%d")  `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//按时间取平台返回注册
func GetPlatformRegisterByTime(proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT dd.date data,SUM(dd.register_count) count
			FROM daily_data dd INNER JOIN agent_product ap ON dd.agent_product_id=ap.id   WHERE ap.pro_id =?
			AND dd.date>=? AND dd.date<= ?
			GROUP BY dd.date `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//按时间取H5放款人数
func GetH5Loan_countByTime(proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT dd.date data, SUM(dd.make_loan_count) count
			FROM 
				daily_data dd
			INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
			WHERE ap.pro_id = ? AND DATE>=? AND DATE<= ?
			GROUP BY dd.date
			HAVING data IS NOT NULL  `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//按时间获取H5申请人数
func GetH5ApplycountByTime(proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT dd.date data, SUM(dd.apply_loan_count) count
			FROM 
				daily_data dd
			INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
			WHERE ap.pro_id = ? AND dd.date>=? AND dd.date<= ?
			GROUP BY dd.date
			HAVING data IS NOT NULL`
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//按时间获取H5产品放款金额
func GetH5LoanmoneyByTime(proid int, params []string) (res []ProductMoneyData, err error) {
	sql := `SELECT dd.date data, SUM(dd.make_loan_amount) money
			FROM daily_data dd
			INNER JOIN agent_product ap ON dd.agent_product_id = ap.id
			WHERE ap.pro_id = ? AND DATE>=? AND DATE<= ?
			GROUP BY dd.date `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//按时间获取H5总支出金额
func GetH5AllloanmoneyByTime(proid int, params []string) (res []ProductMoneyData, err error) {
	sql := `SELECT dd.date data, SUM(pc.cpa_price*register_count) money
			FROM daily_data dd
			LEFT JOIN agent_product ap ON dd.agent_product_id = ap.id
			LEFT JOIN product_cleaning pc ON ap.pro_id=pc.product_id
			WHERE ap.pro_id = ? AND DATE>=? AND DATE<= ?
			GROUP BY dd.date `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//按时间获取H5新老用户总放款人数
func GetH5AllloancountByTime(proid int, params []string) (res []ProductCountDate, err error) {
	sql := `SELECT dd.date data, SUM(dd.make_loan_count) count
			FROM 
				daily_data dd
			INNER JOIN agent_product ap ON dd.agent_product_id=ap.id
			WHERE ap.pro_id = ? AND dd.date>=? AND dd.date<= ?
			GROUP BY dd.date `
	_, err = orm.NewOrm().Raw(sql, proid, params).QueryRows(&res)
	return
}

//获取产品H5CPA
func GetH5CPA(id int) (res float32, err error) {
	err = orm.NewOrm().Raw("SELECT cpa_price FROM product_cleaning WHERE product_id = ?", id).QueryRow(&res)
	return
}

//查询日期判断数据是否存在
func FindDataFromExcel(agentProductId int) (date []time.Time, err error) {
	sql := `SELECT date FROM daily_data WHERE agent_product_id=? `
	_, err = orm.NewOrm().Raw(sql, agentProductId).QueryRows(&date)
	return
}

func UpdateDataFromExcel(dailyDatas []ExcelData, agentProductId int) error {
	o := orm.NewOrm()
	sql := `UPDATE daily_data SET register_count=?,
	apply_count=?,credit_count=?,apply_loan_count=?,make_loan_count=?,make_loan_amount=?,pro_price=?,
	joint_mode=?,cpa_price=?,cps_first_per=?,cps_again_per=?,cpa_define=?
	WHERE agent_product_id =? AND date =?`
	o = orm.NewOrm()
	dailysql, err := o.Raw(sql).Prepare()
	defer dailysql.Close()
	for k, v := range dailyDatas {
		switch v.JointMode {
		case 1:
			switch v.CpaDefine {
			case "注册":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].RegisterCount) * v.CpaPrice
			case "认证":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyCount) * v.CpaPrice
			case "授信":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].CreditCount) * v.CpaPrice
			case "申请借款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyLoanCount) * v.CpaPrice
			case "放款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].MakeLoanCount) * v.CpaPrice
			}
		case 2:
			dailyDatas[k].ProPrice = float64(dailyDatas[k].MakeLoanAmount) * v.CpsFirstPer / 100
		case 3:
			switch v.CpaDefine {
			case "注册":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].RegisterCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer/100
			case "认证":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer/100
			case "授信":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].CreditCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer/100
			case "申请借款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyLoanCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer/100
			case "放款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].MakeLoanCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer/100
			}
		}
		_, err = dailysql.Exec(v.RegisterCount, v.ApplyCount, v.CreditCount, v.ApplyLoanCount, v.MakeLoanCount, v.MakeLoanAmount, dailyDatas[k].ProPrice,
			v.JointMode, v.CpaPrice, v.CpsFirstPer, v.CpsAgainPer, v.CpaDefine,
			agentProductId, v.Date.Format("2006-01-02"))
	}
	return err
}

func InsertDataFromExcel(dailyDatas []ExcelData, agentProductId int) error {
	o := orm.NewOrm()
	sql := `INSERT INTO daily_data (agent_product_id,register_count,apply_count,
	 credit_count,apply_loan_count,
	 make_loan_count,make_loan_amount,date,pro_price,
	 joint_mode,cpa_price,cps_first_per,cps_again_per,cpa_define)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	dailysql, err := o.Raw(sql).Prepare()
	defer dailysql.Close()
	for k, v := range dailyDatas {
		switch v.JointMode {
		case 1:
			switch v.CpaDefine {
			case "注册":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].RegisterCount) * v.CpaPrice
			case "认证":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyCount) * v.CpaPrice
			case "授信":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].CreditCount) * v.CpaPrice
			case "申请借款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyLoanCount) * v.CpaPrice
			case "放款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].MakeLoanCount) * v.CpaPrice
			}
		case 2:
			dailyDatas[k].ProPrice = float64(dailyDatas[k].MakeLoanAmount) * v.CpsFirstPer
		case 3:
			switch v.CpaDefine {
			case "注册":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].RegisterCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer
			case "认证":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer
			case "授信":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].CreditCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer
			case "申请借款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].ApplyLoanCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer
			case "放款":
				dailyDatas[k].ProPrice = float64(dailyDatas[k].MakeLoanCount)*v.CpaPrice + float64(dailyDatas[k].MakeLoanAmount)*v.CpsFirstPer
			}
		}
		_, err = dailysql.Exec(agentProductId, v.RegisterCount, v.ApplyCount, v.CreditCount, v.ApplyLoanCount, v.MakeLoanCount,
			v.MakeLoanAmount, v.Date, dailyDatas[k].ProPrice, v.JointMode, v.CpaPrice, v.CpsFirstPer, v.CpsAgainPer, v.CpaDefine)
	}
	return err
}

//根据产品ID获取详细信息
func GetProImformation(id int) (cleaningData *AgentDailyData, err error) {
	sql := `SELECT pc.joint_mode,pc.cpa_define,pc.cpa_price,pc.cps_price
	      FROM product_cleaning pc
	      INNER JOIN product p
	      ON pc.product_id=p.id
          INNER JOIN agent_product ap
          ON  p.id=ap.pro_id
          WHERE ap.id=? LIMIT 1`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&cleaningData)
	return
}

// 根据productId 获取产品类型
func GetProductCooperationTypeById(id int) (cooperationType int, err error) {
	sql := `SELECT cooperation_type  FROM product WHERE id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&cooperationType)
	return
}
