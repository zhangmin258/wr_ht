package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CleanHistoryShow struct {
	Id             int
	ProId          int
	SettleTime     time.Time //结算时间
	BeginTime      time.Time //开始时间
	EndTime        time.Time //结束时间
	QuantityCount  int       //有效数量
	MakeLoanAmount float64   //放款金额
	SettleMoney    float64   //结算金额
	IsBilling      int64     //是否开票：0，没有；1，开票
	Remark         string    //备注
	DailyDataId    string
}

type CleanHistoryInfo struct {
	Id             int
	ProId          int
	SettleTime     time.Time //结算时间
	BeginTime      time.Time //开始时间
	EndTime        time.Time //结束时间
	QuantityCount  int       //有效数量
	MakeLoanAmount float64   //放款金额
	SettleMoney    float64   //结算金额
	IsBilling      int64     //是否开票：0，没有；1，开票
	Remark         string    //备注
	Name           string    //平台名称
	SelltlePrice   string    //结算价格
	CpaPrice       float64
	CpsPrice       float64
	JointMode      int64
	ProPrice       string
}

//开票信息和后台信息
type BillingInformation struct {
	Id             int
	ProductId      int    //产品id
	CompanyName    string //开票：公司名称
	CompanyAddress string //开票：公司地址
	TaxNumber      string //开票：税号
	BankAccount    string //开票：银行账号
	AccountBank    string //开票：开户银行
	BackUrl        string //后台：后台网址
	BackAccount    string //后台：后台账号
	BackPwd        string //后台:密码
	FromNum        string //后台:渠道号
}

//根据产品id获取结算历史记录
func GetCleanHisByPid(pid, begin, size int) (cleanHistory []*CleanHistoryShow, err error) {
	sql := `SELECT ch.id,ch.pro_id,ch.settle_time,ch.begin_time,ch.end_time,ch.quantity_count,ch.make_loan_amount,ch.settle_money,ch.daily_data_id,
	ch.is_billing,ch.remark
	FROM cleaning_history ch
	INNER JOIN product p
	ON ch.pro_id=p.id
	WHERE p.id = ? ORDER BY ch.end_time DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, pid, begin, size).QueryRows(&cleanHistory)
	return
}

//根据产品id获取结算历史总数
func CleanHisByPidCount(pid int) (count int, err error) {
	sql := `SELECT count(1)
	FROM cleaning_history ch
	INNER JOIN product p
	ON ch.pro_id=p.id
	WHERE p.id = ?`
	err = orm.NewOrm().Raw(sql, pid).QueryRow(&count)
	return
}

//按搜索时间获取结算信息
func GetCleanHisByTime(order, condition string, params []string, begin, size int) (cleanHistoryInfo []*CleanHistoryInfo, err error) {
	sql := `SELECT ch.id,ch.pro_id,ch.settle_time,ch.begin_time,ch.end_time,ch.quantity_count,ch.make_loan_amount,ch.settle_money,p.name,
	ch.is_billing,ch.remark,pc.joint_mode,pc.cpa_price,pc.cps_price
	FROM cleaning_history ch
	INNER JOIN product p
	ON ch.pro_id = p.id
	INNER JOIN product_cleaning pc
	ON pc.product_id = ch.pro_id
	WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	if order != "" {
		sql += order
	}
	sql += ` LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&cleanHistoryInfo)
	return
}

//已结算平台总数
func GetCleanHisCount(condition string, params []string) (count int, err error) {
	sql := `SELECT count(1)
	FROM cleaning_history ch
	WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//根据ID获得开票信息
func GetBilInformation(pid int) (billingInformation *BillingInformation, err error) {
	sql := `SELECT company_name,company_address,tax_number,back_account,account_bank,back_url,back_account,back_pwd,from_num
	      FROM product_cleaning
       	  WHERE product_id=?`
	err = orm.NewOrm().Raw(sql, pid).QueryRow(&billingInformation)
	return
}
