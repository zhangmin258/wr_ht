package models

import (
	"github.com/astaxie/beego/orm"
)

//-----------------------------------------------/**有盾返回数据**/--------------------------------------------------
type YdUserResult struct {
	Uid      int
	YdResult YdResult
}
type YdResult struct {
	Result    Result
	Data      Data   //数据
	Signature string //签名，用于商户进行验签，详见2.4商户验签说明
}
type Result struct {
	Success bool   //结果状态
	Message string //结果描述
}
type Data struct {
	Smses             []SmsItems
	Code              string           //0：正常，1：改月份没有数据，2：运营商网络错误，数据采集失败，3：没有权限查询数据，4：数据采集内容与期望结果不一致，9：未知
	City              string           //城市
	Level             string           //账号星级
	Nets              []NetsItems      //流量详情字段说明
	Last_Modify_Time  string           //最近一次更新时间，格式： yyyy-MM-dd HH:mm:ss
	Mobile            string           //手机号
	Open_Time         string           //入网时间，格式：yyyy-MM-dd
	Message           string           //状态描述
	Packages          []PackagesItems  //套餐字段说明
	Families          []FamiliesItems  //亲情网字段说明
	Available_Balance string           //当前可用余额（单位：分）
	Recharges         []RechargesItems //充值记录字段说明
	Carrier           string           //账号归属运营商标识:CHINA_MOBILE：中国移动 CHINA_TELECOM：中国电信 CHINA_UNICOM：中国联通
	Province          string           //省份
	Calls             []CallsItems     //语音详情字段说明
	Idcard            string           //证件号
	Name              string           //姓名
	Package_Name      string           //套餐名称
	Bills             []BillItems      //账单信息字段说明
	State             string           //账号状态:-1：未知，0：正常，1：单向停机，2：停机，3：预销户，4：销户，5：过户，6：改号，99：号码不存在
}

//-----------------------------------------------/**有盾返回数据--短信详情字段说明**/--------------------------------------------------
type SmsItem struct {
	Peer_Number  string
	Details_Id   string
	Send_Type    string
	Service_Name string
	Fee          string
	Msg_Type     string
	Location     string
	Time         string
}
type SmsItems struct {
	Bill_Month string
	Total_Size string
	Items      []SmsItem
}

//-----------------------------------------------/**有盾返回数据--流量详情字段说明**/--------------------------------------------------
type NetItem struct {
	Details_Id   string
	Duration     string
	Subflow      string
	Net_Type     string
	Service_Name string
	Fee          string
	Location     string
	Time         string
}
type NetsItems struct {
	Bill_Month string
	Total_Size string
	Items      []NetItem
}

//-----------------------------------------------/**有盾返回数据--套餐详情字段说明**/--------------------------------------------------
type PackageItem struct {
	Item  string
	Total string
	Unit  string
	Used  string
}
type PackagesItems struct {
	Bill_Start_Date string
	Items           []PackageItem
	Bill_End_Date   string
}

//-----------------------------------------------/**有盾返回数据--亲情网详情字段说明**/--------------------------------------------------
type FamiliesItems struct {
	Family_Num string
	Items      []FamiliesItem
}
type FamiliesItem struct {
	Long_Number  string // 长号
	Short_Number string //短号
	Member_Type  string //成员类型
	Join_Date    string //加入日期
	Expire_Date  string //失效日期
}

//-----------------------------------------------/**有盾返回数据--充值详情字段说明**/--------------------------------------------------
type RechargesItems struct {
	Details_Id    string
	Amount        string
	Recharge_Time string
	Type          string
}

//-----------------------------------------------/**有盾返回数据--语音详情字段说明**/--------------------------------------------------
type CallsItem struct {
	Peer_Number   string
	Duration      string
	Details_Id    string
	Dial_Type     string
	Fee           string
	Location      string
	Time          string
	Location_Type string
}
type CallsItems struct {
	Bill_Month string
	Total_Size string
	Items      []CallsItem
}

