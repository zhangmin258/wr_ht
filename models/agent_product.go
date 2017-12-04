package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type AgentProductShow struct {
	Id            int       //
	AgentId       int       //下级代理商id
	ProId         int       //产品id
	ProductName   string    //代理产品名称
	JointMode     int       //合作模式
	AgentTime     time.Time //代理日期
	CpaPrice      float64   //推广价格
	CleaningCycle string    //结算周期
}

type AgentProduct struct {
	Id            int `orm:"pk"` //主键
	ProId         int            //产品id
	AgentId       int            //代理商id
	UrlId         string         //product_agent_url表的主键
	JointMode     int            //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CleaningType  int            //结算方式：1，对公；2，对私
	CpaDefine     string         //CPA结算的有效事件
	CpaPrice      float64        //CPA的价格
	CpsPrice      float64        //CPS的价格
	CleaningCycle string         //结算周期
	AgentTime     time.Time      //代理时间
}

type AgentProductEdit struct {
	Id            int       //主键
	ProId         int       //产品id
	ProName       string    //产品名称
	AgentId       int       //代理商id
	UrlId         string    //product_agent_url表的主键
	Url           string    //外放链接
	JointMode     int       //合作模式：1，CPA；2，CPS；3，CPA+CPS
	CleaningType  int       //结算方式：1，对公；2，对私
	CpaDefine     string    //CPA结算的有效事件
	CpaPrice      float64   //CPA的价格
	CpsPrice      float64   //CPS的价格
	CleaningCycle string    //结算周期
	AgentTime     time.Time //代理时间
}

//产品数据录入代理商信息显示
type AgentProInfo struct {
	ProName   string //产品名称
	AgentName string //代理商名称
}

//当前产品的结算信息
type ProCleaningNow struct {
	JointMode     int
	CleaningType  int
	CpaPrice      float64
	CpsFirstPer   float64
	CpsAgainPer   float64
	CpaDefine     string
	CleaningCycle string
}

func GetAgentProductList(uid int, condition string, pars []string, begin, count int) (list []AgentProductShow, err error) {
	sql := `SELECT ap.id,ap.joint_mode,ap.agent_id,ap.pro_id,ap.agent_id,ap.agent_time,ap.cleaning_cycle,ap.cpa_price,p.name AS product_name FROM agent_product ap LEFT JOIN product p
			ON ap.pro_id=p.id WHERE ap.agent_id=?`
	sql += condition
	sql += " ORDER BY ap.agent_time DESC LIMIT ?, ?"
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid, pars, begin, count).QueryRows(&list)
	return
}

// 查询所有下级代理产品数量
func GetAgentProductCount(uid int, condition string, params []string) (count int, err error) {
	sql := `SELECT count(1) FROM agent_product ap
			LEFT JOIN product p
			ON ap.pro_id=p.id
			WHERE ap.agent_id=?`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, uid, params).QueryRow(&count)
	return
}

//保存下级代理产品
func (agentproduct *AgentProduct) AgentProductAdd(ids []int) error {
	o := orm.NewOrm()
	o.Begin()
	sql := `INSERT INTO agent_product (pro_id, agent_id, url_id, joint_mode, cleaning_type, cpa_define, cpa_price, cps_price, cleaning_cycle, agent_time) VALUES (?,?,?,?, ?, ?, ?, ?, ?, ?)`
	_, err := o.Raw(sql, agentproduct.ProId, agentproduct.AgentId, agentproduct.UrlId, agentproduct.JointMode, agentproduct.CleaningType, agentproduct.CpaDefine, agentproduct.CpaPrice,
		agentproduct.CpsPrice, agentproduct.CleaningCycle, time.Now()).Exec()
	//update外放链接状态
	sql = `UPDATE product_agent_url SET is_used = 1 WHERE id =?`
	usql, err := o.Raw(sql).Prepare()
	for _, v := range ids {
		_, err = usql.Exec(v)
	}
	defer func() {
		usql.Close()
		if err != nil {
			o.Rollback()
			return
		} else {
			o.Commit()
		}
	}()
	return nil
}

