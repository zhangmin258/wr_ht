package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	//"strconv"
	"strings"
	"time"
	"zcm_tools/uuid"
)

//添加产品
type ProductForAdd struct {
	Id      int    `orm:"pk"`                //产品ID
	ProName string `orm:"column(name);null"` //产品名称

	OrgId                int       //机构ID
	BusId                int       //商务联系人ID
	ProductType          int       //产品类型 小贷，银行贷
	ProductLogo          string    //logo
	LoanTermCount        string    //期数1.3.9----------手动拼接
	LoanTermCountType    string    `orm:"-"` //期数单位"天""月"
	LoanSpeed            string    //放款速度
	MinMoney             float64   //最小额度
	MaxMoney             float64   //最大额度
	MinLoanTerm          int       //最小贷款期限
	MaxLoanTerm          int       //最大贷款期限
	LoanTermUnit         int       //贷款期限单位 0：天1：月
	LoanDailyFee         float64   //日利率
	LoanTaxFee           float64   //月利率
	LoanFee              float64   `orm:"-"` //利率-------------------从页面接收利率数据，根据单位手动给日利率或者月利率赋值
	MinFee               float64   //最小利率
	MaxFee               float64   //最大利率
	FeeType              int       //利率类型:	1.固定值，2:范围
	FeeMethod            int       //利息方式:	1.月利率，2:日利率
	LoanServiceBeforeFee float64   //贷前服务费
	ServiceFeeType       int       //贷前服务费类型	1:固定值2：范围
	MaxServiceFee        float64   //贷前服务费最大值
	MinServiceFee        float64   //贷前服务费最xiao值
	ManagementFee        float64   //管理费
	MaxManagementFee     float64   //管理费最大值
	MinManagementFee     float64   //管理费最小值
	ManagementFeeType    int       //管理费类型	1:固定值2：范围
	ManagementFeeUnit    int       //管理费单位	1:日2：月
	FeePayTimeType       int       //1:贷前扣款2：贷后扣款
	InterestRateDetails  string    //利率详情
	PaymentMethod        string    //还款方式
	Prepayment           string    //提前还款
	LatePolicy           string    //逾期政策
	IsUse                int8      //是否使用 0使用,1不使用------------1
	CooperationType      int       //合作方式 0:api 1:h5
	CreateTime           time.Time //创建时间
	IsIndexShow          int       //展示位置 0:贷款页 1:首页 2:大额贷款页
	Sort                 int       //首页排序
	LoanSort             int       //贷款页排序
	LargeLoanSort        int       //大额贷款页排序
	LoanProductType      int       //0:小额贷款1：大额贷款
	NonrecurringExpense  float64   //一次性费用
	InterestRatesDetails string    //利率详情（大额贷款）
	RequestedMaterial    string    //所需材料
	LargeLoanFee         float64   //服务费（大额贷款）
	RequestCondition     string    //申请条件
	Code                 string    //产品code
	ReferenceRate        float64   //大额参考利率
}