//-----------------------------------------------/**有盾返回数据--账单信息字段说明**/--------------------------------------------------
type BillItems struct {
	Bill_Start_Date   string
	Notes             string
	Related_Mobiles   string
	Discount          string
	Paid_Fee          string
	Web_Fee           string
	Point             string
	Sms_Fee           string
	Base_Fee          string
	Extra_Discount    string
	Bill_Month        string
	Extra_Service_Fee string
	Total_Fee         string
	Last_Point        string
	Extra_Fee         string
	Actual_Fee        string
	Bill_End_Date     string
	Voice_Fee         string
	Unpaid_Fee        string
}

//-----------------------------------------------/**mysql存储**/--------------------------------------------------

type YdUserSms struct {
	Id          int `orm:"column(id);auto"`
	Uid         int
	Name        string //姓名
	BillMonth   string //详情月份，格式：yyyy-MM
	PeerNumber  string //对方号码
	DetailsId   string //详情标识
	SendType    string //SEND-发送，RECEIVE-收取	items.send_type	N	String	SEND-发送，RECEIVE-收取
	ServiceName string //业务名称，eg：点对点（网内）
	Fee         string //通话费（单位：分）
	MsgType     string //SMS-短信，MMS-彩信
	Location    string //通话地（自己的）
	Time        string //收/发短信时间
}

type YdUserNet struct {
	Id          int `orm:"column(id);auto"`
	Uid         int
	Name        string //姓名
	BillMonth   string //详情月份，格式：yyyy-MM
	DetailsId   string //详情标识
	Duration    string //流量使用时长
	Subflow     string //流量使用量，单位：KB
	NetType     string //网络类型
	ServiceName string //业务名称，eg：点对点（网内）
	Fee         string //通话费（单位：分）
	Location    string //通话地（自己的）
	Time        string //收/发短信时间
}

type YdUserFamily struct {
	Id          int `orm:"column(id);auto"`
	Uid         int
	Name        string //姓名
	FamilyNum   string //亲情网编号
	LongNumber  string // 长号
	ShortNumber string //短号
	MemberType  string //成员类型
	JoinDate    string //加入日期
	ExpireDate  string //失效日期
}

type YdUserPackage struct {
	Id            int `orm:"column(id);auto"`
	Uid           int
	Name          string //姓名
	BillStartDate string //账单起始日，格式：yyyy-MM-dd
	BillEndDate   string //账单结束日，格式：yyyy-MM-dd
	Item          string //套餐项目名称
	Total         string //项目总量
	Unit          string //项目已使用量
	Used          string //单位：语音-分；流量-KB；短/彩信-条
}

type YdUserRecharge struct {
	Id           int `orm:"column(id);auto"`
	Uid          int
	Name         string //姓名
	DetailsId    string //详情标识
	Amount       string //充值金额（单位：分）
	RechargeTime string //充值时间，格式：yyyy-MM-dd HH:mm:ss
	Type         string //充值方式，eg：现金
}

type YdUserCall struct {
	Id           int `orm:"column(id);auto"`
	Uid          int
	Name         string //姓名
	BillMonth    string //详情月份，格式yyyy-MM
	PeerNumber   string //对方号码
	Duration     string //通话时长（单位：秒）
	DetailsId    string //详情标识
	DialType     string //DIAL-主叫，DIALED-被叫
	Fee          string //通话费（单位：分）
	Location     string //通话地（自己的）
	Time         string //收/发短信时间
	LocationType string //通话地类型，eg：省内漫游
}

