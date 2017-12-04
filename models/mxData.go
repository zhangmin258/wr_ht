package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type MdbMxSBData struct {
	Uid        int
	MXSBData   MXSBData
	CreateTime time.Time
}

//社保
type MXSBData struct {
	City     string `json:"city"`
	TaskID   string `json:"task_id"`
	AreaCode string `json:"area_code"`
	BaseInfo struct {
		UserInfoID              string `json:"user_info_id"`
		SocialSecurityNo        string `json:"social_security_no"`
		PersonalNo              string `json:"personal_no"`
		RealName                string `json:"real_name"`
		LastPayDate             string `json:"last_pay_date"`
		IDCard                  string `json:"id_card"`
		IDType                  string `json:"id_type"`
		Phone                   string `json:"phone"`
		BaseNumber              int    `json:"base_number"`
		Nation                  string `json:"nation"`
		BirthDay                string `json:"birth_day"`
		Sex                     string `json:"sex"`
		FirstInsuredDate        string `json:"first_insured_date"`
		PersonnelStatus         string `json:"personnel_status"`
		HouseholdRegistration   string `json:"household_registration"`
		InsuredUnit             string `json:"insured_unit"`
		InsuredUnitCode         string `json:"insured_unit_code"`
		IndustrialInsurance     string `json:"industrial_insurance"`
		UnemploymentInsurance   string `json:"unemployment_insurance"`
		MedicalInsurance        string `json:"medical_insurance"`
		EndowmentInsurance      string `json:"endowment_insurance"`
		MaternityInsurance      string `json:"maternity_insurance"`
		FetchTime               string `json:"fetch_time"`
		UnitType                string `json:"unit_type"`
		PayStatus               string `json:"pay_status"`
		BeginDate               string `json:"begin_date"`
		MedicalInsuranceBalance int    `json:"medical_insurance_balance"`
		WorkTime                string `json:"work_time"`
		Address                 string `json:"address"`
	} `json:"base_info"`
	Insurances []struct {
		InsuranceID           string `json:"insurance_id"`
		InsuranceType         string `json:"insurance_type"`
		InsuranceCode         int    `json:"insurance_code"`
		InsuranceStatus       int    `json:"insurance_status"`
		FirstInsuredDate      string `json:"first_insured_date"`
		ThisInsuredDate       string `json:"this_insured_date"`
		BaseNumber            int    `json:"base_number"`
		CorporationName       string `json:"corporation_name"`
		MonthlyCustomerIncome int    `json:"monthly_customer_income"`
	} `json:"insurances"`
	InsuranceRecord []struct {
		InsuranceType   string `json:"insurance_type"`
		InsuranceCode   int    `json:"insurance_code"`
		DealTime        string `json:"deal_time"`
		Month           string `json:"month"`
		CorporationName string `json:"corporation_name"`
		Description     string `json:"description"`
		BaseNumber      int    `json:"base_number"`
		PersonalPayment int    `json:"personal_payment"`
	} `json:"insurance_record"`
	MedicalInsuranceRecord []struct {
		OrganizationName string `json:"organization_name"`
		Type             string `json:"type"`
		SettlemenTime    string `json:"settlemen_time"`
		Money            int    `json:"money"`
	} `json:"medical_insurance_record"`
}

type MdbMxGJJData struct {
	Uid        int
	MXGJJData  MXGJJData
	CreateTime time.Time
}

//公积金
type MXGJJData struct {
	City     string `json:"city"`
	TaskID   string `json:"task_id"`
	AreaCode string `json:"area_code"`
	UserInfo struct {
		Gender                   string `json:"gender"`
		Birthday                 string `json:"birthday"`
		Balance                  int    `json:"balance"`
		Mobile                   string `json:"mobile"`
		Email                    string `json:"email"`
		CustomerNumber           string `json:"customer_number"`
		RealName                 string `json:"real_name"`
		PayStatus                string `json:"pay_status"`
		IDCard                   string `json:"id_card"`
		CardType                 string `json:"card_type"`
		HomeAddress              string `json:"home_address"`
		CorporationName          string `json:"corporation_name"`
		MonthlyCorporationIncome int    `json:"monthly_corporation_income"`
		MonthlyCustomerIncome    int    `json:"monthly_customer_income"`
		MonthlyTotalIncome       int    `json:"monthly_total_income"`
		LastPayDate              string `json:"last_pay_date"`
		FundBalance              string `json:"fund_balance"`
		SubsidyBalance           int    `json:"subsidy_balance"`
		CorporationNumber        string `json:"corporation_number"`
		CorporationRatio         string `json:"corporation_ratio"`
		CustomerRatio            string `json:"customer_ratio"`
		SubsidyCustomerRatio     string `json:"subsidy_customer_ratio"`
		SubsidyCorporationRatio  string `json:"subsidy_corporation_ratio"`
		BaseNumber               int    `json:"base_number"`
		BeginDate                string `json:"begin_date"`
		GjjNumber                string `json:"gjj_number"`
		SubsidyIncome            int    `json:"subsidy_income"`
	} `json:"user_info"`
	BillRecord []struct {
		Month       string `json:"month"`
		Income      int    `json:"income"`
		Outcome     int    `json:"outcome"`
		Description string `json:"description"`
		Balance     int    `json:"balance"`
		DealTime    string `json:"deal_time"`
	} `json:"bill_record"`
	LoanInfo        []interface{} `json:"loan_info"`
	LoanRepayRecord []interface{} `json:"loan_repay_record"`
}

//从数据库查询用户魔蝎社保公积金信息
func GetUsersMXSBGJJInfo(uid int, task_type string) (data string, createTime time.Time, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT response,creat_time FROM mx_joint_record WHERE uid = ? AND joint_record_type=5 AND task_type=? ORDER BY creat_time DESC LIMIT 1`
	err = o.Raw(sql, uid, task_type).QueryRow(&data, &createTime)
	return
}