type ProductForEdit struct {
	Id                   int     `orm:"pk"`                //产品ID
	ProName              string  `orm:"column(name);null"` //产品名称
	ProductType          int     //产品类型 小贷，银行贷
	ProductLogo          string  //logo
	LoanTermCount        string  //期数1.3.9----------手动拼接
	LoanTermCountType    string  `orm:"-"` //期数单位"天""月"
	LoanSpeed            string  //放款速度
	MinMoney             float64 //最小额度
	MaxMoney             float64 //最大额度
	MinLoanTerm          int     //最小贷款期限
	MaxLoanTerm          int     //最大贷款期限
	LoanTermUnit         int     //贷款期限单位 0：天1：月
	LoanDailyFee         float64 //日利率
	LoanTaxFee           float64 //月利率
	LoanFee              float64 `orm:"-"` //利率-------------------从页面接收利率数据，根据单位手动给日利率或者月利率赋值
	MinFee               float64 //最小利率
	MaxFee               float64 //最大利率
	FeeType              int     //利率类型:	1.固定值，2:范围
	FeeMethod            int     //利息方式:	1.月利率，2:日利率
	LoanServiceBeforeFee float64 //贷前服务费
	ServiceFeeType       int     //贷前服务费类型	1:固定值2：范围
	MaxServiceFee        float64 //贷前服务费最大值
	MinServiceFee        float64 //贷前服务费最小值
	ManagementFee        float64 //管理费
	MaxManagementFee     float64 //管理费最大值
	MinManagementFee     float64 //管理费最小值
	ManagementFeeType    int     //管理费类型	1:固定值2：范围
	ManagementFeeUnit    int     //管理费单位	1:日2：月
	FeePayTimeType       int     //1:贷前扣款2：贷后扣款
	InterestRateDetails  string  //利率详情
	PaymentMethod        string  //还款方式
	Prepayment           string  //提前还款
	LatePolicy           string  //逾期政策
}
type LargeProductForEdit struct {
	Id                   int     `orm:"pk"`                //产品ID
	ProName              string  `orm:"column(name);null"` //产品名称
	ProductType          int     //产品类型 --默认为大额
	ProductLogo          string  //logo
	LoanTermCount        string  //期数1.3.9----------手动拼接
	LoanTermCountType    string  `orm:"-"` //期数单位"天""月"
	LoanSpeed            string  //放款速度
	MinMoney             float64 //最小额度
	MaxMoney             float64 //最大额度
	MinLoanTerm          int     //最小贷款期限
	MaxLoanTerm          int     //最大贷款期限
	LoanTermUnit         int     //贷款期限单位 0：天1：月---大额产品默认为月
	LoanTaxFee           float64 //月利率
	FeeMethod            int     //利息方式:	1.月利率，2:日利率-------大额，设置为月
	LargeLoanFee         float64 //服务费（大额贷款）
	NonrecurringExpense  float64 //一次性费用
	InterestRatesDetails string  //利率详情（大额贷款）
	RequestCondition     string  //申请条件
	RequestedMaterial    string  //所需材料
}

//产品列表展示
type ProductListShow struct {
	Id              int       //产品ID
	OrgName         string    //机构名称
	ProductType     int       //产品类型 小贷，银行贷
	Name            string    //产品名称
	MinMoney        float64   //最小金额
	MaxMoney        float64   //最大金额
	MinLoanTerm     int       //最小贷款期限
	MaxLoanTerm     int       //最大贷款期限
	LoanTermUnit    int       //贷款期限单位 0：天1：月
	CooperationType int       //合作类型0:Api,1:H5
	SuccessRate     float64   //成功率
	IsUse           int8      //是否使用 0使用,1不使用
	CreateTime      time.Time //创建时间
	Sort            int       //排序
}

//产品运营数据展示
type ProductMange struct {
	Id            int    //产品ID
	Name          string //产品名称
	IsPopUp       int    //是否在推荐弹窗展示：0，不展示1，展示
	IsUse         int    //是否使用: 0使用,1不使用
	Sort          int    //首页排序
	FullGuide     int    //导流满员人数
	IsIndexShow   int    //是否首页展示：0，不展示 1：展示 2:大额贷款页展示
	LoanSort      int    //贷款页排序
	RecommendSort int    //产品推广
	LargeLoanSort int    //大额贷款页排序
	Tag           string //产品标签
}

//下级代理关联产品：id ,name
type ProductAgent struct {
	Id              int
	Name            string
	CooperationType int
}

//product 排序
type ProductSort struct {
	Id       int //产品id
	Sort     int //排序
	LoanSort int //贷款页排序
}

//产品及相应代理商
type ProductAndAgent struct {
	AgentProductId  int
	AgentId         int    //代理商Id
	CooperationType int    //产品类型
	ProName         string `orm:"column(name);null"` //产品名称
	OrgName         string //代理商名称
}

//产品管理->编辑页面->申请条件
type ProductAddress struct {
	PAId             int       `orm:"column(id)"` //主键
	Id               int       //产品ID
	MinAge           int       //最小年龄
	MaxAge           int       //最大年龄
	MinZmScore       int       //芝麻分最低要求
	AddressType      int       //地域要求类型
	Address          string    //地域
	AppPermit        int       //手机系统类型
	IsProvidentFund  int64     //是否有本地公积金，1：否，没有2：是有本地公积金
	IsSocialSecurity int64     //是否有本地社保,1:否，没有，2:是，有本地社保
	HouseProperty    int64     //名下房产：1无房，2有房
	CarProperty      int64     //名下车产:1:无车;2:无车，准备购买;3:名下有车
	CreateTime       time.Time //创建时间
}