type YdUserBill struct {
	Id              int `orm:"column(id);auto"`
	Uid             int
	Name            string //姓名
	BillStartDate   string //账单起始日期，格式：yyyy-MM-dd
	Notes           string //备注
	RelatedMobiles  string //本手机关联号码，多个手机号以逗号分隔
	Discount        string //优惠费
	PaidFee         string //本期已付费用
	WebFee          string //网络流量费
	Point           string //本期可用积分
	SmsFee          string //短彩信费
	BaseFee         string //套餐及固定费
	ExtraDiscount   string //其他优惠
	BillMonth       string //账单月，格式：yyyy-MM
	ExtraServiceFee string //增值业务费
	TotalFee        string //总费用
	LastPoint       string //上期可用积分
	ExtraFee        string //其他费用
	ActualFee       string //个人实际费用
	BillEndDate     string //账单结束日期，格式：yyyy-MM-dd
	VoiceFee        string //语音费
	UnpaidFee       string //本期未付费用
}
type YdUserBaseData struct {
	Id               int
	Uid              int
	Code             string //0：正常，1：改月份没有数据，2：运营商网络错误，数据采集失败，3：没有权限查询数据，4：数据采集内容与期望结果不一致，9：未知
	City             string //城市
	Level            string //账号星级
	LastModifyTime   string //最近一次更新时间，格式： yyyy-MM-dd HH:mm:ss
	Mobile           string //手机号
	OpenTime         string //入网时间，格式：yyyy-MM-dd
	Message          string //状态描述
	AvailableBalance string //当前可用余额（单位：分）
	Carrier          string //账号归属运营商标识:CHINA_MOBILE：中国移动 CHINA_TELECOM：中国电信 CHINA_UNICOM：中国联通
	Province         string //省份
	Idcard           string //证件号
	Name             string //姓名
	PackageName      string //套餐名称
	State            string //账号状态:-1：未知，0：正常，1：单向停机，2：停机，3：预销户，4：销户，5：过户，6：改号，99：号码不存在
}

type YDYYSUdcreditDevice_detail struct {
	App_instalment_count string
	Is_rooted            string
	Cheats_device        string
	Is_using_proxy_port  string
	Network_type         string
}

type YDYYSUdcreditUser_features struct {
	User_feature_type  string
	Last_modified_date string
}

type YDYYSUdcreditId_detail struct {
	Birthday string
	Address  string
	Gender   string
	Nation   string
}

type YDYYSUdcreditMobile_detail struct {
	Province string
	City     string
	Isp      string
}

type MdbYDYYSUdcreditData struct {
	Uid                     string
	YDYYSUdcreditYdudcredit Ydudcredit `json:"ydyysudcreditydudcredit"`
}

func UpdateYdUdCredit(uid, result int) error {
	o := orm.NewOrm()
	sql := `UPDATE users_auth SET is_ydudcredit = ?,ydudcredit_time=now() WHERE uid = ?`
	_, err := o.Raw(sql, result, uid).Exec()
	return err
}


func AddYdUserSms(ydUserSmses []YdUserSms) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserSmses), ydUserSmses)
	return err
}

func AddYdUserNet(ydUserNet []YdUserNet) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserNet), ydUserNet)
	return err
}

func AddYdUserPackage(ydUserPackage []YdUserPackage) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserPackage), ydUserPackage)
	return err
}

func AddYdUserRecharge(ydUserRecharge []YdUserRecharge) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserRecharge), ydUserRecharge)
	return err
}

func AddYdUserCall(ydUserCall []YdUserCall) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserCall), ydUserCall)
	return err
}

func AddYdUserBill(ydUserBill []YdUserBill) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserBill), ydUserBill)
	return err
}

func AddYdUserFamily(ydUserFamily []YdUserFamily) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.InsertMulti(len(ydUserFamily), ydUserFamily)
	return err
}

func AddYdUserBaseData(data *YdUserBaseData) error {
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err := o.Insert(data)
	return err
}

type YDResult struct { //成功返回
	Success bool   `json:"success"` //结果状态
	Message string `json:"message"` //结果描述
}

//社保数据查询
type MdbYDSBData struct {
	Uid                 string
	YDSBGetDataResponse YDSBGetDataResponse
}

type YDSBGetDataResponse struct {
	Result    YDResult `json:"result"`    //结果
	Data      YDSBData `json:"data"`      //数据
	Signature string   `json:"signature"` //签名，用于商户进行验签
}

