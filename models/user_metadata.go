package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//用户资信信息
type UsersMetadata struct {
	Id               int       `orm:"column(id);pk"` //pk
	Uid              int       //用户ID
	LoanAmount       float64   //贷款金额
	LoanTerm         int       //贷款期限
	IsProvidentFund  int8      //是否有本地公积金，1：否，没有2：是有本地公积金
	IsSocialSecurity int8      //是否有本地社保,1:否，没有，2:是，有本地社保
	HouseProperty    string    //名下房产
	CarProperty      int8      //名下车产:1:无车;2:无车，准备购买;3:名下有车
	CreditSituation  int8      //信用情况:1:一年内逾期超过3次或超过90天;2:一年内逾期少于3次且少于90天;3:无信用卡或贷款;4:信用良好，无逾期
	MaxRepayment     float64   //可接受最高月还款额度
	SalaryForm       int8      //工资发放形式
	BankAmount       float64   //银行发放工资
	CashAmount       float64   //现金发放工资
	CashSalary       float64   //现金收入
	WorkingAge       int       //工龄
	Turnover         float64   //总流水
	CashTurnover     float64   //现金流水
	ManageAge        string    //经营年限
	ManageLocation   string    //经营地点
	BusinessLicense  string    //是否有营业执照
	Identity         int8      //职业身份
	VerifyRealName   string    //真实姓名
	IdCard           string    //身份证号码
	Account          string    //账号
	Birthday         string    //生日
	IsBindCard       int8      //是否绑卡
	Location         string    //定位到的城市
	Address          string    //定位到的详细地址
	IpLocation       string    //IP定位到的城市
	IpAddress        string    //IP定位到的详细地址
	LongiAndLati     string    //经纬度
	Sex              string    //性别
	IsRealName       int8      //是否实名认证
	RealNameTime     time.Time //实名认证时间
	CreateTime       time.Time //创建时间
	ModifyTime       time.Time //修改时间
	ZmScore          int       //芝麻信用分
	CreditCard       int       //信用卡情况
	Carrieroperator  int       //手机运营商认证:0默认值，1:未认证。2已认证
	BqsMobile        int       //默认:0,授权通过:1,数据获取成功:2,报告数据:3,运营商数据获取失败:9
	BqsAntiFraud     int       //反欺诈0:默认,1:无风险,2:低风险,3:高风险
}

//用户个人资信信息
type UserCreditInfo struct {
	LoanAmount       float64 //贷款金额
	LoanTerm         int     //贷款期限
	IsProvidentFund  int     //是否有本地公积金，1：否，没有2：是有本地公积金
	IsSocialSecurity int     //是否有本地社保,1:否，没有，2:是，有本地社保
	HouseProperty    string  //名下房产
	CarProperty      int     //名下车产:1:无车;2:无车，准备购买;3:名下有车
	CreditSituation  int     //信用情况:1:一年内逾期超过3次或超过90天;2:一年内逾期少于3次且少于90天;3:无信用卡或贷款;4:信用良好，无逾期
	IsCreditCard     int     //是否有信用卡：1：无信用卡，2：有，持卡未满6个月，3：有，持卡超过6个月
}

func GetUsersMetadataById(uid int) (v *UsersMetadata, err error) {
	sql := `SELECT um.*,zm.error_code FROM users_metadata um LEFT JOIN zmxy_auth zm on um.uid=zm.uid where um.uid = ?`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&v)
	return
}

func GetUserCredit(uid int) (userCredit UserCreditInfo, err error) {
	sql := `SELECT loan_amount, loan_term, is_provident_fund, is_social_security, house_property, car_property, credit_situation, is_credit_card 
			FROM users_basedata 
			WHERE uid = ?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&userCredit)
	return
}