//申请条件展示
type FindProAddress struct {
	Id               int    // 主键id
	MinAge           int    //最小年龄
	MaxAge           int    //最大年龄
	MinZmscore       int    //芝麻分最低要求
	AddressType      int    //地域要求类型:1:只贷给，2:不贷给
	Address          string //地域，用,号隔开
	AppPermit        int    //1 ios,2 android,3 wx,4pc5 wd 6 wap
	IsProvidentFund  int64  //是否有本地公积金，1：否，没有2：是有本地公积金
	IsSocialSecurity int64  //是否有本地社保,1:否，没有，2:是，有本地社保
	HouseProperty    int64  //名下房产：1无房，2有房
	CarProperty      int64  //名下车产:1:无车;2:无车，准备购买;3:名下有车
}

//产品认证信息
type ProductAuth struct {
	ProductId int `description:"产品Id"`
	AuthField int `description:"认证项目Id"`
}

func (p *ProductForAdd) TableName() string {
	return "product"
}

// 产品保存时，关联的产品工作表的初始化
func AddProductJob(proId int) (err error) {
	o := orm.NewOrm()
	sql := `INSERT INTO product_job(product_id,job_id) VALUES(?,?)`
	prepare, err := o.Raw(sql).Prepare()
	defer prepare.Close()
	if err != nil {
		return
	}
	for i := 1; i < 6; i++ {
		_, err = prepare.Exec(proId, i)
		if err != nil {
			return
		}
	}
	return
}

//默认绑定本公司代理商
func BindAgent(proId int, agentId int) (err error) {
	sql := `INSERT INTO agent_product
	        (pro_id,agent_id,agent_time)
	        VALUES
	        (?,?,now())`
	_, err = orm.NewOrm().Raw(sql, proId, agentId).Exec()
	return
}

// 保存小额产品数据
func AddProductNew(product *ProductForAdd) (pid int, err error) {
	//处理产品期数：将贷款期限的数字和单位拼接
	unit := product.LoanTermCountType
	if strings.Compare(unit, "月") == 0 {
		product.LoanTermCountType = "个" + product.LoanTermCountType
	}
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType
	//处理贷款利率：判断利率单位是日还是月
	if product.FeeMethod == 1 { //月，赋值给月利率
		product.LoanTaxFee = product.LoanFee
	}
	if product.FeeMethod == 2 { //日，赋值给日利率
		product.LoanDailyFee = product.LoanFee
	}
	//默认冻结
	product.IsUse = 1
	//创建时间
	product.CreateTime = time.Now()
	//新增的产品默认为贷款页展示，顺序为最后
	loanSort, err := GetMaxLoanSort()
	product.LoanSort = loanSort + 1
	//首页,大额贷款页排序设为-1
	product.Sort = -1
	product.LargeLoanSort = -1
	//设置随机code
	product.Code = uuid.NewUUID().Hex()

	//插入数据
	o := orm.NewOrm()
	id, err := o.Insert(product)
	if err != nil {
		return -1, err
	}
	return int(id), err
}

// 保存大额API产品数据
func AddBigAPIProductNew(product *ProductForAdd) (pid int, err error) {
	product.LoanTermCountType = "个月"
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType
	product.LoanTermUnit = 1
	//创建时间
	product.CreateTime = time.Now()
	//新增的大额产品为大额贷款页展示，顺序为最后
	product.IsIndexShow = 2
	maxLargeLoanSort, err := GetMaxLargeLoanSort()
	product.LargeLoanSort = maxLargeLoanSort + 1
	//首页,贷款页排序设为-1
	product.Sort = -1
	product.LoanSort = -1
	//设置随机code
	product.Code = uuid.NewUUID().Hex()
	product.FeeMethod = 1
	product.FeeType = 1
	//插入数据
	o := orm.NewOrm()
	id, err := o.Insert(product)
	if err != nil {
		return -1, err
	}
	return int(id), err
}