type YDSBData struct {
	TaskId                 string                       `json:"task_id"`                  //任务id
	City                   string                       `json:"city"`                     //所属城市
	AreaCode               string                       `json:"area_code"`                //地区编码
	BaseInfo               YDSBBaseInfo                 `json:"base_info"`                //基本信息
	Insurances             []YDSBInsurances             `json:"insurances"`               //保险种类，返回集合(多个险种)
	InsuranceRecord        []YDSBInsuranceRecord        `json:"insurance_record"`         //保险缴存记录(多个险种缴存记录)
	MedicalInsuranceRecord []YDSBMedicalInsuranceRecord `json:"medical_insurance_record"` //医保卡消费记录
}

type YDSBInsurances struct {
	BaseNumber               string `json:"base_number"`                //该保险缴存基数，单位：分
	CorporationName          string `json:"corporation_name"`           //缴存公司名称
	CorporationScale         string `json:"corporation_scale"`          //公司缴存比例
	CustomerScale            string `json:"customer_scale"`             //个人缴存比例
	Description              string `json:"description"`                //描述信息
	FirstInsuredDate         string `json:"first_insured_date"`         //首次参保时间
	ThisInsuredDate          string `json:"this_insured_date"`          //本次参保时间
	InsuranceCode            string `json:"insurance_code"`             //险种编号
	InsuranceStatus          string `json:"insurance_status"`           //参保状态
	InsuranceId              string `json:"insurance_id"`               //保险id(险种)
	InsuranceType            string `json:"insurance_type"`             //保险类型(险种)
	MonthlyCorporationIncome string `json:"monthly_corporation_income"` //公司缴存金额，单位为分
	MonthlyCustomerIncome    string `json:"monthly_customer_income"`    //个人缴存金额，单位为分
	TotalMonths              string `json:"total_months"`               //缴存月数
}
type YDSBInsuranceRecord struct {
	Amount             string `json:"amount"`              //缴存总额，单位为分
	BaseNumber         string `json:"base_number"`         //缴存基数，单位为分
	CorporationName    string `json:"corporation_name"`    //缴存公司名称
	CorporationPayment string `json:"corporation_payment"` //公司缴存金额，单位为分
	CorporationScale   string `json:"corporation_scale"`   //公司缴存比例
	PersonalPayment    string `json:"personal_payment"`    //个人缴存数
	CustomerScale      string `json:"customer_scale"`      //个人缴存比例
	DealTime           string `json:"deal_time"`           //发生时间
	Description        string `json:"description"`         //描述信息
	InsuranceType      string `json:"insurance_type"`      //所属保险
	InsuranceCode      string `json:"insurance_code"`      //险种编号:1：养老保险 2：医疗保险 3：工伤保险 4：生育保险 5：失业保险 6：大病保险 0：其他保险
	Month              string `json:"month"`               //缴存月份
	MonthEnd           string `json:"month_end"`           //缴存月份(一次缴存多个月的，表示最后范围的月份)
	Status             string `json:"status"`              //缴存状态标记
}
type YDSBMedicalInsuranceRecord struct {
	Money            string `json:"money"`             //金额
	OrganizationName string `json:"organization_name"` //医疗机构名称
	SettlemenTime    string `json:"settlemen_time"`    //结算时间
	Type             string `json:"type"`              //医疗类别(门诊,药店，住院)
}
type YDSBBaseInfo struct {
	RealName                string `json:"real_name"`                 //真实姓名
	UserInfoId              string `json:"user_info_id"`              //用户信息Id
	SocialSecurityNo        string `json:"social_security_no"`        //社保编号/社保卡号/医保卡号号
	PersonalNo              string `json:"personal_no"`               //个人编号
	Nation                  string `json:"nation"`                    //民族
	IdType                  string `json:"id_type"`                   //证件类型(默认指身份证)
	IdCard                  string `json:"id_card"`                   //证件号码(默认指身份证)
	Address                 string `json:"address"`                   //家庭住址
	BaseNumber              string `json:"base_number"`               //缴存基数(通常指养老保险基础)
	BeginDate               string `json:"begin_date"`                //开户日期或者指开始参加社保的日期
	BirthDay                string `json:"birth_day"`                 //出生日期
	FirstInsuredDate        string `json:"first_insured_date"`        //首次参保时间
	HouseholdRegistration   string `json:"household_registration"`    //户籍性质
	InsuredUnit             string `json:"insured_unit"`              //参保单位
	InsuredUnitCode         string `json:"insured_unit_code"`         //参保单位编号
	LastPayDate             string `json:"last_pay_date"`             //最后缴纳记录时间
	PayStatus               string `json:"pay_status"`                //缴存状态(是否停缴)
	PersonnelStatus         string `json:"personnel_status"`          //人员状态(在职,离职等)
	Phone                   string `json:"phone"`                     //手机号、电话
	Sex                     string `json:"sex"`                       //性别
	UnitType                string `json:"unit_type"`                 //单位类型
	WorkTime                string `json:"work_time"`                 //参加工作时间
	IndustrialInsurance     string `json:"industrial_insurance"`      //工伤保险
	UnemploymentInsurance   string `json:"unemployment_insurance"`    //失业保险
	MedicalInsurance        string `json:"medical_insurance"`         //医疗保险
	MedicalInsuranceBalance string `json:"medical_insurance_balance"` //医疗保险
	EndowmentInsurance      string `json:"endowment_insurance"`       //养老保险
	MaternityInsurance      string `json:"maternity_insurance"`       //生育保险
	FetchTime               string `json:"fetch_time"`                //采集时间
}

