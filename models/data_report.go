package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

/*
*数据报表
*/

type DataReportShow struct {
	Id                    int
	DataTime              time.Time //数据统计的时间
	Sort                  string    //平台位置
	PlatformRegisterCount int       //我司统计注册
	AccessCount           int       //访问人数uv
	PageLoadCount         int       //页面加载次数pv
	RegisterCount         int       //平台返回注册人数
	ActivateUser          int       //激活人数
	ProId                 int       //产品id
	ProName               string    //产品名称
	CpaPrice              float64   //cpa价格
	CpsFirstPer           float64   //cps价格
	MakeLoanAmount        int       //放款金额
}

type DataReport struct {
	Id                    int
	DataTime              time.Time //数据统计的时间
	Sort                  string    //平台位置
	PlatformRegisterCount int       //我司统计注册
	AccessCount           int       //访问人数uv
	PageLoadCount         int       //页面加载次数pv
	RegisterCount         int       //平台返回注册人数
	ActivateUser          int       //激活人数
	ProId                 int       //产品id
	UvRegister            float64   //uv-注册转化率
	Income                float64   //收入
	PerUvIncome           float64   //每uv收益
	Price                 string    //当日价格
	ProName               string
}

//总数
type DataReportCount struct {
	PvAllCount               int     //pv总数
	UvAllCount               int     //uv总数
	PlatformRegisterAllCount int     //我司统计注册总数
	RegisterAllCount         int     //平台返回注册总数
	ActivateAllCount         int     //我司统计激活总数
	IncomeAll                float64 //收入总数
}

//平均
type DataReportAverage struct {
	PvAverage               float64 //平均pv数
	UvAverage               float64 //平均uv数
	PlatformRegisterAverage float64 //平均我司统计注册数
	RegisterAverage         float64 //平均平台返回注册数
	ActivateAverage         float64 //平均我司统计激活数
	IncomeAverage           float64 //平均收入数
	UvRegisterAverage       float64 //平均UV注册转化
	PerUvIncomeAverage      float64 //平均每UV收益
}

//上线产品
type OnlineProduct struct {
	Id   int
	Name string
}

func GetDataReport(proId int, condition string, params []string) (dataReportShow []DataReportShow, err error) {
	sql := `SELECT apdd.cps_first_per,apdd.cpa_price,apdd.make_loan_amount,
		cr.pro_id,cr.access_count,cr.page_load_count,
		cr.platform_register_count,cr.activate_user,cr.sort,cr.data_time
		FROM click_register cr
		LEFT JOIN (
		SELECT ap.pro_id,dd.cpa_price,dd.cps_first_per,dd.make_loan_amount,dd.date
		FROM daily_data dd
		LEFT JOIN agent_product ap ON dd.agent_product_id= ap.id) apdd ON cr.pro_id=apdd.pro_id AND cr.data_time=apdd.date
		WHERE cr.pro_id =? AND cr.data_time IS NOT NULL `
	sql += condition
	sql += ` GROUP BY cr.data_time,cr.pro_id ORDER BY cr.data_time DESC `
	_, err = orm.NewOrm().Raw(sql, proId, params).QueryRows(&dataReportShow)
	if err != nil {
		beego.Debug(err.Error())
	}
	return
}

func GetDataReportAll(condition string, params []string) (dataReportShow []DataReportShow, err error) {
	sql := `SELECT apdd.cps_first_per,apdd.cpa_price,apdd.make_loan_amount,
		cr.pro_id,cr.access_count,cr.page_load_count,
		cr.platform_register_count,cr.activate_user,cr.sort,cr.data_time,p.name AS pro_name
		FROM click_register cr
		LEFT JOIN product p
		ON p.id=cr.pro_id
		LEFT JOIN (
		SELECT ap.pro_id,dd.cpa_price,dd.cps_first_per,dd.make_loan_amount,dd.date
		FROM daily_data dd
		LEFT JOIN agent_product ap ON dd.agent_product_id= ap.id) apdd ON cr.pro_id=apdd.pro_id AND cr.data_time=apdd.date
		WHERE 1=1 `
	sql += condition
	sql += ` GROUP BY cr.data_time,cr.pro_id ORDER BY cr.data_time DESC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&dataReportShow)
	return
}

//获取上线产品
func GetOnlineProduct(condition string, params []string) (onlineProduct []OnlineProduct, err error) {
	sql := `SELECT pro_id AS id,name FROM click_register cr LEFT JOIN product p ON cr.pro_id=p.id
	WHERE 1=1 `
	sql += condition
	sql += ` GROUP BY cr.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&onlineProduct)
	return
}

//获取上线产品名称
func GetOnlineProductName(condition string, params []string) (name []string, err error) {
	sql := `SELECT name FROM click_register cr LEFT JOIN product p ON cr.pro_id=p.id
	WHERE name != "" `
	sql += condition
	sql += ` GROUP BY cr.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&name)
	return
}

//获取上线的一个默认产品
func GetOnlineProOne() (onlineProduct OnlineProduct, err error) {
	sql := `SELECT p.id,name FROM click_register cr LEFT JOIN product p ON cr.pro_id=p.id LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&onlineProduct)
	return
}

func GetRestCountByProId(proId int, condition string, params []string) (registerCount []RegisterCount, err error) {
	sql := `SELECT ap.pro_id AS pro_id,dd.register_count AS count,dd.date FROM daily_data dd LEFT JOIN agent_product ap ON dd.agent_product_id=ap.id WHERE ap.pro_id=? AND ap.agent_id=0 `
	sql += condition
	sql += ` GROUP BY date`
	_, err = orm.NewOrm().Raw(sql, proId, params).QueryRows(&registerCount)
	return
}

func GetRestCountAll(condition string, params []string) (registerCount []RegisterCount, err error) {
	sql := `SELECT ap.pro_id AS pro_id,dd.register_count AS count,dd.date AS date FROM daily_data dd LEFT JOIN agent_product ap ON dd.agent_product_id=ap.id WHERE  ap.agent_id=0  `
	sql += condition
	sql += ` GROUP BY date,ap.pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&registerCount)
	return
}