// 保存大额产品数据
func AddBigProductNew(product *ProductForAdd) (id int, err error) {
	product.LoanTermCountType = "个" + product.LoanTermCountType
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType
	//新增的大额产品为大额贷款页展示，顺序为最后
	product.IsIndexShow = 2
	maxLargeLoanSort, err := GetMaxLargeLoanSort()
	product.LargeLoanSort = maxLargeLoanSort + 1
	//首页,贷款页排序设为-1
	product.Sort = -1
	product.LoanSort = -1
	product.LoanTermUnit = 1
	product.IsUse = 1
	product.FeeType = 1
	product.FeeMethod = 1
	//设置随机code
	product.Code = uuid.NewUUID().Hex()
	sql := `INSERT INTO product (org_id,bus_id,name,loan_product_type,loan_term_unit,product_logo,loan_term_count,loan_speed,
	min_money,max_money,min_loan_term,max_loan_term,reference_rate,large_loan_fee,nonrecurring_expense,interest_rates_details,
	request_condition,requested_material,sort,loan_sort,cooperation_type,is_use,code,loan_tax_fee,fee_type,fee_method,create_time)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`
	result, err := orm.NewOrm().Raw(sql, product.OrgId, product.BusId, product.ProName, product.LoanProductType, product.LoanTermUnit, product.ProductLogo, product.LoanTermCount,
		product.LoanSpeed, product.MinMoney, product.MaxMoney, product.MinLoanTerm, product.MaxLoanTerm,
		product.ReferenceRate, product.LargeLoanFee, product.NonrecurringExpense, product.InterestRatesDetails,
		product.RequestCondition, product.RequestedMaterial, product.Sort, product.LoanSort, product.CooperationType, product.IsUse, product.Code, product.LoanTaxFee, product.FeeType, product.FeeMethod).Exec()
	ids, err := result.LastInsertId()
	id = int(ids)
	return
}

// 分页查询产品数据
func GetProductList(condition string, params []string, begin, size int) (products []ProductListShow, err error) {
	sql := `SELECT
			p.id, o.name AS org_name, p.name, p.product_type, p.is_use, p.min_money, p.max_money,
			p.min_loan_term, p.max_loan_term, p.success_rate, p.create_time, p.cooperation_type, p.sort,p.loan_term_unit
			FROM product p LEFT JOIN organization o ON o.id = p.org_id
			WHERE p.id != 0 `
	if condition != "" {
		sql += condition
	}
	sql += " order by p.id DESC limit ?, ?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&products)
	if err != nil {
		return nil, err
	}
	return
}

func GetProductListCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) count
			FROM product p LEFT JOIN organization o ON o.id = p.org_id
			WHERE p.id != 0 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 查询所有商品数量
func GetProductCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product p LEFT JOIN organization o ON o.id = p.org_id WHERE p.id != 0 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//查询使用产品数量
func GetUsedProductCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM product WHERE is_use=0 AND id!=0`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//查询未使用产品数量
func GetNotUsedProductCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM product WHERE is_use=1 AND id!=0`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//分页查询产品运营信息
func GetProductMange(order, condition string, params []string, begin, size int) (products []ProductMange, err error) {
	sql := `SELECT p.id,p.name,p.is_pop_up,p.is_use,p.sort,p.loan_sort,p.large_loan_sort,p.full_guide,p.is_index_show from product p
	LEFT JOIN organization o
	ON p.org_id=o.id
	WHERE p.id !=0 `
	if condition != "" {
		sql += condition
	}
	if order != "" {
		sql += order
	}
	sql += " limit ?,?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&products)
	return
}

//产品管理->编辑页面->数据保存
func UpdateProductMange(productmange *ProductMange, address *ProductAddress, sortList []ProductSort, loanSortList []ProductSort) error {
	o := orm.NewOrm()
	o.Begin()
	//修改sort后 更新改变后的sort值
	sql := `UPDATE product SET sort=?
	        WHERE id = ?`
	prepare, err := o.Raw(sql).Prepare()
	for _, v := range sortList {
		_, err = prepare.Exec(v.Sort, v.Id)
	}
	//修改loansort后 更新改变后的loansort值
	sql = `UPDATE product SET loan_sort=?
	        WHERE id = ?`
	prepare, err = o.Raw(sql).Prepare()
	for _, v := range loanSortList {
		_, err = prepare.Exec(v.LoanSort, v.Id)
	}
	//更新product数据
	sql = `UPDATE product SET is_use=?,sort=?,full_guide=?,is_pop_up=?,is_index_show=?,loan_sort=?,recommend_sort=?
	        WHERE id = ?`
	_, err = o.Raw(sql, productmange.IsUse, productmange.Sort, productmange.FullGuide, productmange.IsPopUp, productmange.IsIndexShow,
		productmange.LoanSort, productmange.RecommendSort, productmange.Id).Exec()
	//更新address数据
	if address.PAId == 0 {
		sql = `INSERT INTO product_loan_condition (product_id,min_age,max_age,min_zmscore,address_type,address,app_permit,
		is_provident_fund,is_social_security,
	 house_property,car_property,create_time)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`
		_, err = o.Raw(sql, address.Id, address.MinAge, address.MaxAge, address.MinZmScore, address.AddressType, address.Address,
			address.AppPermit, address.IsProvidentFund, address.IsSocialSecurity, address.HouseProperty, address.CarProperty,
			time.Now()).Exec()
	} else {
		sql = `UPDATE product_loan_condition SET min_age=?,max_age=?,min_zmscore=?,address_type=?,address=?,app_permit=?,modify_time=?,
		is_provident_fund=?,is_social_security=?,
	 house_property=?,car_property=?
		WHERE id=?`
		_, err = o.Raw(sql, address.MinAge, address.MaxAge, address.MinZmScore, address.AddressType, address.Address, address.AppPermit, time.Now(),
			address.IsProvidentFund, address.IsSocialSecurity, address.HouseProperty, address.CarProperty,
			address.PAId).Exec()
	}
	defer func() {
		if err != nil {
			o.Rollback()
			return
		}
		prepare.Close()
		o.Commit()
	}()
	return err
}