//有盾公积金
type MdbYDGJJData struct {
	Uid                  string
	YDGJJGetDataResponse YDGJJGetDataResponse
}

type YDGJJGetDataResponse struct {
	Result    YDResult  `json:"result"`    //结果
	Data      YDGJJData `json:"data"`      //数据
	Signature string    `json:"signature"` //签名，用于商户进行验签
}

type YDGJJData struct {
	TaskId          string                 `json:"task_id"`           //医保卡消费记录
	UserInfo        YDGJJUserInfo          `json:"user_info"`         //公积金基本信息
	BillRecord      []YDGJJBillRecord      `json:"bill_record"`       //公积金缴存记录（按年份查询）
	LoanInfo        []YDGJJloanInfo        `json:"loan_info"`         //公积金贷款信息
	LoanRepayRecord []YDGJJLoanRepayRecord `json:"loan_repay_record"` //贷款还贷记录
}

type YDGJJloanInfo struct {
	Name                     string `json:"name"`                       //贷款人姓名
	Phone                    string `json:"phone"`                      //贷款-联系手机
	Status                   string `json:"status"`                     //贷款状态
	Bank                     string `json:"bank"`                       //承办银行
	LoanType                 string `json:"loan_type"`                  //贷款类型（公积金贷款/商业贷款/组合贷款）
	IdCard                   string `json:"id_card"`                    //贷款人身份证
	HouseAddress             string `json:"house_address"`              //当前贷款购房地址
	MailingAddress           string `json:"mailing_address"`            //通讯地址
	ContractNumber           string `json:"contract_number"`            //贷款合同号
	Periods                  string `json:"periods"`                    //贷款期数
	LoanAmount               string `json:"loan_amount"`                //贷款金额
	MonthlyRepayAmount       string `json:"monthly_repay_amount"`       //月还款额度
	StartDate                string `json:"start_date"`                 //贷款开始时间
	EndDate                  string `json:"end_date"`                   //贷款结束日期
	RepayType                string `json:"repay_type"`                 //还款方式（等额本金）
	DeductDay                string `json:"deduct_day"`                 //每月还款日
	BankAccount              string `json:"bank_account"`               //扣款账号
	BankAccountName          string `json:"bank_account_name"`          //扣款银行账号姓名
	LoanInterestPercent      string `json:"loan_interest_percent"`      //贷款利率
	PenaltyInterestPercent   string `json:"penalty_interest_percent"`   //罚息利率
	CommercialContractNumber string `json:"commercial_contract_number"` //商业贷款合同编号
	CommercialBank           string `json:"commercial_bank"`            //商业贷款银行
	CommercialAmount         string `json:"commercial_amount"`          //商业贷款金额
	SecondBankAccount        string `json:"second_bank_account"`        //第二还款人银行账号
	SecondBankAccountName    string `json:"second_bank_account_name"`   //第二还款人姓名
	SecondIdCard             string `json:"second_id_card"`             //第二还款人身份证
	SecondPhone              string `json:"second_phone"`               //第二还款人手机
	SecondCorporationName    string `json:"second_corporation_name"`    //第二还款人工作单位
	RemainAmount             string `json:"remain_amount"`              //贷款余额
	RemainPeriods            string `json:"remain_periods"`             //剩余期数
	LastRepayDate            string `json:"last_repay_date"`            //最后还款日期
	OverdueCapital           string `json:"overdue_capital"`            //逾期本金
	OverdueInterest          string `json:"overdue_interest"`           //逾期利息
	OverduePenalty           string `json:"overdue_penalty"`            //逾期罚息
	OverdueDays              string `json:"overdue_days"`               //逾期天数
}
type YDGJJUserInfo struct {
	RealName                 string `json:"real_name"`                  //姓名
	Gender                   string `json:"gender"`                     //性别   1：男   0：女   未知，此字段不显示（性别根据身份证推算出来）
	Birthday                 string `json:"birthday"`                   //出生年月
	Mobile                   string `json:"mobile"`                     //手机号码
	Email                    string `json:"email"`                      //邮箱
	CustomerNumber           string `json:"customer_number"`            //客户号
	GjjNumber                string `json:"gjj_number"`                 //公积金账号
	Balance                  string `json:"balance"`                    //账户余额（包含公积金余额跟补贴余额）
	FundBalance              string `json:"fund_balance"`               //公积金余额
	SubsidyBalance           string `json:"subsidy_balance"`            //补贴公积金余额（补贴公积金）
	SubsidyIncome            string `json:"subsidy_income"`             //补贴月缴存
	PayStatus                string `json:"pay_status"`                 //缴存状态:  NONE：未缴纳    NORMAL：正常    SUSPENSE：停缴  CLOSED：注销
	IdCard                   string `json:"id_card"`                    //身份证号码
	CardType                 string `json:"card_type"`                  //证件类型    ID#ARD：身份证  PASSPORT：护照
	HomeAddress              string `json:"home_address"`               //通讯地址
	CorporationNumber        string `json:"corporation_number"`         //企业账户号码
	CorporationName          string `json:"corporation_name"`           //当前缴存企业名称
	MonthlyCorporationIncome string `json:"monthly_corporation_income"` //企业月度缴存
	MonthlyCustomerIncome    string `json:"monthly_customer_income"`    //个人月度缴存
	MonthlyTotalIncome       string `json:"monthly_total_income"`       //月度总缴存
	CorporationRatio         string `json:"corporation_ratio"`          //企业缴存比例
	CustomerRatio            string `json:"customer_ratio"`             //个人缴存比例
	SubsidyCorporationRatio  string `json:"subsidy_corporation_ratio"`  //补贴公积金公司缴存比例
	SubsidyCustomerRatio     string `json:"subsidy_customer_ratio"`     //补贴公积金个人缴存比例
	BaseNumber               string `json:"base_number"`                //缴存基数
	LastPayDate              string `json:"last_pay_date"`              //最新缴存日期
	BeginDate                string `json:"begin_date"`                 //开户日期
}
type YDGJJBillRecord struct {
	Outcome           string `json:"outcome"`            //出账
	Income            string `json:"income"`             //入账
	SubsidyIncome     string `json:"subsidy_income"`     //补贴入账
	Description       string `json:"description"`        //描述
	SubsidyOutcome    string `json:"subsidy_outcome"`    //补贴出账
	Balance           string `json:"balance"`            //余额
	DealTime          string `json:"deal_time"`          //缴存时间
	Month             string `json:"month"`              //缴存年月
	CorporationName   string `json:"corporation_name"`   //缴存公司名称
	CorporationIncome string `json:"corporation_income"` //公司缴存金额
	CustomerIncome    string `json:"customer_income"`    //个人缴存金额
	CorporationRatio  string `json:"corporation_ratio"`  //公司缴存比例
	CustomerRatio     string `json:"customer_ratio"`     //个人缴存比例
	AdditionalIncome  string `json:"additional_income"`  //补缴
}
type YDGJJLoanRepayRecord struct {
	RepayDate      string `json:"repay_date"`      //还款日期
	AccountingDate string `json:"accounting_date"` //记账日期
	RepayAmount    string `json:"repay_amount"`    //还款金额
	RepayCapital   string `json:"repay_capital"`   //还款本金
	RepayInterest  string `json:"repay_interest"`  //还款利息
	RepayPenalty   string `json:"repay_penalty"`   //还款罚息
	ContractNumber string `json:"contract_number"` //贷款合同号
}

