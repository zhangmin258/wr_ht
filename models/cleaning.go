package models

import (
	//"fmt"

	"time"

	"github.com/astaxie/beego/orm"
)

//结算信息
type ProductCleaning struct {
	Id             int       `orm:"pk"` //id
	ProductId      int       //产品id
	JointMode      int       //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CleaningType   int       //结算方式：1，对公；2，对私
	CpaDefine      string    //CPA结算的有效事件
	CpaPrice       float64   //CPA的价格
	CpsPrice       float64   //CPA的价格
	CleaningCycle  string    //结算周期
	CompanyName    string    //开票：公司名称
	CompanyAddress string    //开票：公司地址
	TaxNumber      string    //开票：税号
	BankAccount    string    //开票：银行账号
	AccountBank    string    //开票：开户银行
	BackUrl        string    //后台：后台网址
	BackAccount    string    //后台：后台账号
	BackPwd        string    //后台:密码
	FromNum        string    //后台:渠道号
	CreateTime     time.Time //创建时间
	CpsFirstPer    float64   //S的首借百分比
	CpsAgainPer    float64   //S的复借百分比
}

// 保存清算信息数据
func AddCleaning(cleaning *ProductCleaning) (err error) {
	sql := `INSERT INTO product_cleaning (product_id,joint_mode,cleaning_type,cpa_define,cpa_price,
	cps_price,cleaning_cycle,company_name,company_address,
	tax_number,bank_account,account_bank,back_url,
	back_account,back_pwd,from_num,cps_first_per,cps_again_per,create_time)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`
	_, err = orm.NewOrm().Raw(sql, cleaning.ProductId, cleaning.JointMode, cleaning.CleaningType, cleaning.CpaDefine, cleaning.CpaPrice,
		cleaning.CpsPrice, cleaning.CleaningCycle, cleaning.CompanyName, cleaning.CompanyAddress, cleaning.TaxNumber,
		cleaning.BankAccount, cleaning.AccountBank, cleaning.BackUrl, cleaning.BackAccount,
		cleaning.BackPwd, cleaning.FromNum, cleaning.CpsFirstPer, cleaning.CpsAgainPer).Exec()
	return
}

// 根据产品的id查询清算信息数据
func SearchCleaningById(proId int) (claening ProductCleaning, err error) {
	claening.ProductId = proId
	err = orm.NewOrm().Read(&claening, "ProductId")
	return
}

//根据id更新清算数据
func UpdateCleaningById(cleaning *ProductCleaning) (err error) {

	var time time.Time
	orm.NewOrm().Raw(`SELECT create_time FROM product_cleaning WHERE id=?`, cleaning.Id).QueryRow(&time)
	cleaning.CreateTime = time
	_, err = orm.NewOrm().Update(&cleaning)
	return
}

//根据id更新清算数据 ---H5
func UpdateCleaningH5(cleaning *ProductCleaning) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE product_cleaning SET joint_mode=?,cleaning_type=?,cpa_define=?,cpa_price=?,cps_first_per=?,
		cps_again_per=?,cleaning_cycle=?,company_name=?,company_address=?,tax_number=?,account_bank=?,bank_account=? ,back_url=?,back_account=?,back_pwd=?,from_num=?
		WHERE id = ?`
	_, err = o.Raw(sql, cleaning.JointMode, cleaning.CleaningType, cleaning.CpaDefine, cleaning.CpaPrice, cleaning.CpsFirstPer,
		cleaning.CpsAgainPer, cleaning.CleaningCycle, cleaning.CompanyName, cleaning.CompanyAddress, cleaning.TaxNumber,
		cleaning.AccountBank, cleaning.BankAccount, cleaning.BackUrl, cleaning.BackAccount, cleaning.BackPwd, cleaning.FromNum, cleaning.Id).Exec()
	return
}

//根据id更新清算数据 ---API
func UpdateCleaningAPI(cleaning *ProductCleaning) (err error) {
	//fmt.Println(cleaning)
	o := orm.NewOrm()
	sql := `UPDATE product_cleaning SET joint_mode=?,cleaning_type=?,cpa_define=?,cpa_price=?,cps_first_per=?,cps_again_per=?,
	cleaning_cycle=?,company_name=?,company_address=?,tax_number=?,account_bank=?,bank_account=?,back_url=?,back_account=?,back_pwd=?,from_num=?
	   WHERE id = ?`
	_, err = o.Raw(sql, cleaning.JointMode, cleaning.CleaningType, cleaning.CpaDefine, cleaning.CpaPrice, cleaning.CpsFirstPer, cleaning.CpsAgainPer,
		cleaning.CleaningCycle, cleaning.CompanyName, cleaning.CompanyAddress, cleaning.TaxNumber,
		cleaning.AccountBank, cleaning.BankAccount, cleaning.BackUrl, cleaning.BackAccount, cleaning.BackPwd, cleaning.BackPwd, cleaning.Id).Exec()
	return
}