//根据id查找productaddress
func ProductAddressById(id int) (findproaddress FindProAddress, err error) {
	sql := `SELECT id,min_age,max_age,min_zmscore,address_type,address,app_permit,is_provident_fund,is_social_security,
	 house_property,car_property
	 FROM product_loan_condition WHERE product_id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&findproaddress)
	return
}

//根据id查找productmange
func ProductMangeById(id int) (product *ProductMange, err error) {
	sql := `SELECT id,is_use,sort,full_guide,loan_sort,is_pop_up,is_index_show,recommend_sort,large_loan_sort,tag from product WHERE id =?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&product)
	return

}

//根据id删除产品信息
/*func DeleteProduct(pid int) (err error) {
	sql := `DELETE FROM product where id=?`
	_, err = orm.NewOrm().Raw(sql, pid).Exec()
	return
}*/

// 根据id查找商品信息
func SearchProductById(id int) (product ProductForAdd, err error) {
	product = ProductForAdd{Id: id}
	err = orm.NewOrm().Read(&product)
	return
}

// 更新产品信息
func UpdateProduct(product *ProductForAdd) (err error) {
	//处理产品期数：将贷款期限的数字和单位拼接
	unit := product.LoanTermCountType
	if strings.Compare(unit, "月") == 0 {
		product.LoanTermCountType = "个" + product.LoanTermCountType
	}
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType
	//处理贷款利率：判断利率单位是日还是月
	if product.FeeMethod == 1 { //月，赋值给月利率
		product.LoanTaxFee = product.LoanFee
	}
	if product.FeeMethod == 2 { //日，赋值给日利率
		product.LoanDailyFee = product.LoanFee
	}
	var time time.Time
	orm.NewOrm().Raw(`SELECT create_time FROM product WHERE id=?`, product.Id).QueryRow(&time)
	product.CreateTime = time
	_, err = orm.NewOrm().Update(product)
	return
}

// 更新产品信息——方法2
func UpdateProductNew(product *ProductForAdd) (err error) {
	//处理产品期数：将贷款期限的数字和单位拼接
	unit := product.LoanTermCountType
	if strings.Compare(unit, "月") == 0 {
		product.LoanTermCountType = "个" + product.LoanTermCountType
	}
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType
	//处理贷款利率：判断利率单位是日还是月
	if product.FeeMethod == 1 { //月，赋值给月利率
		product.LoanTaxFee = product.LoanFee
	}
	if product.FeeMethod == 2 { //日，赋值给日利率
		product.LoanDailyFee = product.LoanFee
	}
	_, err = orm.NewOrm().Update(product, "name", "product_type", "product_logo", "loan_term_count", "loan_speed",
		"min_money", "max_money", "min_loan_term", "max_loan_term", "loan_term_unit", "min_service_fee", "max_service_fee",
		"service_fee_type", "management_fee", "management_fee_type", "min_management_fee", "max_management_fee", "management_fee_unit",
		"interest_rate_details", "payment_method", "prepayment", "late_policy", "fee_pay_time_type")
	return
}