// update下级代理数据
func (agentproduct *AgentProduct) EditAgentProduct(urls []int, ids []int) error {
	o := orm.NewOrm()
	o.Begin()
	//外放链接重置为0
	sql := `UPDATE product_agent_url SET is_used = 0 WHERE id =?`
	usql, err := o.Raw(sql).Prepare()
	for _, v := range ids {
		_, err = usql.Exec(v)
	}
	usql.Close()
	sql = `UPDATE agent_product SET pro_id=?, agent_id=?, url_id=?, joint_mode=?, cleaning_type=?, cpa_define=?, cpa_price=?, cps_price=?, cleaning_cycle=? WHERE id = ?`
	_, err = o.Raw(sql, agentproduct.ProId, agentproduct.AgentId, agentproduct.UrlId, agentproduct.JointMode, agentproduct.CleaningType, agentproduct.CpaDefine, agentproduct.CpaPrice,
		agentproduct.CpsPrice, agentproduct.CleaningCycle, agentproduct.Id).Exec()
	//update外放链接状态
	sql = `UPDATE product_agent_url SET is_used = 1 WHERE id =?`
	usql, err = o.Raw(sql).Prepare()
	for _, v := range ids {
		_, err = usql.Exec(v)
	}
	defer func() {
		usql.Close()
		if err != nil {
			o.Rollback()
			return
		} else {
			o.Commit()
		}
	}()
	return nil
}

// 根据id查询代理产品
func GetAgentProductById(agentProduct *AgentProduct) (err error) {
	err = orm.NewOrm().Read(agentProduct)
	return
}

// 根据id刪除代理产品
func DelAgentProById(id int, urls []int) (err error) {
	o := orm.NewOrm()
	o.Begin()
	sql := `DELETE FROM agent_product WHERE id=?`
	_, err = o.Raw(sql, id).Exec()
	//外放地址置为0
	sql = `UPDATE product_agent_url SET has_own = 1 WHERE id =?`
	usql, err := o.Raw(sql).Prepare()
	for _, v := range urls {
		_, err = usql.Exec(v)
	}
	defer func() {
		usql.Close()
		if err != nil {
			o.Rollback()
			return
		} else {
			o.Commit()
		}
	}()
	return nil
}

//根据Id查找urlId
func GetUrlductById(id int) (url string, err error) {
	sql := `SELECT url_id FROM agent_product WHERE id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&url)
	return
}

//添加代理商代理产品
func AddAgentProduct(agentProduct *AgentProduct) (err error) {
	_, err = orm.NewOrm().Insert(agentProduct)
	return
}

//修改代理商代理产品
func UpdateAgentProduct(agentProduct *AgentProduct) (err error) {
	sql := `UPDATE agent_product SET url_id = ?, joint_mode = ?, cleaning_type = ?, cpa_define = ?, cpa_price = ?, cps_price = ?, cleaning_cycle = ? WHERE id = ?`
	_, err = orm.NewOrm().Raw(sql, agentProduct.UrlId, agentProduct.JointMode, agentProduct.CleaningType, agentProduct.CpaDefine, agentProduct.CpaPrice, agentProduct.CpsPrice, agentProduct.CleaningCycle, agentProduct.Id).Exec()
	return
}

//修改代理商代理产品的id字符集
func UpdateProductIdStr(agentProduct *AgentProduct) (err error) {
	sql := `UPDATE agent_product SET url_id = ?, joint_mode = ?, cleaning_type = ?, cpa_define = ?, cpa_price = ?, cps_price = ?, cleaning_cycle = ? WHERE id = ?`
	_, err = orm.NewOrm().Raw(sql, agentProduct.UrlId, agentProduct.JointMode, agentProduct.CleaningType, agentProduct.CpaDefine, agentProduct.CpaPrice, agentProduct.CpsPrice, agentProduct.CleaningCycle, agentProduct.Id).Exec()
	return
}

//根据agentproductid获取产品名称和代理商名称
func GetAgentProInfo(agentProductId int) (agentProInfo *AgentProInfo, err error) {
	sql := `SELECT a.org_name AS agent_name,p.name AS pro_name FROM agent_product ap
	LEFT JOIN agent a
	ON ap.agent_id=a.id
	LEFT JOIN product p
	ON ap.pro_id=p.id
	WHERE ap.id=?`
	err = orm.NewOrm().Raw(sql, agentProductId).QueryRow(&agentProInfo)
	return
}

//根据产品id获得当前产品的结算信息
func GetProCleaningNow(id int) (proCleaningNow ProCleaningNow, err error) {
	sql := `SELECT dd.joint_mode,dd.cpa_price,dd.cps_first_per,dd.cps_again_per,dd.cpa_define,pc.cleaning_type,pc.cleaning_cycle FROM daily_data dd
			LEFT JOIN agent_product ap
			ON dd.agent_product_id = ap.id
			LEFT JOIN product_cleaning pc
			ON pc.product_id=ap.pro_id
			WHERE ap.pro_id =? AND ap.agent_id =0
			ORDER BY dd.date DESC LIMIT 1`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&proCleaningNow)
	return
}
