package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type ClickRegister struct {
	Name                  string
	ProId                 int
	RegisterCount         int     //平台返回注册人数
	AccessCount           int     //访问人数
	PlatformRegisterCount int     //我司统计注册
	ActivateUser          int     //激活人数
	DataRisk              float64 //数据风险
	Sort                  string  //历史位置
	CooperationType       int
	Count                 int
	Sorting               int  //位置建议调整
	ClickEarnings         float64 //点击收益
	AllEarnings           float64 //总收益
	CpaPrice              float64 //产品单价
	CreateTime            time.Time
}

type RegisterCount struct {
	ProId int
	Date  time.Time
	Count int
}

type ProIncome struct {
	ProductId     int     //产品ID
	ClickEarnings float64 //点击收益
}
type Pros []*ProIncome

type ProductInfos []ClickRegister

func (I ProductInfos) Len() int {
	return len(I)
}
func (I ProductInfos) Less(i, j int) bool {
	return I[i].AccessCount > I[j].AccessCount
}
func (I ProductInfos) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

func (I Pros) Len() int {
	return len(I)
}
func (I Pros) Less(i, j int) bool {
	return I[i].ClickEarnings > I[j].ClickEarnings
}
func (I Pros) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//获得首页的产品数量
func GetProCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM product WHERE is_index_show=1`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

func GetClickRegister(condition string, params []string) (clickRegister []ClickRegister, err error) {
	sql := `SELECT pc.cpa_price,p.name,pro_id,SUM(register_count) AS register_count,SUM(access_count) AS access_count,SUM(platform_register_count) AS platform_register_count,SUM(activate_user) AS activate_user,c.sort,p.cooperation_type,count(c.sort) AS count
			FROM click_register c
			LEFT JOIN product p
			ON p.id=c.pro_id
			LEFT JOIN product_cleaning pc ON pc.product_id = p.id
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY c.pro_id,c.sort ORDER BY pro_id,count DESC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&clickRegister)

	return
}

func QueryRegisterCount(condition string, params []string) (registerCount []RegisterCount, err error) {
	sql := `SELECT ap.pro_id AS pro_id,SUM(dd.register_count) AS count FROM daily_data dd LEFT JOIN agent_product ap ON dd.agent_product_id=ap.id WHERE 1=1 `
	sql += condition
	sql += ` GROUP BY pro_id`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&registerCount)
	return
}