// 更新产品信息
func UpdateBigProduct(product *ProductForAdd) (err error) {
	//处理产品期数：将贷款期限的数字和单位拼接
	product.LoanTermCountType = "个" + product.LoanTermCountType
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType
	sql := `UPDATE product SET name=?,loan_product_type=?,product_logo=?,loan_term_count=?,loan_speed=?,
	min_money=?,max_money=?,min_loan_term=?,max_loan_term=?,reference_rate=?,large_loan_fee=?,nonrecurring_expense=?,interest_rates_details=?,
	request_condition=?,requested_material=?,loan_tax_fee=? WHERE id = ?`
	_, err = orm.NewOrm().Raw(sql, product.ProName, product.LoanProductType, product.ProductLogo, product.LoanTermCount,
		product.LoanSpeed, product.MinMoney, product.MaxMoney, product.MinLoanTerm, product.MaxLoanTerm,
		product.ReferenceRate, product.LargeLoanFee, product.NonrecurringExpense, product.InterestRatesDetails,
		product.RequestCondition, product.RequestedMaterial, product.LoanTaxFee, product.Id).Exec()
	return
}

//获取所有产品的名称和id
func GetProductIdAndName(condition string, name string) (products []ProductAgent, err error) {
	sql := `SELECT id,name FROM product WHERE 1=1 `
	sql += condition
	_, err = orm.NewOrm().Raw(sql, name).QueryRows(&products)
	return
}

//获取所有H5产品的名称和id
func GetProductH5IdAndName(condition string, name string) (products []ProductAgent, err error) {
	sql := `SELECT id,name FROM product WHERE 1=1 AND cooperation_type = 1 AND id != 0 AND name !="微融"`
	sql += condition
	_, err = orm.NewOrm().Raw(sql, name).QueryRows(&products)
	return
}

//获取第一个H5产品的名称和id
func GetProductH5IdAndNameFirst() (product ProductAgent, err error) {
	sql := `SELECT id,name,cooperation_type FROM product WHERE 1=1 AND cooperation_type = 1 AND id!=0 AND name !="微融" LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&product)
	return
}

//获取数据分析-放款人数
func GetAnalyisLoanCount(params []string) (count int) {
	sql := `SELECT count(1) FROM `
	fmt.Println(sql)
	return
}

//获取数据分析-产品名称
func GetProductNameById(pid int) (productName string, err error) {
	sql := `SELECT name FROM product WHERE id=?`
	err = orm.NewOrm().Raw(sql, pid).QueryRow(&productName)
	/**
	获取所有产品及代理商信息
	*/ /*
		func GetProducts() (productAndAgents []ProductAndAgent, err error) {
			sql := `SELECT dd.agent_product_id,ap.agent_id,p.cooperation_type,p.name,a.org_name  FROM product p
					LEFT JOIN agent_product ap ON p.id = ap.pro_id
					LEFT JOIN agent a ON ap.agent_id = a.id
					LEFT JOIN daily_data dd ON  ap.id = dd.agent_product_id`
			_, err = orm.NewOrm().Raw(sql).QueryRows(&productAndAgents)
			return
		}*/
	return
}

type ProductForSelect struct {
	Id              int    //产品ID
	ProName         string `orm:"column(name)"` //产品名称
	CooperationType int    //产品类型
}

/**
获取所有H5产品ID和名称
*/
func GetProductsInfo(condition, params string) (products []ProductForSelect, err error) {
	sql := `SELECT id,name,cooperation_type FROM product WHERE id >0 `
	sql += condition
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&products)
	return
}

/**
根据产品ID查询产品类型
*/
func GetProductCooperationById(id int) (cooperationType int, err error) {
	sql := `SELECT cooperation_type  FROM product WHERE id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&cooperationType)
	return
}

/**
根据产品ID获取第一个代理商ID
*/
func GetFirstProductAgentIdByProductId(id int) (agentId int, err error) {
	sql := `SELECT a.id FROM agent a INNER JOIN agent_product ap ON a.id = ap.agent_id WHERE pro_id = ? LIMIT 0,1`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&agentId)
	return
}
func GetProductIdByName(productName string) (id int, err error) {
	sql := `SELECT id FROM product WHERE name = ?`
	err = orm.NewOrm().Raw(sql, productName).QueryRow(&id)
	return
}

//获取要改变排序的list
func SelectProductSort(condition string, params []int) (proSort []ProductSort, err error) {
	sql := `SELECT id,sort,loan_sort FROM product WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&proSort)
	return
}

func GetProductMangeById(id int) (pro ProductMange, err error) {
	sql := `SELECT id,sort,loan_sort FROM product WHERE id=? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&pro)
	return
}

//获取最大sort值
func GetMaxSort() (maxSort, maxLoanSort, maxLargeLoanSort int, err error) {
	sql := `SELECT MAX(sort),MAX(loan_sort),MAX(large_loan_sort) FROM product`
	err = orm.NewOrm().Raw(sql).QueryRow(&maxSort, &maxLoanSort, &maxLargeLoanSort)
	return
}