//返回的社保信息
type UserSecurityData struct {
	City                            string //所属城市
	AreaCode                        string //地区编码
	SocialSecurityNo                string //社保卡号
	PersonalNo                      string //个人编号
	RealName                        string //姓名
	BaseNumber                      string //缴存基数
	LastPayDate                     string //最后缴存日期
	IdType                          string //证件类型
	IdCard                          string //证件号码
	Nation                          string //民族
	BirthDay                        string //生日
	Phone                           string //手机号码
	Address                         string //家庭住址
	PersonnelStatus                 string //人员状态
	HouseholdRegistration           string //户口属性
	FirstInsuredDate                string //首次参保时间
	WorkTime                        string //参加工作时间
	InsuredUnit                     string //参保单位
	InsuredUnitCode                 string //参保单位编码
	UnitType                        string //单位类型
	PayStatus                       string //缴存状态
	BeginDate                       string //开户日期
	IndustrialInsuranceBaseNumber   string //工伤保险基数
	UnemploymentInsuranceBaseNumber string //失业保险基数
	MedicalInsuranceBaseNumber      string //医疗保险基数
	EndowmentInsuranceBaseNumber    string //养老保险基数
	MaternityInsuranceBaseNumber    string //生育保险基数
	FetchTime                       string //抓取时间
	MedicalInsuranceBalance         string //医疗保险账户余额
}