//获取使用中产品的最大sort值
func GetMaxUsedSortAndLoanSort() (maxSort, maxLoanSort, maxLargeLoanSort int, err error) {
	sql := `SELECT MAX(sort),MAX(loan_sort),MAX(large_loan_sort) FROM product WHERE is_use=0`
	err = orm.NewOrm().Raw(sql).QueryRow(&maxSort, &maxLoanSort, &maxLargeLoanSort)
	return
}

//获取最大sort和loanSort
func GetMaxLoanSort() (loanSort int, err error) {
	o := orm.NewOrm()
	sql := `SELECT loan_sort FROM product ORDER BY loan_sort DESC LIMIT 0,1`
	err = o.Raw(sql).QueryRow(&loanSort)
	return
}

//获取最大largeLoanSort
func GetMaxLargeLoanSort() (largeLoanSort int, err error) {
	o := orm.NewOrm()
	sql := `SELECT large_loan_sort FROM product ORDER BY large_loan_sort DESC LIMIT 0,1`
	err = o.Raw(sql).QueryRow(&largeLoanSort)
	return
}

func GetProductSortInfoById(id int) (pro ProductMange, err error) {
	sql := `SELECT id,sort,loan_sort,large_loan_sort,is_index_show,is_use FROM product WHERE id=? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&pro)
	o := orm.NewOrm()
	o.Begin()
	o.Rollback()
	o.Commit()

	return
}

//降位排序
func DownSort(oldSort, newSort int, types string) (err error) {
	var sql string
	if types == "index" {
		sql = `UPDATE product p SET p.sort=p.sort-1 WHERE p.sort>? AND p.sort<=? AND p.sort!=-1`
	}
	if types == "loan" {
		sql = `UPDATE product p SET p.loan_sort=p.loan_sort-1 WHERE p.loan_sort>? AND p.loan_sort<=? AND p.loan_sort!=-1`
	}
	if types == "largeLoan" {
		sql = `UPDATE product p SET p.large_loan_sort=p.large_loan_sort-1 WHERE p.large_loan_sort>? AND p.large_loan_sort<=? AND p.large_loan_sort!=-1`
	}
	_, err = orm.NewOrm().Raw(sql, oldSort, newSort).Exec()
	return
}

//提高排序
func UpSort(oldSort, newSort int, types string) (err error) {
	var sql string
	if types == "index" {
		sql = `UPDATE product p SET  p.sort= p.sort+1 WHERE  p.sort<? AND  p.sort>=? AND  p.sort!=-1`
	}
	if types == "loan" {
		sql = `UPDATE product p SET  p.loan_sort= p.loan_sort+1 WHERE  p.loan_sort<? AND  p.loan_sort>=? AND  p.loan_sort!=-1`
	}
	if types == "largeLoan" {
		sql = `UPDATE product p SET  p.large_loan_sort= p.large_loan_sort+1 WHERE  p.large_loan_sort<? AND  p.large_loan_sort>=? AND  p.large_loan_sort!=-1`
	}
	_, err = orm.NewOrm().Raw(sql, oldSort, newSort).Exec()
	return
}

//获取产品的展示信息
func GetProductMangeInfo(id int) (pro *ProductMange, err error) {
	o := orm.NewOrm()
	sql := `SELECT is_use,full_guide,is_pop_up,is_index_show,sort,loan_sort,large_loan_sort,recommend_sort FROM product WHERE id=?`
	err = o.Raw(sql, id).QueryRow(&pro)
	return
}

//修改产品信息
func UpdateProductInfo(pro *ProductMange) (err error) {
	sql := `UPDATE product SET is_use=?,full_guide=?,is_pop_up=?,is_index_show=?,sort=?,loan_sort=?,large_loan_sort=?,recommend_sort=?,tag=?  WHERE id=?`
	_, err = orm.NewOrm().Raw(sql, pro.IsUse, pro.FullGuide, pro.IsPopUp, pro.IsIndexShow, pro.Sort, pro.LoanSort, pro.LargeLoanSort, pro.RecommendSort, pro.Tag, pro.Id).Exec()
	if err != nil {
		return
	}
	return
}

//产品管理->编辑页面->数据保存
func UpdateProductAddress(address *ProductAddress) (err error) {
	o := orm.NewOrm()
	var sql string
	//更新address数据

	if address.PAId == 0 {
		sql = `INSERT INTO product_loan_condition (product_id,min_age,max_age,min_zmscore,address_type,address,app_permit,
		is_provident_fund,is_social_security,
	 house_property,car_property,create_time)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`
		_, err = o.Raw(sql, address.Id, address.MinAge, address.MaxAge, address.MinZmScore, address.AddressType, address.Address,
			address.AppPermit, address.IsProvidentFund, address.IsSocialSecurity, address.HouseProperty, address.CarProperty,
			time.Now()).Exec()
	} else {
		sql = `UPDATE product_loan_condition SET min_age=?,max_age=?,min_zmscore=?,address_type=?,address=?,app_permit=?,modify_time=?,
		is_provident_fund=?,is_social_security=?,
	 house_property=?,car_property=?
		WHERE id=?`
		_, err = o.Raw(sql, address.MinAge, address.MaxAge, address.MinZmScore, address.AddressType, address.Address, address.AppPermit, time.Now(),
			address.IsProvidentFund, address.IsSocialSecurity, address.HouseProperty, address.CarProperty,
			address.PAId).Exec()
	}
	return
}

//保存产品认证
func SaveProductAuth(ProductId int, authFields []string) error {
	//

	o := orm.NewOrm()
	o.Begin()
	sql := "INSERT INTO product_auth (product_id,auth_field) VALUES(?,?)"
	psql, err := o.Raw(sql).Prepare()
	if err != nil {
		return err
	}
	for _, val := range authFields {
		_, err := psql.Exec(ProductId, val)
		if err != nil {
			return err
		}
	}
	defer func() {
		psql.Close()
		if err != nil {
			o.Rollback()
			return
		} else {
			o.Commit()
		}
	}()
	return nil
}

// 更新大额API产品信息
func UpdateBigAPIProduct(product *ProductForAdd) error {
	/*fmt.Println("---")
	fmt.Println(product.LoanTermCount)*/
	product.LoanTermCountType = "个月"
	loanTermCount := strings.Replace(product.LoanTermCount, ",", product.LoanTermCountType+",", -1)
	product.LoanTermCount = loanTermCount + product.LoanTermCountType

	o := orm.NewOrm()

	sql := `UPDATE product SET name=?, cooperation_type = ?, product_logo = ?,loan_term_count = ?,loan_speed =?, min_money=?, max_money=?,
		min_loan_term=?,max_loan_term=?,reference_rate=?,large_loan_fee=?,nonrecurring_expense =?,interest_rates_details=?,request_condition=?,requested_material=?,loan_tax_fee=? WHERE id = ? `
	_, err := o.Raw(sql, product.ProName, product.CooperationType, product.ProductLogo, product.LoanTermCount, product.LoanSpeed, product.MinMoney, product.MaxMoney,
		product.MinLoanTerm, product.MaxLoanTerm, product.ReferenceRate, product.LargeLoanFee, product.NonrecurringExpense, product.InterestRatesDetails, product.RequestCondition,
		product.RequestedMaterial, product.LoanTaxFee, product.Id).Exec()

	return err
}

//根据id查询产品认证
func QueryProAuthBYId(product_id int) (list []ProductAuth, err error) {
	o := orm.NewOrm()
	sql := `SELECT product_id, auth_field FROM product_auth WHERE product_id = ? `
	_, err = o.Raw(sql, product_id).QueryRows(&list)
	return
}

//根据product_id 统计产品认证数目
func StaProAuthNumByProid(product_id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM product_auth WHERE product_id = ? `
	err = o.Raw(sql, product_id).QueryRow(&count)
	return
}

//根据proid删除所有对象
func DeleteProductAuthByProId(productId int) (err error) {
	o := orm.NewOrm()
	sql := `DELETE FROM product_auth WHERE product_id = ? `
	_, err = o.Raw(sql, productId).Exec()
	return
}

//根据产品id查询产品类型
func GetProductTypeById(productId int) (loanProductType int, err error) {
	sql := `SELECT loan_product_type FROM product WHERE id=?`
	err = orm.NewOrm().Raw(sql, productId).QueryRow(&loanProductType)
	return
}

//根据产品类型查询产品标签
func GetTagsByProductTypeAndId(loanProductType int) (tag []string, err error) {
	sql := `SELECT tag FROM product WHERE loan_product_type=? AND tag IS NOT NULL GROUP BY tag`
	_, err = orm.NewOrm().Raw(sql, loanProductType).QueryRows(&tag)
	return
}