//从数据库查询用户社保公积金信息
func GetUsersYdInfo(uid int, taskType string) (data string, err error) {
	o := orm.NewOrm()
	sql := `SELECT response FROM youdun_joint_record WHERE joint_record_type=2 AND task_type=? AND success=0 AND uid = ? ORDER BY creat_time DESC LIMIT 1`
	err = o.Raw(sql, taskType, uid).QueryRow(&data)
	return
}

func GetUsersTotalCount() (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM users`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//获取社保公积金认证用户数量
func GetSecurityFundUsers() (fundCount, securityCount, securityFundCount, totalCount int, err error) {
	o := orm.NewOrm()
	//社保
	sql1 := `SELECT COUNT(1) FROM users_auth WHERE is_sb_mx = 2 `
	err = o.Raw(sql1).QueryRow(&fundCount)
	if err != nil {
		return
	}
	//公积金
	sql2 := `SELECT COUNT(1) FROM users_auth WHERE is_gjj_mx = 2 `
	err = o.Raw(sql2).QueryRow(&securityCount)
	if err != nil {
		return
	}
	//社保公积金
	sql3 := `SELECT COUNT(1) FROM users_auth WHERE is_sb_mx = 2 AND is_gjj_mx = 2 `
	err = o.Raw(sql3).QueryRow(&securityFundCount)
	if err != nil {
		return
	}
	//总人数
	sql4 := `SELECT COUNT(1) FROM users `
	err = o.Raw(sql4).QueryRow(&totalCount)
	if err != nil {
		return
	}
	return
}

type UserSecurityFund struct {
	Num  int    `json:"Count"`
	Type string `json:"Data"`
}

func SecurityFundInit() (s []UserSecurityFund) {
	t1 := UserSecurityFund{Type: "security"}
	t2 := UserSecurityFund{Type: "fund"}
	t3 := UserSecurityFund{Type: "security_fund"}
	t4 := UserSecurityFund{Type: "total_count"}
	s = []UserSecurityFund{t1, t2, t3, t4}
	return
}
